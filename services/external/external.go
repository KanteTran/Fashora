package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fashora-backend/config"
	"fashora-backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"io"
	_ "io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func CallTryOnAPI(personImageURL, clothImageURL, maskURL string) models.APIResponse {
	// Create a multipart form request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	fields := map[string]string{
		"person_image_url": personImageURL,
		"cloth_image_url":  clothImageURL,
		"mask_url":         maskURL,
		"access_token":     RefreshTokenGcp(), // Ensure RefreshTokenGcp() provides the access token
	}

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return models.CreateErrorResponse(
				http.StatusInternalServerError,
				fmt.Sprintf("failed to add field '%s': %v", key, err))
		}
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return models.CreateErrorResponse(
			http.StatusInternalServerError, fmt.Sprintf("failed to close writer: %v", err))
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", config.AppConfig.ModelGenAPI, body)
	if err != nil {
		return models.CreateErrorResponse(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to create request: %v", err))
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.CreateErrorResponse(
			http.StatusBadGateway,
			fmt.Sprintf("failed to send request: %v", err))
	}
	defer resp.Body.Close()

	// Read response body
	bodyNew, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.CreateErrorResponse(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to read response body: %v", err))
	}

	// Decode JSON response
	var responseJSON map[string]interface{}
	err = json.Unmarshal(bodyNew, &responseJSON)
	if err != nil {
		return models.CreateErrorResponse(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to decode JSON: %v", err))
	}

	// Extract "result_url"
	resultURL, ok := responseJSON["result_url"].(string)
	if !ok {
		return models.CreateErrorResponse(
			http.StatusInternalServerError,
			"result_url not found in response or is not a string")
	}

	return models.APIResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Successfully fetched result URL",
		Data: map[string]string{
			"result_url": resultURL,
		},
	}
}

func RefreshTokenGcp() string {
	credentialsFilePath := config.AppConfig.GscKeyFile

	data, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		log.Fatalf("Failed to read credentials file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatalf("Failed to parse credentials file: %v", err)
	}

	tokenSource := config.TokenSource(context.Background())

	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Failed to retrieve token: %v", err)
	}

	return token.AccessToken
}

func HomePage(c *gin.Context) {
	var stores []models.Stores

	if err := models.DB.Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stores"})
		return
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"stores": stores,
	})
}

func CreateStorePage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_store.html", nil)
}
