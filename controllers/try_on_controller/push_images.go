package try_on_controller

import (
	"cloud.google.com/go/storage"
	"context"
	"fashora-backend/config"
	"fashora-backend/services/external"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"time"
)

func UploadImages(c *gin.Context) {
	// List of form keys and corresponding bucket names
	images := []struct {
		formKey    string
		bucketName string
	}{
		{"people", config.AppConfig.GscFolderPeople},
		{"clothes", config.AppConfig.GscFolderClothes},
		{"mask", config.AppConfig.GscFolderMask},
	}

	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: external.RefreshTokenGcp(),
	})
	client, err := storage.NewClient(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GCS client"})
		return
	}
	defer client.Close()

	imageURLs := make(map[string]string)

	for _, image := range images {
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

		imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.AppConfig.GscBucketName, objectName)
		imageURLs[image.formKey] = imageURL
	}

	apiResponse := external.CallTryOnAPI(imageURLs["people"], imageURLs["clothes"], imageURLs["mask"])

	c.JSON(apiResponse.Status, apiResponse)

}
