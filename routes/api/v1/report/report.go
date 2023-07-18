package report

import (
	"go-ekyc/handlers"
	"go-ekyc/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterReportRoutes(v1Group *gin.RouterGroup, appHandler *handlers.ApplicationHandler) {
	customerHandler := appHandler.CustomerHandler
	authGroup := v1Group.Group("/report")
	{
		authGroup.POST("/", middlewares.AuthMiddleware(&customerHandler.CustomerService),
			customerHandler.GetAggregatedReport)
		authGroup.POST("/get-all-reports", middlewares.AuthMiddleware(&customerHandler.CustomerService),
			customerHandler.GetAggregatedReportForAllCustomers)
	}
}
