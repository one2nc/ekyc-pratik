package handlers

import (
	"encoding/json"
	"fmt"
	"go-ekyc/handlers/requests"
	"go-ekyc/model"
	service "go-ekyc/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CustomerControllers struct {
	CustomerService service.CustomerService
	PlansService    service.PlansService
}

func newCustomerController(customerService *service.CustomerService, plansService *service.PlansService) *CustomerControllers {
	return &CustomerControllers{
		CustomerService: *customerService,
		PlansService:    *plansService,
	}
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

func (i *CustomerControllers) GetAggregatedReport(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := customer.(model.Customer)

	var reportsRequest requests.ReportsRequest

	if err := c.ShouldBindJSON(&reportsRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	fmt.Println(reportsRequest.StartDate, reportsRequest.EndDate)
	startDate, err := time.Parse("2006-01-02", reportsRequest.StartDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid dates. valid format is yyyy-mm-dd"})
		return
	}
	endDate, err := time.Parse("2006-01-02", reportsRequest.EndDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid dates. valid format is yyyy-mm-dd"})
		return
	}
	results, err := i.CustomerService.GetAggregateReportForCustomer(startDate, endDate, customerModel.ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	jsonBody, err := json.Marshal(results)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "error while fetching report"})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonBody)
	c.Abort()

}
