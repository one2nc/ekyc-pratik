package image

import (
	"go-ekyc/controllers"
	"go-ekyc/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterImageRoutes(v1Group *gin.RouterGroup, appController *controllers.ApplicationController) {
	imageController := appController.ImageController
	imageGroup := v1Group.Group("/image")
	{
		imageGroup.POST("/upload",
			middlewares.AuthMiddleware(&imageController.CustomerService),
			imageController.UplaodImage)
		imageGroup.POST("/face-match",
			middlewares.AuthMiddleware(&imageController.CustomerService),
			imageController.FaceMatch)
	}
}
