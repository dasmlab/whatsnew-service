package main

import (

        "github.com/gin-gonic/gin"
        "whatsnew-service/api"
)


func initializeRoutes(router *gin.Engine) {
    log.Infof("DASMLAB WhatsNew - Initializing  API routes..")
    router.GET("/get",  api.GetCachedCommits)
}

