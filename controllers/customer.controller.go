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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	plan, err := cc.PlansService.FetchPlansByName(signupRequest.Plan)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	accessKey := helper.GenerateRandomString(10)
	secretKey := helper.GenerateRandomString(20)
	customerData := model.Customer{
		Name:      signupRequest.Name,
		Email:     signupRequest.Email,
		PlanID:    plan.ID,
		AccessKey: helper.GetMD5Hash(accessKey),
		SecretKey: helper.GetMD5Hash(secretKey),
	}

	customer, err := cc.CustomerService.GetCustomerByEmail(customerData.Email)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	if customer != (model.Customer{}) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "Email Already taken",
		})
		return
	}

	err = cc.CustomerService.CreateCustomer(&customerData)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"access_key": accessKey,
		"secret_key": secretKey,
	})

}

func newCustomerController(customerService *service.CustomerService, plansService *service.PlansService) *CustomerControllers {
	return &CustomerControllers{
		CustomerService: *customerService,
		PlansService:    *plansService,
	}
}
