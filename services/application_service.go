package service

import "go-ekyc/repository"

type ApplicationService struct {
	CustomerService *CustomerService
	PlansService *PlansService
	ImageService *ImageService
	MinioService *MinioService
	FaceMatchScoreService *FaceMatchScoreService
	OCRService *OCRService
}

func NewApplicationService(applicationRepository *repository.ApplicationRepository) *ApplicationService {
	
	return &ApplicationService{
		CustomerService: newCustomerService(applicationRepository.CustomerRepository),
		PlansService: newPlansService(applicationRepository.PlansRepository),
		ImageService: newImageService(applicationRepository.ImageRepository),
		FaceMatchScoreService: newFaceMatchScoreService(applicationRepository.FaceMatchScoreRepository),
		OCRService: newOCRService(applicationRepository.OCRRepository),
	}
}
