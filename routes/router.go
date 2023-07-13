package routes

import (
	"go-ekyc/handlers"
	"go-ekyc/routes/api"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, appController *handlers.ApplicationController) {
	api.RegisterAPIRoutes(router, appController)
}
