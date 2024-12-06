package inventory_controller

import (
	"errors"
	"fashora-backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func AddInventory(c *gin.Context) {
	storeID := c.PostForm("store_id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	imageURL := c.PostForm("image_url")
	userID := c.PostForm("user_id")

	if storeID == "" || name == "" || url == "" || imageURL == "" || userID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Missing required fields (store_id, name, url, image_url, or user_id)",
			Data:    nil,
		})
		return
	}

	inventory := models.Inventory{
		StoreID:  storeID,
		Name:     name,
		URL:      url,
		ImageURL: imageURL,
		UserID:   userID,
	}

	if err := models.DB.Create(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to add inventory",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Inventory added successfully",
		Data:    inventory,
	})
}

func DeleteInventory(c *gin.Context) {
	id := c.PostForm("item_id")

	if id == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Inventory ID is required to delete",
			Data:    nil,
		})
		return
	}

	if err := models.DB.Where("id = ?", id).Delete(&models.Inventory{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("Inventory with ID %s not found", id),
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Status:  http.StatusInternalServerError,
				Message: "Failed to delete inventory",
				Data:    err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: fmt.Sprintf("Inventory with ID %s deleted successfully", id),
		Data:    nil,
	})
}

func ListInventories(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Missing required parameter: user_id",
			Data:    nil,
		})
		return
	}

	var inventories []models.Inventory

	if err := models.DB.Where("user_id = ?", userID).Find(&inventories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch inventories for user",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Image URLs fetched successfully",
		Data:    inventories,
	})
}
