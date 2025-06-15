package api

import (
	"os"
	"bufio"
	"strings"
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	
	// Our stuff
	"whatsnew-service/logutil"
)

type CommitEntry struct {
        ID      int    `json:"id"`
        Project string `json:"project"`
        Title   string `json:"title"`
        Date    string `json:"date"`
}

var cachedCommits []CommitEntry
var githubAccessToken string
var log = logutil.InitLogger("whatsnew-svc-api")


func ReadTargetRepos(path string) ([]string, error) {
        file, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        defer file.Close()

        var repos []string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                line := strings.TrimSpace(scanner.Text())
                if line != "" {
                        repos = append(repos, line)
                }
        }
        return repos, scanner.Err()
}

func fetchCommitsFromGitHub(repo string) ([]CommitEntry, error) {
        repoName := strings.TrimPrefix(repo, "https://github.com/")
        apiURL := fmt.Sprintf("https://api.github.com/repos/%s/commits?per_page=2", repoName)

        req, _ := http.NewRequest("GET", apiURL, nil)
        req.Header.Set("Authorization", "Bearer "+githubAccessToken)
        req.Header.Set("Accept", "application/vnd.github+json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                log.Errorf("GitHub API error: %v", err)
                return nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
                log.Errorf("GitHub API status: %s for repo %s", resp.Status, repoName)
                return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
        }

        var apiResponse []struct {
                Commit struct {
                        Message string    `json:"message"`
                        Author  struct {
                                Date time.Time `json:"date"`
                        } `json:"author"`
                } `json:"commit"`
        }
        err = json.NewDecoder(resp.Body).Decode(&apiResponse)
        if err != nil {
                return nil, err
        }

        var commits []CommitEntry
        for i, c := range apiResponse {
                commits = append(commits, CommitEntry{
                        ID:      i + 1,
                        Project: repoName,
                        Title:   c.Commit.Message,
                        Date:    c.Commit.Author.Date.Format("2006-01-02"),
                })
        }
        return commits, nil
}

func RefreshCommitCache(repoList []string) {
        var newCommits []CommitEntry
        idCounter := 1

        for _, repo := range repoList {
                commits, err := fetchCommitsFromGitHub(repo)
                if err != nil {
                        continue
                }
                for _, c := range commits {
                        c.ID = idCounter
                        idCounter++
                        newCommits = append(newCommits, c)
                }
        }
        cachedCommits = newCommits
}

