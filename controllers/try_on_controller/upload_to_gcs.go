package try_on_controller

import (
	"context"
	"fmt"
	"mime/multipart"
	"sync"

	"cloud.google.com/go/storage"

	"fashora-backend/models"
)

func uploadImagesToGCS(ctx context.Context,
	client *storage.Client,
	files map[string]*multipart.FileHeader,
	images []models.Image) (map[string]string, error) {
	var wg sync.WaitGroup
	results := make(chan struct {
		formKey  string
		imageURL string
		err      error
	}, len(images))

	// Concurrently upload files
	for _, image := range images {
		wg.Add(1)
		go func(image models.Image, file *multipart.FileHeader) {
			defer wg.Done()

			// Open file content
			fileContent, err := file.Open()
			if err != nil {
				results <- struct {
					formKey  string
					imageURL string
					err      error
				}{image.FormKey, "", fmt.Errorf("Unable to open file for %s: %v", image.FormKey, err)}
				return
			}
			defer fileContent.Close()

			// Upload to GCS
			imageURL, err := uploadToGCS(ctx, client, fileContent, image.FormKey, file.Filename)
			results <- struct {
				formKey  string
				imageURL string
				err      error
			}{image.FormKey, imageURL, err}
		}(image, files[image.FormKey])
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	imageURLs := make(map[string]string)
	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		imageURLs[res.formKey] = res.imageURL
	}
	return imageURLs, nil
}
