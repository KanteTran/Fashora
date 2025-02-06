package scoring

import (
	"bytes"
	"fashora-backend/config"
	"fashora-backend/logger"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fashora-backend/services/prompt"
	"fashora-backend/utils"
	"fmt"
	"golang.org/x/image/webp"
	"image/png"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/adrium/goheif"
	"github.com/gin-gonic/gin"
	_ "golang.org/x/image/webp"

	"encoding/json"
	"io"
	"log"
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

// Detect image format
func detectFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpeg", ".jpg":
		return "jpeg"
	case ".png":
		return "png"
	case ".heic", ".heif":
		return "heic"
	case ".webp":
		return "webp"
	default:
		return "unknown"
	}
}

// Convert HEIC to PNG
func convertHEICToPNG(file multipart.File) ([]byte, error) {
	// Read the HEIC file into memory
	heicData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read HEIC file: %v", err)
	}

	// Create an io.Reader from the byte slice
	reader := bytes.NewReader(heicData)

	// Decode the HEIC image using goheif
	img, err := goheif.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HEIC image: %v", err)
	}

	// Convert the decoded image to PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %v", err)
	}

	// Return the PNG data
	return buf.Bytes(), nil
}

// Convert WebP to PNG
func convertWebPToPNG(img io.Reader) ([]byte, error) {
	decoded, err := webp.Decode(img)
	if err != nil {
		return nil, fmt.Errorf("failed to decode WebP image: %v", err)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %v", err)
	}

	return buf.Bytes(), nil
}

// Prepare image: detect format & convert if needed
func prepareImage(file multipart.File, filename string) ([]byte, string, error) {
	format := detectFormat(filename)

	switch format {
	case "jpeg", "png":
		imgData, err := io.ReadAll(file)
		return imgData, format, err

	case "heic":
		imgData, err := convertHEICToPNG(file)
		return imgData, "png", err

	case "webp":
		imgData, err := convertWebPToPNG(file)
		return imgData, "png", err

	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", format)
	}
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

	imgData, imgFormat, err := prepareImage(file, fileHeader.Filename)
	logger.Info(imgFormat)

	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not read image file")
		return
	}
	GeminiApp := external.InitGemini(config.AppConfig.Model.GeminiAPI)
	outfitEvalPrompt, _ := models.PromptLoader.GetPrompt(config.AppConfig.Prompt.OutfitEvalPrompt)
	rawJSON, _ := GeminiApp.GeminiFashionScore(imgFormat, imgData, prompt.ConvertPromptToString(outfitEvalPrompt))

	// filter "```json\n" và "```\n" in response from model
	cleanedJSON := strings.TrimPrefix(rawJSON, "```json\n")
	cleanedJSON = strings.TrimSuffix(cleanedJSON, "```\n")

	// Parse JSON to object
	var evaluation ScoreResponse
	err = json.Unmarshal([]byte(cleanedJSON), &evaluation)
	if err != nil {
		log.Fatalf("Error when parse JSON: %v", err)
	}

	utils.SendSuccessResponse(c, http.StatusOK, "Evaluated complete", evaluation)

}
