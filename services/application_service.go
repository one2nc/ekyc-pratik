package service

import (
	"go-ekyc/repository"
)

type ApplicationService struct {
	CustomerService       *CustomerService
	ImageService          *ImageService
	MinioService          IMinioService
	OCRService            *OCRService
}

func NewApplicationService(applicationRepository *repository.ApplicationRepository, minioService IMinioService) *ApplicationService {

	return &ApplicationService{
		CustomerService:       newCustomerService(applicationRepository.CustomerRepository, applicationRepository.PlansRepository, applicationRepository.ImageRepository, applicationRepository.FaceMatchScoreRepository, applicationRepository.OCRRepository, applicationRepository.DailyReportsRepository, applicationRepository.RedisRepository),
		ImageService:          newImageService(applicationRepository.ImageRepository, applicationRepository.PlansRepository, applicationRepository.FaceMatchScoreRepository, applicationRepository.OCRRepository, newOCRService(applicationRepository.OCRRepository), minioService),
		OCRService:            newOCRService(applicationRepository.OCRRepository),
		MinioService:          minioService,
	}
}
