package service

import "go-ekyc/repository"

type ApplicationService struct {
	CustomerService       *CustomerService
	PlansService          *PlansService
	ImageService          *ImageService
	MinioService          *MinioService
	FaceMatchScoreService *FaceMatchScoreService
	OCRService            *OCRService
}

func NewApplicationService(applicationRepository *repository.ApplicationRepository) *ApplicationService {

	return &ApplicationService{
		CustomerService:       newCustomerService(applicationRepository.CustomerRepository, applicationRepository.PlansRepository, applicationRepository.ImageRepository, applicationRepository.FaceMatchScoreRepository, applicationRepository.OCRRepository,applicationRepository.DailyReportsRepository,applicationRepository.RedisRepository),
		PlansService:          newPlansService(applicationRepository.PlansRepository),
		ImageService:          newImageService(applicationRepository.ImageRepository, applicationRepository.PlansRepository, applicationRepository.FaceMatchScoreRepository, applicationRepository.OCRRepository, newOCRService(applicationRepository.OCRRepository)),
		FaceMatchScoreService: newFaceMatchScoreService(applicationRepository.FaceMatchScoreRepository),
		OCRService:            newOCRService(applicationRepository.OCRRepository),
	}
}
