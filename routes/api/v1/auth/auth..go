package auth

import (
	"go-ekyc/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(v1Group *gin.RouterGroup, appHandler *handlers.ApplicationHandler) {
	customerHandler := appHandler.CustomerHandler
	authGroup := v1Group.Group("/auth")
	{
		authGroup.POST("/signup",
			customerHandler.RegisterCustomer)
	}
}
