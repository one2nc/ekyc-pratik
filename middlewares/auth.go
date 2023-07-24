package middlewares

import (
	"go-ekyc/handlers/response"
	"go-ekyc/helper"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errorMessages": []string{"Invalid access or secret key"}})
			return
		}

		customer, err := customerService.GetCustomerByCredendials(helper.GetMD5Hash(accessKey), helper.GetMD5Hash(secretKey))
		if err != nil {
			status, errMsg := response.GetHttpStatusAndError(err)
			c.AbortWithStatusJSON(status, gin.H{
				"errorMessages": []string{errMsg.Error()},
			})
			return

			
		}
		c.Set("customer", customer)

		c.Next()

	}
}
