package config

type VersionConfig struct {
	MinimalVersion string
	LatestVersion  string
}

func loadVersionConfig() VersionConfig {
	return VersionConfig{
		MinimalVersion: GetEnv("MINIMAL_VERSION", "1.0.0"),
		LatestVersion:  GetEnv("LATEST_VERSION", "1.0.0"),
	}
}
