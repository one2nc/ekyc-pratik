package controllers

import (
	"go-ekyc/helper"
	"go-ekyc/model"
	service "go-ekyc/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerControllers struct {
	CustomerService service.CustomerService
	PlansService    service.PlansService
}

func (cc *CustomerControllers) RegisterCustomer(c *gin.Context) {
	type SignupRequest struct {
		Name  string `json:"name" binding:"required"`
		Plan  string `json:"plan" binding:"required,oneof=basic advanced enterprise"`
		Email string `json:"email" binding:"required,email"`
	}
	var signupRequest SignupRequest

	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	plan, err := cc.PlansService.FetchPlansByName(signupRequest.Plan)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	customerData := model.Customer{
		Name:      signupRequest.Name,
		Email:     signupRequest.Email,
		PlanID:    plan.ID,
		AccessKey: helper.GenerateRandomString(10),
		SecretKey: helper.GenerateRandomString(20),
	}

	customer, err := cc.CustomerService.GetCustomerByEmail(customerData.Email)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	if customer != (model.Customer{}) {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": "Email Already taken",
		})
		return
	}

	err = cc.CustomerService.CreateCustomer(&customerData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secretKey": customerData.SecretKey,
		"accessKey": customerData.AccessKey,
	})

}

func newCustomerController(customerService *service.CustomerService, plansService *service.PlansService) *CustomerControllers {
	return &CustomerControllers{
		CustomerService: *customerService,
		PlansService:    *plansService,
	}
}
