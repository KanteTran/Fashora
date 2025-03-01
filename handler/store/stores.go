package store

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"fashora-backend/config"
	"fashora-backend/database"
	"fashora-backend/handler/tagging"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fashora-backend/utils"
)

type ImplHandlerStore interface {
	CreateStore(c *gin.Context)
	AddItemPage(c *gin.Context)
	AddItem(c *gin.Context)
	ListStores(c *gin.Context)
	GetStoreItemsById(c *gin.Context)
	GetItemsById(c *gin.Context)
}

type HandlerStore struct {
	db database.DBAdapter
}

func NewHandlerStore() ImplHandlerStore {
	return &HandlerStore{
		db: database.GetDBInstance(),
	}
}

func (h *HandlerStore) CreateStore(c *gin.Context) {
	phone := c.PostForm("phone")
	storeName := c.PostForm("store_name")
	address := c.PostForm("address")
	description := c.PostForm("description")
	typeStore := c.PostForm("type")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get image"})
		return
	}

	tx := h.db.DB().Begin()
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
		Type:        typeStore,
	}

	if err := h.db.DB().Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create store"})
		return
	}

	err = external.CreateFoldersIfNotExists(config.AppConfig.GCS.BucketName, fmt.Sprintf("stores/%s", store.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create store's cloud folder"})
		tx.Rollback()
		return
	}

	fileName := fmt.Sprintf("%s/%s/%s", config.AppConfig.GCS.BucketName, fmt.Sprintf("stores/%s", store.Id), file.Filename)
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

func (h *HandlerStore) AddItemPage(c *gin.Context) {
	var stores []models.Stores
	if err := h.db.DB().Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stores"})
		return
	}

	c.HTML(http.StatusOK, "add_item.html", gin.H{
		"stores": stores,
	})
}

func (h *HandlerStore) AddItem(c *gin.Context) {
	storeID := c.PostForm("store_id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	description := c.PostForm("description")

	file, err := c.FormFile("image")

	tags := tagging.TagClothes(file)

	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Could not get image")
		return
	}

	var store models.Stores
	err = h.db.DB().Where("id = ?", storeID).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Store does not exist")
		} else {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Database error")
		}
	}

	err = external.CreateFoldersIfNotExists(config.AppConfig.GCS.BucketName, fmt.Sprintf("stores/%s/items", storeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not create store's items folder: %s", err)})
		return
	}

	fileName := fmt.Sprintf(
		"%s/%s/%s", config.AppConfig.GCS.BucketName,
		fmt.Sprintf("stores/%s/items", store.Id),
		file.Filename)
	imageUrl, err := external.UploadImageToGCS(fileName, file)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Could not upload image: %s", err))
		return
	}

	item := models.Item{
		StoreID:     storeID,
		Name:        name,
		URL:         url,
		ImageURL:    imageUrl,
		Description: description,
		Tags:        pq.Int64Array(tags[:]),
	}

	if err := h.db.DB().Create(&item).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not add item to store")
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "Item added successfully", item)
	c.Redirect(http.StatusFound, "/stores/add-item?success=true")
}

func (h *HandlerStore) ListStores(c *gin.Context) {
	var stores []models.Stores

	storeType := c.Query("type")

	query := h.db.DB()
	if storeType != "" {
		query = query.Where("type = ?", storeType)
	}

	if err := query.Find(&stores).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch stores")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Stores fetched successfully", stores)
}

func (h *HandlerStore) GetStoreItemsById(c *gin.Context) {
	storeID := c.Query("id")

	var store models.Stores
	if err := h.db.DB().Where("Id = ?", storeID).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Store with ID %s not found", storeID))
			return
		}

		utils.SendErrorResponse(c, http.StatusNotFound, "Failed to fetch store")
		return

	}

	var items []models.Item
	if err := h.db.DB().Where("store_id = ?", storeID).Find(&items).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch items for the store")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Store and items fetched successfully",
		gin.H{
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
		})
}

func (h *HandlerStore) GetItemsById(c *gin.Context) {
	itemID := c.Query("id")

	if itemID == "" {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Item ID is missing in the request")
		return
	}

	var item models.Item
	if err := h.db.DB().Where("id = ?", itemID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Item with ID %s not found", itemID))
			return
		}

		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch item")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Item fetched successfully", gin.H{
		"id":          item.ID,
		"store_id":    item.StoreID,
		"name":        item.Name,
		"url":         item.URL,
		"image_url":   item.ImageURL,
		"description": item.Description,
	})
}
