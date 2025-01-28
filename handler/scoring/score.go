package scoring

import (
	"fashora-backend/config"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fashora-backend/services/prompt"
	"fashora-backend/utils"

	"github.com/gin-gonic/gin"

	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type ScoreResponse struct {
	StyleDescription       string `json:"style_description"`
	BodyShapeSkinTone      string `json:"body_shape_skin_tone"`
	ColorHarmonyScore      int    `json:"color_harmony_score"`
	QualityScore           int    `json:"quality_score"`
	BalanceScore           int    `json:"balance_score"`
	StyleMatchScore        int    `json:"style_match_score"`
	BodyShapeFitScore      int    `json:"body_shape_fit_score"`
	PracticalityScore      int    `json:"practicality_score"`
	ComfortScore           int    `json:"comfort_score"`
	SkinHairToneMatchScore int    `json:"skin_hair_tone_match_score"`
	Conclusion             string `json:"conclusion"`
	SuggestedImprovements  string `json:"suggested_improvements"`
}

func ScoreImage(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Could not get image")
		return
	}
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not open image file")
		return
	}
	defer file.Close()

	// Read the file content into a byte slice
	imgData, err := io.ReadAll(file)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not read image file")
		return
	}
	GeminiApp := external.InitGemini(config.AppConfig.Model.GeminiAPI)
	outfitEvalPrompt, _ := models.PromptLoader.GetPrompt(config.AppConfig.Prompt.OutfitEvalPrompt)
	rawJSON, _ := GeminiApp.GeminiFashionScore(imgData, prompt.ConvertPromptToString(outfitEvalPrompt))

	// filter "```json\n" và "```\n" in response from model
	cleanedJSON := strings.TrimPrefix(rawJSON, "```json\n")
	cleanedJSON = strings.TrimSuffix(cleanedJSON, "```\n")

	// Parse JSON to object
	var evaluation ScoreResponse
	err = json.Unmarshal([]byte(cleanedJSON), &evaluation)
	if err != nil {
		log.Fatalf("Lỗi khi parse JSON: %v", err)
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Evaluated complete", evaluation)

}
