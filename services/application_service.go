package service

import "go-ekyc/repository"

type ApplicationService struct {
	CustomerService *CustomerService
	PlansService *PlansService
	ImageService *ImageService
	MinioService *MinioService
}

func NewApplicationService(applicationRepository *repository.ApplicationRepository) *ApplicationService {
	
	return &ApplicationService{
		CustomerService: newCustomerService(applicationRepository.CustomerRepository),
		PlansService: newPlansService(applicationRepository.PlansRepository),
		ImageService: newImageService(applicationRepository.ImageRepository),
	}
}
