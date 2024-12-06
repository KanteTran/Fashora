package external

import (
	_ "io"
)

//func CallTryOnAPI(personImageURL, clothImageURL, maskURL string) models.Response {
//	body := &bytes.Buffer{}
//	writer := multipart.NewWriter(body)
//
//	fields := map[string]string{
//		"person_image_url": personImageURL,
//		"cloth_image_url":  clothImageURL,
//		"mask_url":         maskURL,
//		"access_token":     RefreshTokenGcp(), // Ensure RefreshTokenGcp() provides the access token
//	}
//
//	for key, value := range fields {
//		if err := writer.WriteField(key, value); err != nil {
//			return models.CreateErrorResponse(
//				http.StatusInternalServerError,
//				fmt.Sprintf("failed to add field '%s': %v", key, err))
//		}
//	}
//
//	if err := writer.Close(); err != nil {
//		return models.CreateErrorResponse(
//			http.StatusInternalServerError, fmt.Sprintf("failed to close writer: %v", err))
//	}
//
//	req, err := http.NewRequest("POST", config.AppConfig.Model.GenAPI, body)
//	if err != nil {
//		return models.CreateErrorResponse(
//			http.StatusInternalServerError,
//			fmt.Sprintf("failed to create request: %v", err))
//	}
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return models.CreateErrorResponse(
//			http.StatusBadGateway,
//			fmt.Sprintf("failed to send request: %v", err))
//	}
//	defer resp.Body.Close()
//
//	bodyNew, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return models.CreateErrorResponse(
//			http.StatusInternalServerError,
//			fmt.Sprintf("failed to read response body: %v", err))
//	}
//
//	var responseJSON map[string]interface{}
//	err = json.Unmarshal(bodyNew, &responseJSON)
//	if err != nil {
//		return models.CreateErrorResponse(
//			http.StatusInternalServerError,
//			fmt.Sprintf("failed to decode JSON: %v", err))
//	}
//
//	resultURL, ok := responseJSON["result_url"].(string)
//	if !ok {
//		return models.CreateErrorResponse(
//			http.StatusInternalServerError,
//			"result_url not found in response or is not a string")
//	}
//
//	return models.Response{
//		Success: true,
//		Status:  http.StatusOK,
//		Message: "Successfully fetched result URL",
//		Data: map[string]string{
//			"result_url": resultURL,
//		},
//	}
//}
//
//func RefreshTokenGcp() string {
//	credentialsFilePath := config.AppConfig.GCS.KeyFile
//
//	data, err := os.ReadFile(credentialsFilePath)
//	if err != nil {
//		log.Fatalf("Failed to read credentials file: %v", err)
//	}
//
//	config, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/cloud-platform")
//	if err != nil {
//		log.Fatalf("Failed to parse credentials file: %v", err)
//	}
//
//	tokenSource := config.TokenSource(context.Background())
//
//	token, err := tokenSource.Token()
//	if err != nil {
//		log.Fatalf("Failed to retrieve token: %v", err)
//	}
//
//	return token.AccessToken
//}
//
//func HomePage(c *gin.Context) {
//	var stores []models.Stores
//
//	if err := models.DB.Find(&stores).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch stores"})
//		return
//	}
//
//	c.HTML(http.StatusOK, "home.html", gin.H{
//		"stores": stores,
//	})
//}
//
//func CreateStorePage(c *gin.Context) {
//	c.HTML(http.StatusOK, "create_store.html", nil)
//}
