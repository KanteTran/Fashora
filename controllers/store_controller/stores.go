package store_controller

import (
	"errors"
	"fashora-backend/config"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateStore(c *gin.Context) {
	phone := c.PostForm("phone")
	storeName := c.PostForm("store_name")
	address := c.PostForm("address")
	description := c.PostForm("description")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get image"})
		return
	}

	tx := models.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	store := models.Stores{
		Phone:       phone,
		StoreName:   storeName,
		Address:     address,
		Description: description,
		Status:      1,
	}

	if err := models.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create store"})
		return
	}

	err = external.CreateFoldersIfNotExists(config.AppConfig.GscBucketName, fmt.Sprintf("stores/%s", store.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create store's cloud folder"})
		tx.Rollback()
		return
	}

	fileName := fmt.Sprintf("%s/%s/%s", config.AppConfig.GscBucketName, fmt.Sprintf("stores/%s", store.Id), file.Filename)
	url, err := external.UploadImageToGCS(fileName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not upload image: %s", err)})
		tx.Rollback()
		return
	}

	store.UrlImage = url
	if err := tx.Save(&store).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update store with image URL"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Redirect(http.StatusFound, "/stores")
}

func AddItemPage(c *gin.Context) {
	var stores []models.Stores
	if err := models.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stores"})
		return
	}

	c.HTML(http.StatusOK, "add_item.html", gin.H{
		"stores": stores,
	})
}

func AddItem(c *gin.Context) {
	storeID := c.PostForm("store_id")
	name := c.PostForm("name")
	url := c.PostForm("url")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get image"})
		return
	}

	var store models.Stores
	err = models.DB.Where("id = ?", storeID).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Store does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
	}

	err = external.CreateFoldersIfNotExists(config.AppConfig.GscBucketName, fmt.Sprintf("stores/%s/items", storeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not create store's items folder: %s", err)})
		return
	}

	fileName := fmt.Sprintf("%s/%s/%s", config.AppConfig.GscBucketName, fmt.Sprintf("stores/%s/items", store.Id), file.Filename)
	imageUrl, err := external.UploadImageToGCS(fileName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not upload image: %s", err)})
		return
	}

	item := models.Item{
		StoreID:  storeID,
		Name:     name,
		URL:      url,
		ImageURL: imageUrl,
	}

	if err := models.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Item added successfully",
		"item_id":   item.ID,
		"store_id":  item.StoreID,
		"name":      item.Name,
		"url":       item.URL,
		"image_url": imageUrl,
	})
	c.Redirect(http.StatusFound, "/stores/add-item?success=true")
}

func ListStores(c *gin.Context) {
	var stores []models.Stores

	if err := models.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch stores",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Stores fetched successfully",
		"data":    stores,
	})
}

func GetStoreItemsById(c *gin.Context) {
	// Lấy `id` từ tham số URL
	//storeID := c.Param("id")
	storeID := c.Query("id")

	fmt.Printf(storeID)
	// Kiểm tra xem cửa hàng có tồn tại không
	var store models.Stores
	if err := models.DB.Where("Id = ?", storeID).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := models.Response{
				Success: false,
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("Store with ID %s not found", storeID),
			}
			c.JSON(http.StatusNotFound, response)
			return
		}
		response := models.Response{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch store",
			Data:    err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Truy vấn danh sách items của store
	var items []models.Item
	if err := models.DB.Where("store_id = ?", storeID).Find(&items).Error; err != nil {
		response := models.Response{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch items for the store",
			Data:    err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Tạo phản hồi JSON chứa thông tin cửa hàng và danh sách items
	response := models.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Store and items fetched successfully",
		Data: gin.H{
			"store": gin.H{
				"id":          store.Id,
				"store_name":  store.StoreName,
				"phone":       store.Phone,
				"address":     store.Address,
				"description": store.Description,
				"url_image":   store.UrlImage,
				"status":      store.Status,
			},
			"items": items,
		},
	}
	c.JSON(http.StatusOK, response)
}

func GetItemsById(c *gin.Context) {
	itemID := c.Query("id")
	if itemID == "" {
		response := models.Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Item ID is missing in the request",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var item models.Item
	if err := models.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := models.Response{
				Success: false,
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("Item with ID %s not found", itemID),
			}
			c.JSON(http.StatusNotFound, response)
			return
		}

		response := models.Response{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to fetch item",
			Data:    err.Error(),
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := models.Response{
		Success: true,
		Status:  http.StatusOK,
		Message: "Item fetched successfully",
		Data: gin.H{
			"id":        item.ID,
			"store_id":  item.StoreID,
			"name":      item.Name,
			"url":       item.URL,
			"image_url": item.ImageURL,
		},
	}
	c.JSON(http.StatusOK, response)
}
