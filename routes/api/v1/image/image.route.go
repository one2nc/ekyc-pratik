package image

import (
	"go-ekyc/controllers"
	"go-ekyc/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterImageRoutes(v1Group *gin.RouterGroup, appController *controllers.ApplicationController) {
	imageController := appController.ImageController
	authGroup := v1Group.Group("/image")
	{
		authGroup.POST("/upload",
			middlewares.AuthMiddleware(&imageController.CustomerService),
			imageController.UplaodImage)
	}
}
