package config

type GCSConfig struct {
	BucketName     string
	FolderPeople   string
	FolderMask     string
	KeyFile        string
	StoreImagePath string
}

func loadGCSConfig() GCSConfig {
	return GCSConfig{
		BucketName:     GetEnv("GSC_BUCKET_NAME", ""),
		FolderPeople:   GetEnv("GSC_FOLDER_PEOPLE", ""),
		FolderMask:     GetEnv("GSC_FOLDER_MASK", ""),
		KeyFile:        GetEnv("GSC_KEY_FILE", ""),
		StoreImagePath: GetEnv("GSC_STORES_IMAGE", ""),
	}
}
