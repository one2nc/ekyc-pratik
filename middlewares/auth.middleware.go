package middlewares

import (
	service "go-ekyc/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(customerService *service.CustomerService) func(*gin.Context) {

	return func(c *gin.Context) {
		accessKey := c.GetHeader("Access-Key")
		secretKey := c.GetHeader("Secret-Key")

		// Validate access and secret keys
		if accessKey == "" || secretKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errorMessage": "Invalid access or secret key"})
			return
		}

		customer, err := customerService.GetCustomerByCredendials(accessKey, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"errorMessage": err.Error(),
			})

			return
		}
		c.Set("customer", customer)

		c.Next()

	}
}
