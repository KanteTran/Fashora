package config

type JWTConfig struct {
	Secret          string
	ExpirationHours string
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:          GetEnv("JWT_SECRET", "default_jwt_secret"),
		ExpirationHours: GetEnv("JWT_EXPIRATION_HOURS", "72"),
	}
}
