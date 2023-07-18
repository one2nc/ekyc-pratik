package v1

import (
	"go-ekyc/handlers"
	"go-ekyc/routes/api/v1/auth"
	"go-ekyc/routes/api/v1/image"
	"go-ekyc/routes/api/v1/report"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(apiGroup *gin.RouterGroup, appHandler *handlers.ApplicationHandler) {
	v1Group := apiGroup.Group("/v1")
	auth.RegisterAuthRoutes(v1Group, appHandler)
	image.RegisterImageRoutes(v1Group, appHandler)
	report.RegisterReportRoutes(v1Group, appHandler)
}
