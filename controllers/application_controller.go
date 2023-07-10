package controllers

import service "go-ekyc/services"

type ApplicationController struct {
	CustomerController *CustomerControllers
	ImageController    *ImageControllers
}

func NewApplicationController(applicationService *service.ApplicationService) *ApplicationController {
	return &ApplicationController{
		ImageController:    newImageController(applicationService.CustomerService, applicationService.PlansService, applicationService.ImageService),
		CustomerController: newCustomerController(applicationService.CustomerService, applicationService.PlansService),
	}
}
