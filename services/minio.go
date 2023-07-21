package service

import (
	"go-ekyc/config"
	"io"

	"github.com/minio/minio-go"
)

type IMinioService interface {
	UploadFileToMinio( objName string, file io.Reader, objSize int64, contentType string) error
}

type MinioService struct {
	minioClient *minio.Client
	minioConfig config.MinioConfig
}

func NewMinioService(config config.MinioConfig) (IMinioService, error) {

	minioClient, err := minio.New(config.Endpoint, config.AccessKey, config.SecretKey, false)
	if err != nil {
		return nil, err
	}

	return &MinioService{
		minioClient: minioClient,
		minioConfig: config,
	}, nil
}
func (m *MinioService) UploadFileToMinio( objName string, file io.Reader, objSize int64, contentType string) error {

	_, err := m.minioClient.PutObject(m.minioConfig.ImageBucket, objName, file, objSize, minio.PutObjectOptions{ContentType: contentType})
	return err
}
