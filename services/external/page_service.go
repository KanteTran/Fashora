package external

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fashora-backend/database"
	"fashora-backend/models"
	"fashora-backend/utils"
)

func HomePage(c *gin.Context) {
	var stores []models.Stores

	if err := database.GetDBInstance().DB().Find(&stores).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not fetch stores")
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"stores": stores,
	})
}

// CreateStorePage renders the page for creating a new store
func CreateStorePage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_store.html", nil)
}
