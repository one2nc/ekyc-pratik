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

func InitiateServer(serverConfig ServerConfig, appController *handlers.ApplicationController) {
	r := gin.Default()
	route.RegisterRoutes(r, appController)
	r.Run(serverConfig.Address + ":" + serverConfig.Port)
}
