package external

import (
	"fashora-backend/models"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HomePage renders the homepage with a list of stores
func HomePage(c *gin.Context) {
	var stores []models.Stores

	if err := models.DB.Find(&stores).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not fetch stores")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Stores fetched successfully", stores)
}

// CreateStorePage renders the page for creating a new store
func CreateStorePage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_store.html", nil)
}
