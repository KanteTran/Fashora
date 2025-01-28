package config

type FireBaseConfig struct {
	FileKey string
}

func loadFireBaseConfig() FireBaseConfig {
	return FireBaseConfig{
		FileKey: GetEnv("FIREBASE_FILE_CONFIG", ""),
	}
}
