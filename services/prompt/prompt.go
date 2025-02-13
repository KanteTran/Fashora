package prompt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"fashora-backend/logger"
)

// Prompt Define the structure for each prompt
type Prompt struct {
	Prompt    string                 `json:"prompt"`
	Criteria  []Criterion            `json:"criteria"`
	Context   map[string]interface{} `json:"context"`
	Responses []Response             `json:"responses"`
}

type Criterion struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Response struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

// PromptLoader handles loading and accessing prompts from a JSON file
type PromptLoader struct {
	prompts map[string]Prompt
}

// NewPromptLoader creates a new instance of PromptLoader
func NewPromptLoader(filePath string) (*PromptLoader, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		logger.Errorf("Error opening file: %s", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Errorf("Error closing file: %v", err)
		}
	}(file)

	// Read file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Errorf("Error reading file: %v", err)
		return nil, err
	}

	// Parse JSON into a map
	var data map[string]Prompt
	err = json.Unmarshal(content, &data)
	if err != nil {
		logger.Errorf("Error parsing JSON: %v", err)
		return nil, err
	}

	// Return a new PromptLoader
	return &PromptLoader{prompts: data}, nil
}

// GetPrompt retrieves a specific prompt by name
func (pl *PromptLoader) GetPrompt(name string) (Prompt, error) {
	prompt, exists := pl.prompts[name]
	if !exists {
		return Prompt{}, errors.New("prompt not found")
	}
	return prompt, nil
}

func ConvertPromptToString(p Prompt) string {
	var sb strings.Builder

	sb.WriteString(p.Prompt)
	for key, value := range p.Context {
		sb.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
	}
	for _, criterion := range p.Criteria {
		sb.WriteString("- " + criterion.Description + "\n")
	}
	for _, response := range p.Responses {
		sb.WriteString(response.Description + "\n")
	}

	return sb.String()
}

func ConvertPromptToString_Recommend(p Prompt,
	userProfile string, useCase string) string {
	var sb strings.Builder

	sb.WriteString(p.Prompt)
	for key, value := range p.Context {
		sb.WriteString(fmt.Sprintf("- %s: %v\n", key, value))
	}
	sb.WriteString(fmt.Sprintf("- %s: %s\n", "Thông tin user bao gồm như sau:\n", userProfile))
	//sb.WriteString(fmt.Sprintf("- %s: %v\n", "Sinh năm", Birthday))
	//sb.WriteString(fmt.Sprintf("- %s: %v\n", "Chiều cao(cm)", Height))
	//sb.WriteString(fmt.Sprintf("- %s: %v\n", "Cân nặng(kg)", Weight))
	//sb.WriteString(fmt.Sprintf("- %s: %v\n", "Màu da", SkinTone))
	//sb.WriteString(fmt.Sprintf("- %s: %v\n", "Giới tính(0: male, 1: female, 2: other)", Gender))
	sb.WriteString(fmt.Sprintf("- %s: %v\n", "Mục đích sử dụng", useCase))

	for _, criterion := range p.Criteria {
		sb.WriteString("- " + criterion.Description + "\n")
	}
	for _, response := range p.Responses {
		sb.WriteString(response.Description + "\n")
	}

	return sb.String()
}
