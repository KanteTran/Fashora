package external

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fashora-backend/config"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"strings"
)

func CreateFoldersIfNotExists(bucketName string, folderPath string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.AppConfig.GscKeyFile))
	if err != nil {
		return fmt.Errorf("failed to create gsc keyfile: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	// Split the folder path into levels
	folders := strings.Split(folderPath, "/")

	// Iterate through the levels and create each one if it doesn't exist
	currentPath := ""
	for _, folder := range folders {
		if folder == "" {
			continue
		}

		// Update the current path
		currentPath += folder + "/"

		// Check if the folder exists
		query := &storage.Query{Prefix: currentPath, Delimiter: "/"}
		it := bucket.Objects(ctx, query)
		_, err := it.Next()

		if err == nil {
			fmt.Printf("Folder '%s' already exists in bucket '%s'\n", currentPath, bucketName)
			continue
		} else if !errors.Is(err, iterator.Done) {
			return fmt.Errorf("error checking if folder '%s' exists: %v", currentPath, err)
		}

		// Create the folder
		obj := bucket.Object(currentPath)
		w := obj.NewWriter(ctx)
		if _, err := w.Write([]byte{}); err != nil {
			return fmt.Errorf("failed to create folder '%s': %v", currentPath, err)
		}
		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close writer for folder '%s': %v", currentPath, err)
		}

		fmt.Printf("Folder '%s' created in bucket '%s'\n", currentPath, bucketName)
	}

	return nil
}
func UploadImageToGCS(fileURL string, file *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	// Parse the bucket name and object path from fileURL
	parts := strings.SplitN(fileURL, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid file URL: must contain bucket name and object path")
	}

	bucketName := parts[0]
	objectPath := parts[1]

	// Initialize GCS client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.AppConfig.GscKeyFile))
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	print("oke r ne" +
		"")
	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Get the bucket and object
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectPath)
	writer := obj.NewWriter(ctx)
	writer.ContentType = file.Header.Get("Content-Type")
	//writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}} // Make the file public

	// Copy the file to GCS
	if _, err := io.Copy(writer, src); err != nil {
		return "", fmt.Errorf("failed to upload file to GCS: %v", err)
	}

	// Close the writer
	print("ddkdmk")
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Generate and return the public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
	fmt.Printf("File uploaded successfully: %s\n", publicURL)
	return publicURL, nil
}
