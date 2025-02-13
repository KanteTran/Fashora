package recommend

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"fashora-backend/config"
	"fashora-backend/logger"
	"fashora-backend/models"
	"fashora-backend/services/external"
	"fashora-backend/services/prompt"
	"fashora-backend/utils"
)

type FashionRecommendation struct {
	Advise string `json:"Advise"`
	Bo1    Outfit `json:"Bộ 1"`
	Bo2    Outfit `json:"Bộ 2"`
	Bo3    Outfit `json:"Bộ 3"`
}

type Outfit struct {
	Advise string  `json:"Advise"`
	Items  [][]int `json:"items"`
}

func GenTagRecommend(c *gin.Context) {
	useCase := c.PostForm("use_case")
	userProfile := c.PostForm("user_profile")
	//birthday := c.PostForm("birthday")
	//height := c.PostForm("height")
	//weight := c.PostForm("weight")
	//skinTone := c.PostForm("skin_tone")
	//gender := c.PostForm("gender")
	//user, err := auth_service.GetAuthenticatedUser(c)
	//if err != nil {
	//	utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
	//	return
	//}

	GeminiApp := external.InitGemini(config.AppConfig.Model.GeminiAPI)
	recommendPrompt, _ := models.PromptLoader.GetPrompt(config.AppConfig.Prompt.RecommendTags)
	rawJSON, _ := GeminiApp.GeminiFashionTags(
		prompt.ConvertPromptToString_Recommend(recommendPrompt, userProfile, useCase))

	// filter "```json\n" và "```\n" in response from model
	cleanedJSON := strings.TrimPrefix(rawJSON, "```json\n")
	cleanedJSON = strings.TrimSuffix(cleanedJSON, "```\n")

	logger.Info(cleanedJSON)
	var recommendation FashionRecommendation
	err := json.Unmarshal([]byte(cleanedJSON), &recommendation)
	if err != nil {
		log.Fatalf("Error when parse JSON: %v", err)
		utils.SendErrorResponse(c, http.StatusInternalServerError, rawJSON)
	} else {
		utils.SendSuccessResponse(c, http.StatusOK, "Evaluated complete", recommendation)
	}

}
