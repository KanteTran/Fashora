package try_on_controller

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"

	"fashora-backend/config"
	"fashora-backend/models"
)

func readFilesFromRequest(c *gin.Context, images []models.Image) (map[string]*multipart.FileHeader, error) {
	files := make(map[string]*multipart.FileHeader)
	for _, image := range images {
		file, err := c.FormFile(image.FormKey)
		if err != nil {
			return nil, fmt.Errorf("No file uploaded for %s", image.FormKey)
		}
		files[image.FormKey] = file
	}
	return files, nil
}

func uploadToGCS(ctx context.Context,
	client *storage.Client,
	fileContent io.Reader,
	formKey,
	fileName string) (string, error) {
	objectName := fmt.Sprintf("%s/%d_%s", formKey, time.Now().Unix(), fileName)
	bucket := client.Bucket(config.AppConfig.GCS.BucketName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)

	// Write content to GCS
	if _, err := io.Copy(writer, fileContent); err != nil {
		return "", fmt.Errorf("Failed to upload file %s: %v", formKey, err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("Failed to finalize upload for %s: %v", formKey, err)
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.AppConfig.GCS.BucketName, objectName)
	return imageURL, nil
}

// CreateGCSClient initializes a GCS client using a token
func CreateGCSClient(ctx context.Context, accessToken string) (*storage.Client, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})
	client, err := storage.NewClient(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("Failed to create GCS client: %v", err)
	}
	return client, nil
}
