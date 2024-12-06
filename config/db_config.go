package config

type DBConfig struct {
	User     string
	Password string
	DB       string
	Host     string
	Port     string
}

func loadDBConfig() DBConfig {
	return DBConfig{
		User:     GetEnv("POSTGRES_USER", "postgres"),
		Password: GetEnv("POSTGRES_PASSWORD", "password"),
		DB:       GetEnv("POSTGRES_DB", "postgres_db"),
		Host:     GetEnv("POSTGRES_HOST", "localhost"),
		Port:     GetEnv("POSTGRES_PORT", "5432"),
	}
}
