package auth

import (
	"go-ekyc/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(v1Group *gin.RouterGroup, appController *controllers.ApplicationController) {
	customerController := appController.CustomerController
	authGroup := v1Group.Group("/auth")
	{
		authGroup.POST("/signup",
			customerController.RegisterCustomer)
	}
}
