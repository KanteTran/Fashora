package config

import (
	"fashora-backend/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres DbPostGreSQLConfig
	JWT      JWTConfig
	GCS      GCSConfig
	Server   ServerConfig
	Model    ModelConfig
	FireBase FireBaseConfig
	Version  VersionConfig
	Prompt   PromptConfig
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logger.Info("No .env file found, reading from environment variables")
	}

	AppConfig = Config{
		Postgres: loadDBConfig(),
		JWT:      loadJWTConfig(),
		GCS:      loadGCSConfig(),
		Server:   loadServerConfig(),
		Model:    loadModelConfig(),
		FireBase: loadFireBaseConfig(),
		Version:  loadVersionConfig(),
		Prompt:   loadPromptConfig(),
	}
}
