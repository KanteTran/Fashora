package try_on_controller

import (
	"cloud.google.com/go/storage"
	"context"
	"fashora-backend/config"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/googleapis/enterprise-certificate-proxy/client"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"sync"
	"time"
)

func UploadImagesV2(c *gin.Context) {
	ctx := context.Background()
	fmt.Printf("start - Request: %s, Time: %v\n", c.Request.URL.Path, time.Now())

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: external.RefreshTokenGcp(),
	})
	client, err := storage.NewClient(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GCS client"})
		return
	}
	defer client.Close()
	images := []models.Image{
		{"people", config.AppConfig.GscFolderPeople},
		{"clothes", config.AppConfig.GscFolderClothes},
		{"mask", config.AppConfig.GscFolderMask},
	}
	fmt.Println(images)
	imageURLs := make(map[string]string)
	errorChannel := make(chan error, len(images)) // Channel để thu thập lỗi
	results := make(chan struct {
		formKey  string
		imageURL string
		err      error
	}, len(images)) // Channel để nhận kết quả upload
	var wg sync.WaitGroup

	for _, image := range images {
		fmt.Println(image.BucketName)
		fmt.Printf(image.FormKey)
		wg.Add(1)
		go func(image models.Image) {
			defer wg.Done()

			// Nhận file từ form
			file, err := c.FormFile(image.FormKey)
			if err != nil {
				results <- struct {
					formKey  string
					imageURL string
					err      error
				}{image.FormKey, "",
					fmt.Errorf("no file uploaded for %s", image.FormKey)}
				return
			}

			// Mở nội dung file
			fileContent, err := file.Open()
			if err != nil {
				results <- struct {
					formKey  string
					imageURL string
					err      error
				}{image.FormKey, "",
					fmt.Errorf("Unable to open file for %s: %v", image.FormKey, err)}
				return
			}
			defer fileContent.Close()
			//fmt.Println("oke den day roi ne")
			// Tạo objectName cho file trong GCS
			objectName := fmt.Sprintf("%s/%d_%s", image.BucketName, time.Now().Unix(), file.Filename)
			bucket := client.Bucket(config.AppConfig.GscBucketName)
			object := bucket.Object(objectName)
			writer := object.NewWriter(ctx)

			// Upload file lên GCS
			if _, err := io.Copy(writer, fileContent); err != nil {
				results <- struct {
					formKey  string
					imageURL string
					err      error
				}{image.FormKey, "",
					fmt.Errorf("Failed to upload image %s to GCS: %v", image.FormKey, err)}
				return
			}

			// Đóng writer
			if err := writer.Close(); err != nil {
				results <- struct {
					formKey  string
					imageURL string
					err      error
				}{image.FormKey, "",
					fmt.Errorf("Failed to finalize upload for %s: %v", image.FormKey, err)}
				return
			}

			// Thành công
			imageURL := fmt.Sprintf(
				"https://storage.googleapis.com/%s/%s",
				config.AppConfig.GscBucketName,
				objectName)
			results <- struct {
				formKey  string
				imageURL string
				err      error
			}{image.FormKey, imageURL, nil}
		}(image)
	}

	// Chờ tất cả các goroutine hoàn thành
	go func() {
		wg.Wait()
		close(results)
	}()

	// Xử lý kết quả
	for res := range results {
		if res.err != nil {
			errorChannel <- res.err
		} else {
			imageURLs[res.formKey] = res.imageURL
		}
	}

	if len(errorChannel) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": <-errorChannel})
		return
	}

	c.JSON(http.StatusOK, gin.H{"imageURLs": imageURLs})
}
