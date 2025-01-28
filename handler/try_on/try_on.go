package try_on

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"fashora-backend/config"
	"fashora-backend/logger"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fashora-backend/utils"
)

func UploadImages(c *gin.Context) {
	ctx := context.Background()
	logger.Infof("Start - Request: %s, Time: %v\n", c.Request.URL.Path, time.Now())

	// Create GCS client
	client, err := CreateGCSClient(ctx, external.RefreshTokenGcp())
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create GCS client")
		return
	}
	defer client.Close()

	// Image configurations
	images := []models.Image{
		{
			FormKey:    "people",
			BucketName: config.AppConfig.GCS.FolderPeople,
		},
		{
			FormKey:    "mask",
			BucketName: config.AppConfig.GCS.FolderMask,
		},
		{
			FormKey:    "clothes",
			BucketName: config.AppConfig.GCS.FolderClothes,
		},
	}

	// Step 1: Read all files from the request
	files, err := readFilesFromRequest(c, images)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Step 2: Upload files to GCS
	imageURLs, err := uploadImagesToGCS(ctx, client, files, images)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Step 3: Call external TryOn API
	external.CallTryOnAPI(c, imageURLs["people"], imageURLs["clothes"], imageURLs["mask"])

}
