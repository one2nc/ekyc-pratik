package handlers

import (
	"go-ekyc/handlers/requests"
	"go-ekyc/handlers/response"
	"go-ekyc/helper"
	"go-ekyc/model"
	"go-ekyc/repository"
	service "go-ekyc/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerHandlers struct {
	CustomerService service.CustomerService
}

func newCustomerHandler(customerService *service.CustomerService) *CustomerHandlers {
	return &CustomerHandlers{
		CustomerService: *customerService,
	}
}

func (cc *CustomerHandlers) RegisterCustomer(c *gin.Context) {

	var signupRequest requests.SignupRequest

	if err := c.ShouldBindJSON(&signupRequest); err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": helper.ErrorParser(err, &signupRequest)})
		return
	}

	result, err := cc.CustomerService.RegisterCustomer(service.RegisterServiceInput{
		PlanName:      signupRequest.Plan,
		CustomerEmail: signupRequest.Email,
		CustomerName:  signupRequest.Name,
	})
	if err != nil {
		status, errMsg := response.GetHttpStatusAndError(err)
		c.AbortWithStatusJSON(status, gin.H{
			"errorMessages": []string{errMsg.Error()},
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"access_key": result.AccessKey,
		"secret_key": result.SecretKey,
	})

}

func (i *CustomerHandlers) GetAggregatedReport(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := customer.(model.Customer)

	var reportsRequest requests.ReportsRequest

	if err := c.ShouldBindQuery(&reportsRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": helper.ErrorParser(err, &reportsRequest)})
		return
	}
	startDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.StartDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{"Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"}})
		return
	}
	endDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.EndDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{"Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"}})
		return
	}
	if endDate.Before(startDate) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{"end_date should to be greater that start_date"}})
		return
	}
	results, err := i.CustomerService.GetAggregateReportForCustomer(startDate, endDate, []uuid.UUID{customerModel.ID})

	if err != nil {
		status, errMsg := response.GetHttpStatusAndError(err)
		c.AbortWithStatusJSON(status, gin.H{
			"errorMessages": []string{errMsg.Error()},
		})
		return
	}

	if len(results) > 0 {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"report": results[0]})
	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"report": map[string]interface{}{}})
	}

}
func (i *CustomerHandlers) GetAggregatedReportForAllCustomers(c *gin.Context) {
	

	var reportsRequest requests.ReportsRequest

	if err := c.ShouldBindQuery(&reportsRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": helper.ErrorParser(err, &reportsRequest)})
		return
	}
	startDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.StartDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{"Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"}})
		return
	}
	endDate, err := time.Parse("2006-01-02 15:04:05", reportsRequest.EndDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{"Invalid dates. valid format is yyyy-mm-dd hr:mm:ss"}})
		return
	}
	if endDate.Before(startDate) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{"end_date should to be greater that start_date"}})
		return
	}
	results, err := i.CustomerService.GetAggregateReportForCustomer(startDate, endDate, []uuid.UUID{})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{err.Error()}})
		return
	}

	response := struct {
		Reports   []repository.CustomerAggregatedReport `json:"reports"`
		StartDate time.Time                             `json:"start_date_of_report"`
		EndDate   time.Time                             `json:"end_date_of_report"`
	}{
		StartDate: startDate,
		EndDate:   endDate,
		Reports:   results,
	}

	c.AbortWithStatusJSON(http.StatusOK, response)

}
