package config

type ModelConfig struct {
	GenAPI          string
	SEGMENT         string
	GeminiAPI       string
	GeminiModelName string
	Prompt          string
}

func loadModelConfig() ModelConfig {
	return ModelConfig{
		GenAPI:          GetEnv("MODEL_GEN_API", ""),
		SEGMENT:         GetEnv("MODEL_SEGMENT", ""),
		GeminiAPI:       GetEnv("GEMINI_API_KEY", ""),
		GeminiModelName: GetEnv("GEMINI_MODEL_NAME", ""),
		Prompt:          GetEnv("MODEL_PROMPT", ""),
	}
}
