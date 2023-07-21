package service

import (
	"go-ekyc/config"
	"go-ekyc/repository"
)

type ApplicationService struct {
	CustomerService       *CustomerService
	PlansService          *PlansService
	ImageService          *ImageService
	MinioService          IMinioService
	FaceMatchScoreService *FaceMatchScoreService
	OCRService            *OCRService
}

func NewApplicationService(applicationRepository *repository.ApplicationRepository) (*ApplicationService, error) {
	minioConfig := config.NewMinioConfig()
	minioService, err := NewMinioService(minioConfig)

	if err != nil {
		return nil, err
	}
	return &ApplicationService{
		CustomerService:       newCustomerService(applicationRepository.CustomerRepository, applicationRepository.PlansRepository, applicationRepository.ImageRepository, applicationRepository.FaceMatchScoreRepository, applicationRepository.OCRRepository, applicationRepository.DailyReportsRepository, applicationRepository.RedisRepository),
		PlansService:          newPlansService(applicationRepository.PlansRepository),
		ImageService:          newImageService(applicationRepository.ImageRepository, applicationRepository.PlansRepository, applicationRepository.FaceMatchScoreRepository, applicationRepository.OCRRepository, newOCRService(applicationRepository.OCRRepository), minioService),
		FaceMatchScoreService: newFaceMatchScoreService(applicationRepository.FaceMatchScoreRepository),
		OCRService:            newOCRService(applicationRepository.OCRRepository),
		MinioService: minioService,
	}, nil
}
