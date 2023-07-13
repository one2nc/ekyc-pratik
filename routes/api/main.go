package api

import (
	"go-ekyc/handlers"
	v1 "go-ekyc/routes/api/v1"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(router *gin.Engine, appController *handlers.ApplicationController) {
	apiGroup := router.Group("/api")

	v1.RegisterV1Routes(apiGroup, appController)

}
