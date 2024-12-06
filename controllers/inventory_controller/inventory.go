package inventory_controller

import (
	"errors"
	"fashora-backend/models"
	"fashora-backend/services/auth_service"
	"fashora-backend/utils"
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
	user, err := auth_service.GetAuthenticatedUser(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	userID = user.Id

	if storeID == "" || name == "" || url == "" || imageURL == "" || userID == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Missing required fields (store_id, name, url, image_url, or user_id)")
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
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to add inventory")
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Inventory added successfully", inventory)
	return
}

func DeleteInventory(c *gin.Context) {
	id := c.PostForm("item_id")
	user, err := auth_service.GetAuthenticatedUser(c)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if id == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Inventory ID is required to delete")
		return
	}

	if err := models.DB.Where("id = ? AND user_id = ?", id, user.Id).Delete(&models.Inventory{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Inventory with ID %s not found", id))
			return
		} else {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete inventory")
			return
		}
	}
	utils.SendSuccessResponse(c, http.StatusOK, fmt.Sprintf("Inventory with ID %s deleted successfully", id), nil)

}

func ListInventories(c *gin.Context) {
	user, _ := auth_service.GetAuthenticatedUser(c)
	if user.Id == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var inventories []models.Inventory

	if err := models.DB.Where("user_id = ?", user.Id).Find(&inventories).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to get inventory")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Image URLs fetched successfully", inventories)
	return
}
