package config

type ModelConfig struct {
	GenAPI string
}

func loadModelConfig() ModelConfig {
	return ModelConfig{
		GenAPI: GetEnv("MODEL_GEN_API", ""),
	}
}
