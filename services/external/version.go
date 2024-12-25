package external

import (
	"fashora-backend/config"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Version(c *gin.Context) {
	utils.SendSuccessResponse(c, http.StatusOK, "Get version completely", gin.H{
		"minimal_version": config.AppConfig.Version.MinimalVersion,
		"latest_version":  config.AppConfig.Version.LatestVersion,
	})
}
