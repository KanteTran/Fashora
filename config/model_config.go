package config

type ModelConfig struct {
	GenAPI  string
	SEGMENT string
}

func loadModelConfig() ModelConfig {
	return ModelConfig{
		GenAPI:  GetEnv("MODEL_GEN_API", ""),
		SEGMENT: GetEnv("MODEL_SEGMENT", ""),
	}
}
