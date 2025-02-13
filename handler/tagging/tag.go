package tagging

import (
	"fashora-backend/config"
	"fashora-backend/handler/scoring"
	"fashora-backend/logger"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fashora-backend/services/prompt"

	"mime/multipart"
	"strconv"
	"strings"
)

//func TagImage(c *gin.Context) {
//	fileHeader, err := c.FormFile("image")
//	if err != nil {
//		utils.SendErrorResponse(c, http.StatusBadRequest, "Could not get image")
//		return
//	}
//	rawJSON := TagClothes(fileHeader)
//	l
//	if err != nil {
//		log.Fatalf("Error when parse JSON: %v", err)
//	}
//
//	utils.SendSuccessResponse(c, http.StatusOK, "Evaluated complete")
//
//	logger.Info("Image uploaded successfully")
//}

func stringToInt64List(input string) ([]int64, error) {
	// Remove brackets [] if present
	input = strings.Trim(input, "[]")

	// Split the string into parts
	parts := strings.Split(input, ", ")

	var result []int64
	for _, part := range parts {
		// Parse the string into int64
		num, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}
func TagClothes(fileHeader *multipart.FileHeader) []int64 {
	imgData, imgFormat, _ := scoring.PrepareImage(fileHeader)
	GeminiApp := external.InitGemini(config.AppConfig.Model.GeminiAPI)
	logger.Info(config.AppConfig.Prompt.TagClothes)
	tagClothesPrompt, _ := models.PromptLoader.GetPrompt(config.AppConfig.Prompt.TagClothes)
	rawJSON, _ := GeminiApp.GeminiFashionScore(imgFormat, imgData, prompt.ConvertPromptToString(tagClothesPrompt))
	output, _ := stringToInt64List(rawJSON)
	return output
}
