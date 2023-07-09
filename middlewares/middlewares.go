package middlewares

import (
	"fmt"
	service "go-ekyc/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(customerService *service.CustomerService) func(*gin.Context) {

	return func(c *gin.Context) {
		fmt.Println("middleware applied")
	}
}
