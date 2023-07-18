package handlers

import (
	"encoding/json"
	"fmt"
	"go-ekyc/handlers/requests"
	"go-ekyc/model"
	"go-ekyc/repository"
	service "go-ekyc/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	startDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.StartDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"})
		return
	}
	endDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.EndDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"})
		return
	}
	if endDate.Before(startDate) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "end_date should to be greater that start_date"})
		return
	}
	results, err := i.CustomerService.GetAggregateReportForCustomer(startDate, endDate, []uuid.UUID{customerModel.ID})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	report := results[0]
	report.TotalInvoiceAmount = report.TotalBaseCharge + report.TotalAPICallCharges

	jsonBody, err := json.Marshal(report)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "error while fetching report"})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonBody)
	c.Abort()

}
func (i *CustomerControllers) GetAggregatedReportForAllCustomers(c *gin.Context) {
	// customer, _ := c.Get("customer")
	// customerModel := customer.(model.Customer)

	var reportsRequest requests.ReportsRequest

	if err := c.ShouldBindJSON(&reportsRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	fmt.Println(reportsRequest.StartDate, reportsRequest.EndDate)
	startDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.StartDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"})
		return
	}
	endDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.EndDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"})
		return
	}
	if endDate.Before(startDate) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "end_date should to be greater that start_date"})
		return
	}
	results, err := i.CustomerService.GetAggregateReportForCustomer(startDate, endDate, []uuid.UUID{})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	for i, result := range results {
		results[i].TotalInvoiceAmount = result.TotalBaseCharge + result.TotalAPICallCharges
		results[i].StartDate = startDate
		results[i].EndDate = endDate
	}

	jsonResponse := struct {
		Reports   []repository.CustomerAggregatedReport `json:"reports"`
		StartDate time.Time                             `json:"start_date_of_report"`
		EndDate   time.Time                             `json:"end_date_of_report"`
	}{
		StartDate: startDate,
		EndDate:   endDate,
		Reports:   results,
	}
	jsonBody, err := json.Marshal(jsonResponse)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "error while fetching report"})
		return
	}
	c.Data(http.StatusOK, "application/json", jsonBody)
	c.Abort()

}
