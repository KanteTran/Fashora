package config

type PromptConfig struct {
	PromptFileName   string
	OutfitEvalPrompt string
}

func loadPromptConfig() PromptConfig {
	return PromptConfig{
		PromptFileName:   GetEnv("FILE_PROMPT", ""),
		OutfitEvalPrompt: GetEnv("OUTFIT_EVALUATION", ""),
	}
}
