package external

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"

	"fashora-backend/config"
	"fashora-backend/logger"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func CreateFoldersIfNotExists(bucketName string, folderPath string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.AppConfig.GCS.KeyFile))
	if err != nil {
		return fmt.Errorf("failed to create gsc keyfile: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	folders := strings.Split(folderPath, "/")

	currentPath := ""
	for _, folder := range folders {
		if folder == "" {
			continue
		}

		currentPath += folder + "/"

		query := &storage.Query{Prefix: currentPath, Delimiter: "/"}
		it := bucket.Objects(ctx, query)
		_, err := it.Next()

		if err == nil {
			log.Printf("Folder '%s' already exists in bucket '%s'\n", currentPath, bucketName)
			continue
		} else if !errors.Is(err, iterator.Done) {
			return fmt.Errorf("error checking if folder '%s' exists: %v", currentPath, err)
		}

		obj := bucket.Object(currentPath)
		w := obj.NewWriter(ctx)
		if _, err := w.Write([]byte{}); err != nil {
			return fmt.Errorf("failed to create folder '%s': %v", currentPath, err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close writer for folder '%s': %v", currentPath, err)
		}

		logger.Infof("Folder '%s' created in bucket '%s'\n", currentPath, bucketName)
	}

	return nil
}
func UploadImageToGCS(fileURL string, file *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	parts := strings.SplitN(fileURL, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid file URL: must contain bucket name and object path")
	}

	bucketName := parts[0]
	objectPath := parts[1]

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.AppConfig.GCS.KeyFile))
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectPath)
	writer := obj.NewWriter(ctx)
	writer.ContentType = file.Header.Get("Content-Type")

	if _, err := io.Copy(writer, src); err != nil {
		return "", fmt.Errorf("failed to upload file to GCS: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
	logger.Infof("File uploaded successfully: %s\n", publicURL)
	return publicURL, nil
}
