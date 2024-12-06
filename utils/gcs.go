package utils

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fashora-backend/config"
	"fashora-backend/models"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// LoadServiceAccount đọc và parse file Service Account JSON
func LoadServiceAccount(keyFile string) (*models.ServiceAccount, error) {
	// Đọc nội dung file JSON
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %v", err)
	}

	// Khai báo cấu trúc Service Account
	var sa models.ServiceAccount
	if err := json.Unmarshal(data, &sa); err != nil {
		return nil, fmt.Errorf("failed to parse service account file: %v", err)
	}

	return &sa, nil
}

// UploadToGCS uploads a file to Google Cloud Storage and returns the public URL
func UploadToGCS(fileContent io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(config.AppConfig.GCS.BucketName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, fileContent); err != nil {
		log.Printf("Failed to upload file to GCS (object: %s): %v", objectName, err)
		return "", fmt.Errorf("failed to upload file to GCS: %v", err)
	}

	if err := writer.Close(); err != nil {
		log.Printf("Failed to finalize upload (object: %s): %v", objectName, err)
		return "", fmt.Errorf("failed to finalize upload: %v", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.AppConfig.GCS.BucketName, objectName)
	log.Printf("Successfully uploaded file to GCS. Public URL: %s", publicURL)
	return publicURL, nil
}

// GenerateSignedURL generates a signed URL for accessing an object in GCS
func GenerateSignedURL(objectName string) (string, error) {
	saData, err := LoadServiceAccount(config.AppConfig.GCS.KeyFile)
	if err != nil {
		log.Printf("Error loading service account: %v", err)
		return "", err
	}

	options := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute), // URL valid for 15 minutes
		GoogleAccessID: saData.ClientEmail,
		PrivateKey:     []byte(saData.PrivateKey),
	}

	url, err := storage.SignedURL(config.AppConfig.GCS.BucketName, objectName, options)
	if err != nil {
		log.Printf("Failed to generate signed URL (object: %s): %v", objectName, err)
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}

	log.Printf("Successfully generated signed URL: %s", url)
	return url, nil
}

// ExtractObjectName extracts the object name from a full GCS URL
func ExtractObjectName(fullURL string) (string, error) {
	prefix := fmt.Sprintf("https://storage.googleapis.com/%s/", config.AppConfig.GCS.BucketName)
	if !strings.HasPrefix(fullURL, prefix) {
		log.Printf("Invalid URL format: %s", fullURL)
		return "", fmt.Errorf("invalid URL format")
	}
	objectName := strings.TrimPrefix(fullURL, prefix)
	log.Printf("Extracted object name from URL: %s", objectName)
	return objectName, nil
}

// GenerateObjectName generates an object name for GCS
func GenerateObjectName(filename, folder string) string {
	objectName := fmt.Sprintf("%s/%d_%s", folder, time.Now().Unix(), filename)
	log.Printf("Generated object name: %s", objectName)
	return objectName
}
