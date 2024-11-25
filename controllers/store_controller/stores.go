package store_controller

import (
	"fashora-backend/config"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

	print("bat dau nhanh xong r ne")

	// Upload file lên GCS
	err = external.CreateFoldersIfNotExists(config.AppConfig.GscBucketName, fmt.Sprintf("stores/%s", storeName))
	if err != nil {
		return
	}
	fileName := fmt.Sprintf("%s/%s/%s", config.AppConfig.GscBucketName, fmt.Sprintf("stores/%s", storeName), file.Filename)
	url, err := external.UploadImageToGCS(fileName, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not upload image: %s", err)})
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
		UrlImage:    url, // Lưu URL vào DB
	}

	// Lưu vào DB
	if err := models.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create store"})
		return
	}
	print("luuw oke")

	c.Redirect(http.StatusFound, "/stores")
}

func AddItemPage(c *gin.Context) {
	var stores []models.Store
	// Lấy danh sách store từ database
	if err := models.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stores"})
		return
	}

	// Render trang HTML với danh sách stores
	c.HTML(http.StatusOK, "add_item.html", gin.H{
		"stores": stores,
	})
}

//func AddItem(c *gin.Context) {
//	// Lấy dữ liệu từ form
//	storeID := c.PostForm("store_id")
//	name := c.PostForm("name")
//	url := c.PostForm("url")
//	productCode := c.PostForm("product_code")
//
//	// Nhận file ảnh từ form
//	file, err := c.FormFile("image")
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not get image file"})
//		return
//	}
//
//	// Upload ảnh lên GCS
//	bucketName := config.AppConfig.GscBucketName
//	folderName := fmt.Sprintf("stores/%s/items", storeID) // Tạo thư mục trong GCS theo store ID
//	fileName := fmt.Sprintf("%s/%s", folderName, file.Filename)
//
//	// Upload ảnh và lấy URL
//	imageURL, err := external.UploadImageToGCS(bucketName, fileName, file)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not upload image: %s", err)})
//		return
//	}
//
//	// Tạo một bản ghi item mới
//	item := models.Item{
//		StoreID:     parseID(storeID),
//		Name:        name,
//		URL:         url,
//		ImageURLs:   []string{imageURL}, // Chỉ chứa URL của ảnh đã upload
//		ProductCode: productCode,
//	}
//
//	// Lưu vào DB
//	if err := models.DB.Create(&item).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item"})
//		return
//	}
//
//	// Trả về phản hồi thành công
//	c.JSON(http.StatusCreated, gin.H{
//		"message":   "Item added successfully",
//		"item_id":   item.ID,
//		"store_id":  item.StoreID,
//		"name":      item.Name,
//		"url":       item.URL,
//		"image_url": imageURL,
//	})
//}
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
