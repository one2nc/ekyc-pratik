package report

import (
	"go-ekyc/handlers"
	"go-ekyc/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(v1Group *gin.RouterGroup, appController *handlers.ApplicationController) {
	customerController := appController.CustomerController
	authGroup := v1Group.Group("/report")
	{
		authGroup.POST("/", middlewares.AuthMiddleware(&customerController.CustomerService),
			customerController.GetAggregatedReport)
	}
}
