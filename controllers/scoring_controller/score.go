package scoring_controller

import (
	"fashora-backend/config"
	"fashora-backend/services/external"
	"fashora-backend/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

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
	//TODO: We need to put prompt some where can update immediately
	// Define the fashion scoring prompt
	prompt := `
    Please analyze the uploaded outfit image based on the following criteria:
    1. Rate the style: Is it casual, formal, or sporty? Score it out of 10 and explain your reasoning.
    2. Evaluate the color combination: Does it look harmonious? Score it out of 10 and briefly explain why.
    3. Assess the occasion suitability: Is it appropriate for a casual outing, formal meeting, or a party? Score it out of 10 and provide a short explanation.
    4. Trend analysis: How trendy or timeless does this outfit look? Score it out of 10 and explain your reasoning.
    5. Accessory matching: Do the accessories match the outfit? Score it out of 10 and describe why they do or do not match.
    `
	GeminiApp := external.InitGemini(config.AppConfig.Model.GeminiAPI)
	score, err := GeminiApp.GeminiFashionScore(imgData, prompt)
	if err != nil {
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Evaluated complete", score)

}
