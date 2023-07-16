package handlers

import (
	"go-ekyc/handlers/requests"
	service "go-ekyc/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerControllers struct {
	CustomerService service.CustomerService
	PlansService    service.PlansService
}

func (cc *CustomerControllers) RegisterCustomer(c *gin.Context) {

	var signupRequest requests.SignupRequest

	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	result, err := cc.CustomerService.RegisterCustomer(service.RegisterServiceInput{
		PlanName:      signupRequest.Plan,
		CustomerEmail: signupRequest.Email,
		CustomerName:  signupRequest.Name,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"access_key": result.AccessKey,
		"secret_key": result.SecretKey,
	})

}

func newCustomerController(customerService *service.CustomerService, plansService *service.PlansService) *CustomerControllers {
	return &CustomerControllers{
		CustomerService: *customerService,
		PlansService:    *plansService,
	}
}
