package config

import (
	"os"

	"fashora-backend/logger"
)

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	logger.Infof("Environment variable %s not found. Using default value: %s", key, defaultValue)
	return defaultValue
}
