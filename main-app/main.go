package main

import (
	// STD
	"os"
	"time"
	"strings"

	// 3PPs
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
        swaggerFiles "github.com/swaggo/files"
        "github.com/Depado/ginprom"
        "github.com/gin-contrib/cors"

	// Our Stuff
	"whatsnew-service/logutil"
	"whatsnew-service/docs"
	"whatsnew-service/api"

)

// VARS
const version = "0.0.1"
var component_name = "whatsnew-svc-main"
var log = logutil.InitLogger(component_name)




// @title MCP Explorer - MCP Server APIs
// @version 0.0.1
// @description APIs for MCP Server Instantiation, Configuration and Handling
// @BasePath /


func main() {
	log.Infof("DASMLAB WhatsNew Service - Starting %s", component_name)
	docs.SwaggerInfo.Version = version

	// Set gin Prod mode
	gin.SetMode(gin.ReleaseMode)


	// Check ENV Vars
	svcSecret := "whatsnew-svc-secret.txt"
	pemFile := "whatsnew-github.pem"


	// Github OAUTH SETUP
	var appID, instID string
	data, err := os.ReadFile(svcSecret)
	if err != nil {
		log.Fatalf("Failed to read: %s: %v", svcSecret, err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "APP_ID=") {
			appID = strings.TrimPrefix(line, "APP_ID=")
		}
		if strings.HasPrefix(line, "INSTALLATION_ID=") {
			instID = strings.TrimPrefix(line, "INSTALLATION_ID=")
		}
	}

	ghAuth := &api.GitHubAppAuth{
		AppID: appID,
		InstallationID: instID,
		PrivateKeyPath: pemFile,
	}


	// Read in the list of repos to monitor
	repoList, err := api.ReadTargetRepos("repos/target_repos.txt")
	if err != nil {
		log.Fatalf("Failed to read repos: %v", err)
	}

	// Start up the CommitCache Backgrounder
	go func() {
		for {
			token, err := ghAuth.GetAccessToken()
			if err != nil {
				log.Errorf("Token Fetch failed: %v", err)
				time.Sleep(1 * time.Minute)
				continue
			}
			log.Info("Refreshing cache...")
			api.SetGitHubAccessToken(token)
			api.RefreshCommitCache(repoList)
			time.Sleep(1 * time.Minute)
		}
	}()


	// Primary App Router
	mainRouter := gin.Default()
	
	// Allow CORS
	mainRouter.Use(cors.Default()) // Allows all - Depado rocks!

	// Metrics (out of band) Router
	metricsRouter := gin.Default()

	// ginprim hooks
	p := ginprom.New(
		ginprom.Engine(metricsRouter),
		ginprom.Subsystem("gin"),
                        ginprom.Path("/metrics"),
        )

        // Wrap our mainRouter
        mainRouter.Use(p.Instrument())

        // Add our Custom Metrics - turned off for now
        //metrics.RegisterCustomMetrics()

	// Add Swagger UI Route
	mainRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Init Routes
	initializeRoutes(mainRouter)

	// Launch metricsROuter

	go func() {
		log.Infof("Starting metrics server on :9200")
		if err := metricsRouter.Run(":9200"); err != nil {
			log.Fatalf("Metrics Server Error: %v", err)
		}
	}()

	// Launch MainRouter


	log.Info("Start main Server listening on :10020")
	if err := mainRouter.Run(":10020"); err != nil {
		log.Fatalf("Main Server Error: %v", err)
	}

}

