package scoring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/adrium/goheif"
	"github.com/gin-gonic/gin"
	"golang.org/x/image/webp"
	_ "golang.org/x/image/webp"

	"fashora-backend/config"
	"fashora-backend/logger"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fashora-backend/services/prompt"
	"fashora-backend/utils"
)

type ScoreResponse struct {
	StyleDescription       string `json:"style_description"`
	BodyShapeSkinTone      string `json:"body_shape_skin_tone"`
	ColorHarmonyScore      string `json:"color_harmony_score"`
	QualityScore           string `json:"quality_score"`
	BalanceScore           string `json:"balance_score"`
	StyleMatchScore        string `json:"style_match_score"`
	BodyShapeFitScore      string `json:"body_shape_fit_score"`
	PracticalityScore      string `json:"practicality_score"`
	ComfortScore           string `json:"comfort_score"`
	SkinHairToneMatchScore string `json:"skin_hair_tone_match_score"`
	Conclusion             string `json:"conclusion"`
	SuggestedImprovements  string `json:"suggested_improvements"`
}

// Detect image format
func detectFormat(filename string) string {
	logger.Info("detecting format")
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
func convertHEICToPNG(file *multipart.FileHeader) ([]byte, error) {
	// Read the HEIC file into memory

	fileHeader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileHeader.Close()
	imgData, err := io.ReadAll(fileHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to read HEIC file: %v", err)
	}

	// Create an io.Reader from the byte slice
	reader := bytes.NewReader(imgData)

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
func convertWebPToPNG(file *multipart.FileHeader) ([]byte, error) {

	fileHeader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func(fileHeader multipart.File) {
		err := fileHeader.Close()
		if err != nil {
		}
	}(fileHeader)
	img, _ := io.ReadAll(fileHeader)

	// Create an io.Reader from the byte slice
	reader := bytes.NewReader(img)
	decoded, err := webp.Decode(reader)
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
func PrepareImage(fileHeader *multipart.FileHeader) ([]byte, string, error) {
	format := detectFormat(fileHeader.Filename)

	switch format {
	case "jpeg", "png":
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return nil, "", err
		}
		defer file.Close()
		imgData, err := io.ReadAll(file)
		return imgData, format, err

	case "heic":
		imgData, err := convertHEICToPNG(fileHeader)
		return imgData, "png", err

	case "webp":
		imgData, err := convertWebPToPNG(fileHeader)
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

	imgData, imgFormat, err := PrepareImage(fileHeader)
	logger.Infof("Image file read successfully, si	ze: %d bytes", len(imgData))

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
		logger.Errorf("Error when parse JSON: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, rawJSON)
	} else {
		utils.SendSuccessResponse(c, http.StatusOK, "Evaluated complete", evaluation)
	}
}
