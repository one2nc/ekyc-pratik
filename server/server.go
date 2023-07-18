package server

import (
	"go-ekyc/handlers"
	route "go-ekyc/routes"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Port    string
	Address string
}

func InitiateServer(serverConfig ServerConfig, appHandler *handlers.ApplicationHandler) {
	r := gin.Default()
	route.RegisterRoutes(r, appHandler)
	r.Run(serverConfig.Address + ":" + serverConfig.Port)
}
