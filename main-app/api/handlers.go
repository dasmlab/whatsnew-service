package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCachedCommits godoc
// @Summary Get latest commits from cache
// @Description Returns cached GitHub commit info across configured repos
// @Tags whatsnew
// @Produce json
// @Success 200 {object} map[string][]CommitEntry
// @Router /get [get]
func GetCachedCommits(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"latest_commits": cachedCommits,
	})
}

