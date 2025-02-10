package recommend

import (
	"fashora-backend/database"
	"fashora-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type GetItemsByTagsRequest struct {
	Tags []int64 `json:"tags"`
}

func GetItemsByTags(c *gin.Context) {
	var req GetItemsByTagsRequest

	// Parse JSON input từ request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.Tags) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tags list cannot be empty"})
		return
	}

	var items []models.Item
	err := database.GetDBInstance().DB().Where("tags && ?", pq.Array(req.Tags)).Find(&items).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
	})
}
