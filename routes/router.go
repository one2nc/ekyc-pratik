package routes

import (
	"go-ekyc/handlers"
	"go-ekyc/routes/api"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, appHandler *handlers.ApplicationHandler) {
	api.RegisterAPIRoutes(router, appHandler)
}
