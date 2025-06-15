package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type CommitEntry struct {
	ID      int    `json:"id"`
	Project string `json:"project"`
	Title   string `json:"title"`
	Date    string `json:"date"`
}

var cachedCommits []CommitEntry
var githubAccessToken string

func SetGitHubAccessToken(tok string) {
	githubAccessToken = tok
}

func getAccessibleRepos(token string) ([]string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/installation/repositories", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list repos: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Repositories []struct {
			FullName string `json:"full_name"`
		} `json:"repositories"`
	}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid repo response: %v", err)
	}

	var repos []string
	for _, r := range result.Repositories {
		repos = append(repos, r.FullName)
	}
	return repos, nil
}

func fetchCommits(repo string) ([]CommitEntry, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/commits?per_page=2", repo)

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+githubAccessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("GitHub API error for %s: %v", repo, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Warnf("Non-200 response for %s: %d\nBody: %s", repo, resp.StatusCode, body)
		return nil, fmt.Errorf("GitHub API error %d", resp.StatusCode)
	}

	var apiResp []struct {
		Commit struct {
			Message string    `json:"message"`
			Author  struct {
				Date time.Time `json:"date"`
			} `json:"author"`
		} `json:"commit"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var commits []CommitEntry
	for i, c := range apiResp {
		commits = append(commits, CommitEntry{
			ID:      i + 1,
			Project: repo,
			Title:   c.Commit.Message,
			Date:    c.Commit.Author.Date.Format("2006-01-02"),
		})
	}
	return commits, nil
}

func RefreshCommitCacheDynamic(repoList []string) {
	log.Infof("🔄 Refreshing commit cache for %d repos...", len(repoList))

	var results []CommitEntry
	idCounter := 1

	for _, repo := range repoList {
		log.Debugf("📦 Processing repo: %s", repo)
		commits, err := fetchCommits(repo)
		if err != nil {
			log.Warnf("⚠️ Skipped %s: %v", repo, err)
			continue
		}
		for _, c := range commits {
			c.ID = idCounter
			idCounter++
			results = append(results, c)
		}
	}

	cachedCommits = results
	log.Infof("✅ Commit cache updated. Total commits: %d", len(results))
}

func ListAccessibleRepos(token string) ([]string, error) {
	url := "https://api.github.com/installation/repositories"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Warnf("GitHub API returned %d\nBody: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var result struct {
		Repositories []struct {
			FullName string `json:"full_name"`
		} `json:"repositories"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON unmarshal failed: %v", err)
	}

	var repos []string
	for _, r := range result.Repositories {
		repos = append(repos, r.FullName)
	}

	log.Infof("🔍 Retrieved %d accessible repos", len(repos))
	return repos, nil
}
