package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PostgresUser       string
	PostgresPassword   string
	PostgresDB         string
	PostgresHost       string
	PostgresPort       string
	JWTSecret          string
	JwtExpirationHours string
}

var AppConfig Config

func LoadConfig() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	// Populate Config struct with environment variables
	AppConfig = Config{
		PostgresUser:       getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword:   getEnv("POSTGRES_PASSWORD", "password"),
		PostgresDB:         getEnv("POSTGRES_DB", "postgres_db"),
		PostgresHost:       getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:       getEnv("POSTGRES_PORT", "5432"),
		JWTSecret:          getEnv("JWT_SECRET", "default_jwt_secret"),
		JwtExpirationHours: getEnv("JWT_EXPIRATION_HOURS", "72"),
	}
}

// getEnv reads an environment variable or returns a default value if not found
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
