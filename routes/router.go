package routes

import (
	"go-ekyc/controllers"
	"go-ekyc/routes/api"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, appController *controllers.ApplicationController) {
	api.RegisterAPIRoutes(router, appController)
}
