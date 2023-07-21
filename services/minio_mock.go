package service

import (
	"io"
)

type MinioMockService struct {
}

func NewMinioMockService() (*MinioMockService, error) {

	return &MinioMockService{}, nil
}
func (m *MinioMockService) UploadFileToMinio(objName string, file io.Reader, objSize int64, contentType string) error {

	return nil
}
