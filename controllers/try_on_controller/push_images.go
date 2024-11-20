package try_on_controller

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fashora-backend/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var clients = make(map[string]*websocket.Conn)
var userChannels = make(map[string]chan string)

func handleWebSocket(c *gin.Context) {
	// Lấy userID từ query string
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	// Nâng cấp kết nối lên WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}
	defer conn.Close()

	// Lưu kết nối WebSocket vào map với userID là key
	clients[userID] = conn

	// Tạo channel cho userID nếu chưa tồn tại
	if _, exists := userChannels[userID]; !exists {
		userChannels[userID] = make(chan string)
		log.Printf("Channel created for user: %s", userID)
	}

	// Goroutine để lắng nghe tin nhắn từ channel và gửi tới WebSocket client
	go func(userID string) {
		for message := range userChannels[userID] {
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Error sending WebSocket message:", err)
				return
			}
		}
	}(userID)

	// Đợi tin nhắn từ client qua WebSocket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket connection closed:", err)
			delete(clients, userID)
			delete(userChannels, userID)
			break
		}
		fmt.Printf("Received message from user %s: %s\n", userID, string(msg))
	}
}

func UploadImages(c *gin.Context) {
	// List of form keys and corresponding bucket names
	images := []struct {
		formKey    string
		bucketName string
	}{
		{"people", config.AppConfig.GscFolderPeople},
		{"clothes", config.AppConfig.GscFolderClothes},
		{"posh", config.AppConfig.GscFolderPosh},
	}

	// Initialize Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.AppConfig.GscKeyFile))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GCS client"})
		return
	}
	defer client.Close()

	// Initialize a map to store the URLs of uploaded images
	imageURLs := make(map[string]string)

	// Loop over each image and upload it to the corresponding bucket
	for _, image := range images {
		// Retrieve the file from the form
		file, err := c.FormFile(image.formKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No file uploaded for %s", image.formKey)})
			return
		}

		// Open the file
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to open file for %s", image.formKey)})
			return
		}
		defer fileContent.Close()

		// Define the object name in GCS
		objectName := fmt.Sprintf("%s/%d_%s", image.formKey, time.Now().Unix(), file.Filename)

		// Upload the file to GCS
		bucket := client.Bucket(config.AppConfig.GscBucketName)
		object := bucket.Object(objectName)
		writer := object.NewWriter(ctx)
		if _, err := io.Copy(writer, fileContent); err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload image %s to GCS", err)})
			return
		}
		if err := writer.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to finalize upload for %s", err)})
			return
		}

		// Construct the public URL (assuming public bucket settings)
		imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.AppConfig.GscBucketName, objectName)
		imageURLs[image.formKey] = imageURL
	}

	// Respond with URLs of uploaded images
	c.JSON(http.StatusCreated, gin.H{
		"message": "Images uploaded successfully",
		"urls":    imageURLs,
	})

}

type ImageRequest struct {
	Image1 []byte `json:"image1"` // First image in binary form
	Image2 []byte `json:"image2"` // Second image in binary form
	Image3 []byte `json:"image3"` // Third image in binary form
}

type ImageResponse struct {
	Status   string `json:"status"`
	ImageURL string `json:"image_url"` // Processed image URL or data
}

type ServiceAccount struct {
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

// GetImageURL generates a signed URL to access the image in GCS
func GetImageURL(c *gin.Context) {
	fullURL := c.Query("filename")
	if fullURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename URL is required"})
		return
	}

	// Extract object name from full URL
	objectName, err := extractObjectName(fullURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define GCS bucket name
	bucketName := config.AppConfig.GscBucketName

	// Generate signed URL for the image
	url, err := generateSignedURL(bucketName, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate signed URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image URL generated successfully",
		"url":     url,
	})
}

func extractObjectName(fullURL string) (string, error) {
	prefix := "https://storage.googleapis.com/fashionira/"
	if !strings.HasPrefix(fullURL, prefix) {
		return "", fmt.Errorf("Invalid URL format")
	}
	return strings.TrimPrefix(fullURL, prefix), nil
}

func generateSignedURL(bucketName, objectName string) (string, error) {
	// Load the service account credentials from JSON file
	data, err := ioutil.ReadFile("smart-exchange-441906-p0-c0277d140202.json")
	if err != nil {
		return "", fmt.Errorf("failed to read service account file: %v", err)
	}

	var sa ServiceAccount
	if err := json.Unmarshal(data, &sa); err != nil {
		return "", fmt.Errorf("failed to parse service account file: %v", err)
	}

	// Set up options for the signed URL
	options := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute), // URL valid for 15 minutes
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     []byte(sa.PrivateKey),
	}

	// Generate the signed URL
	url, err := storage.SignedURL(bucketName, objectName, options)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}
	return url, nil
}
