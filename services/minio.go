package service

import (
	"go-ekyc/config"
	"io"

	"github.com/minio/minio-go"
)

type MinioService struct {
	minioClient *minio.Client
	MinioConfig config.MinioConfig
}
func NewMinioService(config config.MinioConfig ) (*MinioService, error) {

	
	minioClient, err := minio.New(config.Endpoint, config.AccessKey, config.SecretKey, false)
	if err != nil {
		return nil, err
	}

	return &MinioService{
		minioClient: minioClient,
		MinioConfig: config,
	}, nil
}
func (m *MinioService) UploadFileToMinio(bucketName string, objName string, file io.Reader, objSize int64, contentType string) error {

	_, err := m.minioClient.PutObject(bucketName, objName, file, objSize, minio.PutObjectOptions{ContentType: contentType})
	return err
}


