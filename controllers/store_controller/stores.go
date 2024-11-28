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
	// Lấy dữ liệu từ form
	phone := c.PostForm("phone")
	storeName := c.PostForm("store_name")
	address := c.PostForm("address")
	description := c.PostForm("description")
	password := c.PostForm("password")
	status := c.PostForm("status")

	// Chuyển đổi status sang số nguyên
	storeStatus := parseID(status)

	// Nhận file ảnh từ form
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

	// Tạo một bản ghi store mới
	store := models.Store{
		Phone:       phone,
		StoreName:   storeName,
		Address:     address,
		Description: description,
		Password:    password,
		Status:      storeStatus,
	}

	// Lưu vào DB
	if err := models.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create store"})
		return
	}

	// Upload file lên GCS
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
	// Update the store record with the image URL
	if err := tx.Save(&store).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update store with image URL"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Redirect(http.StatusFound, "/stores")
}

func AddItemPage(c *gin.Context) {
	var stores []models.Store
	if err := models.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stores"})
		return
	}

	c.HTML(http.StatusOK, "add_item.html", gin.H{
		"stores": stores,
	})
}

func AddItem(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form data"})
		return
	}

	// Lấy dữ liệu từ form
	storeID := c.PostForm("store_id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	productCode := c.PostForm("product_code")

	// Nhận file ảnh từ form
	files := form.File["images[]"]

	var store models.Store
	err = models.DB.Where("id = ?", storeID).First(&store).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Store does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
	}

	// Tạo thư mục trong GCS theo store ID
	err = external.CreateFoldersIfNotExists(config.AppConfig.GscBucketName, fmt.Sprintf("stores/%s/items", storeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not create store's items folder: %s", err)})
		return
	}

	// Upload ảnh và lấy URL
	var imageURLs []string
	folderName := fmt.Sprintf("%s/stores/%s/items", config.AppConfig.GscBucketName, storeID)
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", folderName, file.Filename)
		imageURL, _ := external.UploadImageToGCS(filePath, file)
		imageURLs = append(imageURLs, imageURL)
	}

	// Không upload được ảnh nào
	if len(imageURLs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not upload any image of item"})
		return
	}

	// Tạo một bản ghi item mới
	item := models.Item{
		StoreID:     parseID(storeID),
		Name:        name,
		URL:         url,
		ImageURLs:   imageURLs, // Chỉ chứa URL của ảnh đã upload
		ProductCode: productCode,
	}

	// Lưu vào DB
	if err := models.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item"})
		return
	}

	// Trả về phản hồi thành công
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Item added successfully",
		"item_id":    item.ID,
		"store_id":   item.StoreID,
		"name":       item.Name,
		"url":        item.URL,
		"image_urls": imageURLs,
	})
}

// Helper to parse IDs
func parseID(input string) int {
	var id int
	fmt.Sscanf(input, "%d", &id)
	return id
}

func ListStores(c *gin.Context) {
	var stores []models.Store

	// Lấy danh sách các cửa hàng từ DB
	if err := models.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch stores",
			"error":   err.Error(),
		})
		return
	}

	// Trả về danh sách cửa hàng dưới dạng JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Stores fetched successfully",
		"data":    stores,
	})
}
