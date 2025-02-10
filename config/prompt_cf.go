package config

type PromptConfig struct {
	PromptFileName   string
	OutfitEvalPrompt string
	TagClothes       string
	RecommendTags    string
}

func loadPromptConfig() PromptConfig {
	return PromptConfig{
		PromptFileName:   GetEnv("FILE_PROMPT", ""),
		OutfitEvalPrompt: GetEnv("OUTFIT_EVALUATION", ""),
		TagClothes:       GetEnv("TAG_CLOTHES", ""),
		RecommendTags:    GetEnv("RECOMMEND_PROMPT", ""),
	}
}
