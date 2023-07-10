package config

import "os"

type MinioConfig struct {
	AccessKey string
	SecretKey string
	Endpoint string
	ImageBucket string
}

func NewMinioConfig() MinioConfig {
	return MinioConfig{
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Endpoint: os.Getenv("MINIO_IMAGE_ENDPOINT"),
		ImageBucket: os.Getenv("MINIO_IMAGE_BUCKET_NAME"),
	}
}
