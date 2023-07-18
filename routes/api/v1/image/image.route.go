package image

import (
	"go-ekyc/handlers"
	"go-ekyc/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterImageRoutes(v1Group *gin.RouterGroup, appHandler *handlers.ApplicationHandler) {
	imageHandler := appHandler.ImageHandler
	imageGroup := v1Group.Group("/image")
	{
		imageGroup.POST("/upload",
			middlewares.AuthMiddleware(&imageHandler.CustomerService),
			imageHandler.UplaodImage)
		imageGroup.POST("/face-match",
			middlewares.AuthMiddleware(&imageHandler.CustomerService),
			imageHandler.FaceMatch)
		imageGroup.POST("/ocr",
			middlewares.AuthMiddleware(&imageHandler.CustomerService),
			imageHandler.GetOcrData)
	}
}
