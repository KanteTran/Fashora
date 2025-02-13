package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"

	"fashora-backend/config"
	"fashora-backend/logger"
	"fashora-backend/utils"
)

// CallTryOnAPI sends images to the Try-On API and retrieves the result URL
func CallTryOnAPI(c *gin.Context, personImageURL, clothImageURL, maskURL string) {
	body, contentType, err := createMultipartData(map[string]string{
		"person_image_url": personImageURL,
		"cloth_image_url":  clothImageURL,
		"mask_url":         maskURL,
		"access_token":     RefreshTokenGcp(),
	})
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := sendPostRequest(config.AppConfig.Model.GenAPI, body, contentType)
	if err != nil {
		utils.SendErrorResponse(nil, http.StatusBadGateway, err.Error())
		return
	}
	defer resp.Body.Close()

	processAPIResponse(c, resp)
}

// RefreshTokenGcp retrieves a fresh access token for Google Cloud Platform
func RefreshTokenGcp() string {
	credentialsFilePath := config.AppConfig.GCS.KeyFile

	data, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		logger.Infof("Failed to read credentials file: %v\n", err)
		os.Exit(1)
	}
	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		logger.Infof("Failed to parse credentials file: %v\n", err)
		os.Exit(1)
	}

	tokenSource := config.TokenSource(context.Background())
	token, err := tokenSource.Token()
	if err != nil {
		fmt.Printf("Failed to retrieve token: %v\n", err)
		os.Exit(1)
	}

	return token.AccessToken
}

// createMultipartData generates a multipart form body
func createMultipartData(fields map[string]string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", fmt.Errorf("failed to add field '%s': %v", key, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close writer: %v", err)
	}

	return body, writer.FormDataContentType(), nil
}

// sendPostRequest sends a POST request with the given body and content type
func sendPostRequest(url string, body *bytes.Buffer, contentType string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	return resp, nil
}

// processAPIResponse processes the response from the Try-On API
func processAPIResponse(c *gin.Context, resp *http.Response) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("failed to read response body: %v", err))
		return
	}

	var responseJSON map[string]interface{}
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		utils.SendErrorResponse(nil, http.StatusInternalServerError,
			fmt.Sprintf("failed to decode JSON: %v", err))
		return
	}

	resultURL, ok := responseJSON["result_url"].(string)
	if !ok {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "result_url not found in response or is not a string")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Successfully fetched result URL", map[string]string{
		"result_url": resultURL,
	})
}
