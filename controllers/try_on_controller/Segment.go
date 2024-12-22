package try_on_controller

import (
	"bytes"
	"fashora-backend/config"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
)

func Segment(c *gin.Context) {
	// Retrieve the uploaded file from the form-data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to retrieve file: "+err.Error())
		return
	}
	defer file.Close()

	// Create a buffer to store the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create the form file field for "file"
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create form file: "+err.Error())
		return
	}

	// Copy the file content into the multipart part
	_, err = io.Copy(part, file)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to copy file content: "+err.Error())
		return
	}

	// Close the writer to finalize the multipart body
	err = writer.Close()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to finalize form data: "+err.Error())
		return
	}

	// Make the HTTP POST request to the external API
	req, err := http.NewRequest("POST", config.AppConfig.Model.SEGMENT, body)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create HTTP request: "+err.Error())
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadGateway, "Failed to send request to external API: "+err.Error())
		return
	}
	defer resp.Body.Close()

	// Read the response from the API
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to read response body: "+err.Error())
		return
	}

	// Return the API response to the client
	if resp.StatusCode != http.StatusOK {
		utils.SendErrorResponse(c, resp.StatusCode, string(respBody))
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Get segment completely", string(respBody))
}
