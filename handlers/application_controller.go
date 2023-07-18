package handlers

import service "go-ekyc/services"

type ApplicationHandler struct {
	CustomerHandler *CustomerHandlers
	ImageHandler    *ImageHandlers
}

func NewApplicationHandler(applicationService *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		ImageHandler:    newImageHandler(applicationService.CustomerService, applicationService.ImageService),
		CustomerHandler: newCustomerHandler(applicationService.CustomerService, applicationService.PlansService),
	}
}
