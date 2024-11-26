package image_controller

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fashora-backend/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
		return
	}
	defer fileContent.Close()

	// Initialize Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.AppConfig.GscKeyFile))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GCS client"})
		return
	}
	defer client.Close()

	// Define GCS bucket and object name
	bucketName := config.AppConfig.GscBucketName // Replace with your bucket name
	objectName := fmt.Sprintf("images/%d_%s", time.Now().Unix(), file.Filename)

	// Upload the file to GCS
	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)
	if _, err := io.Copy(writer, fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to GCS,	"})
		return
	}
	if err := writer.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize upload"})
		return
	}

	// Construct the public URL
	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Image uploaded successfully",
		"url":     imageURL,
	})
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
	prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", config.AppConfig.GscBucketName)
	if !strings.HasPrefix(fullURL, prefix) {
		return "", fmt.Errorf("Invalid URL format")
	}
	return strings.TrimPrefix(fullURL, prefix), nil
}

func generateSignedURL(bucketName, objectName string) (string, error) {
	// Load the service account credentials from JSON file
	data, err := ioutil.ReadFile(fmt.Sprintf("%s", config.AppConfig.GscKeyFile))
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
