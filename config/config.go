package config

import (
	"github.com/joho/godotenv"

	"log"
)

type Config struct {
	Postgres DBConfig
	JWT      JWTConfig
	GCS      GCSConfig
	Server   ServerConfig
	Model    ModelConfig
	FireBase FireBaseConfig
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	AppConfig = Config{
		Postgres: loadDBConfig(),
		JWT:      loadJWTConfig(),
		GCS:      loadGCSConfig(),
		Server:   loadServerConfig(),
		Model:    loadModelConfig(),
		FireBase: loadFireBaseConfig(),
	}
}
