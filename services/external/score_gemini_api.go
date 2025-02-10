package external

import (
	"context"
	"fmt"

	"fashora-backend/config"
	"fashora-backend/logger"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiApp struct {
	geminiKey string
	ctx       context.Context
	client    *genai.Client
}

// InitGemini Initialize Gemini client
func InitGemini(key string) *GeminiApp {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		logger.Error("Cannot create client Gemini")
	}
	return &GeminiApp{key, ctx, client}
}

// GeminiFashionScore Send an image to Gemini and get fashion scoring
func (app *GeminiApp) GeminiFashionScore(imgFormat string, imgData []byte, prompt string) (string, error) {
	// Set up the generative model with a custom temperature
	model := app.client.GenerativeModel(config.AppConfig.Model.GeminiModelName)
	temp := float32(0.8)
	model.Temperature = &temp

	// Create the input data for the model (image and prompt)
	data := []genai.Part{
		genai.ImageData(imgFormat, imgData), // Assuming the image is in JPEG format
		genai.Text(prompt),
	}

	logger.Info("Processing fashion scoring for the uploaded image...")

	// Generate content based on the input
	resp, err := model.GenerateContent(app.ctx, data...)
	if err != nil {
		logger.Errorf("Error during processing: %s", err)
		return "", err
	}

	// Extract and return the response
	return printResponse(resp), nil
}

// GeminiFashionScore Send an image to Gemini and get fashion scoring
func (app *GeminiApp) GeminiFashionTags(prompt string) (string, error) {
	// Set up the generative model with a custom temperature
	model := app.client.GenerativeModel(config.AppConfig.Model.GeminiModelName)
	temp := float32(0.8)
	model.Temperature = &temp

	// Create the input data for the model (image and prompt)
	data := []genai.Part{
		genai.Text(prompt),
	}

	resp, err := model.GenerateContent(app.ctx, data...)
	if err != nil {
		logger.Errorf("Error during processing: %s", err)
		return "", err
	}

	// Extract and return the response
	return printResponse(resp), nil
}

// Print the response from Gemini
func printResponse(resp *genai.GenerateContentResponse) string {
	var result string
	for _, cand := range resp.Candidates {
		for _, part := range cand.Content.Parts {
			result += fmt.Sprintf("%v", part)
		}
	}
	return result
}
