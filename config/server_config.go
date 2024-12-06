package config

type ServerConfig struct {
	Host string
	Port string
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Host: GetEnv("HOST_SERVER", "localhost"),
		Port: GetEnv("PORT_SERVER", "8080"),
	}
}
