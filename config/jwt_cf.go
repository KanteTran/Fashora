package config

type JWTConfig struct {
	Secret                 string
	AccessTokenExpiration  string
	RefreshTokenExpiration string
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:                 GetEnv("JWT_SECRET", "default_jwt_secret"),
		AccessTokenExpiration:  GetEnv("JWT_ACCESS_TOKEN_EXPIRATION_MINITUNES", "72"),
		RefreshTokenExpiration: GetEnv("JWT_REFRESH_TOKEN_EXPIRATION_HOURS", "72"),
	}
}
