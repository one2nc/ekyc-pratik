package v1

import (
	"go-ekyc/handlers"
	"go-ekyc/routes/api/v1/auth"
	"go-ekyc/routes/api/v1/image"
	"go-ekyc/routes/api/v1/report"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(apiGroup *gin.RouterGroup, appController *handlers.ApplicationController) {
	v1Group := apiGroup.Group("/v1")
	auth.RegisterAuthRoutes(v1Group, appController)
	image.RegisterImageRoutes(v1Group, appController)
	report.RegisterReportRoutes(v1Group, appController)
}
