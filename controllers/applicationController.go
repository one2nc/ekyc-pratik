package controllers

import service "go-ekyc/services"

type ApplicationController struct{
	CustomerController *CustomerControllers
}

func NewApplicationController (applicationService *service.ApplicationService) *ApplicationController{
	return &ApplicationController{
		CustomerController: newCustomerController(applicationService.CustomerService,applicationService.PlansService),
	}
}