package v1

import (
	"go-ekyc/controllers"
	"go-ekyc/routes/api/v1/auth"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(apiGroup *gin.RouterGroup, appController *controllers.ApplicationController) {
	v1Group := apiGroup.Group("/v1")
	auth.RegisterAuthRoutes(v1Group, appController)
}
