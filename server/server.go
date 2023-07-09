package server

import (
	"go-ekyc/controllers"
	route "go-ekyc/routes"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Port    string
	Address string
}

func InitiateServer(serverConfig ServerConfig, appController *controllers.ApplicationController) {
	r := gin.Default()
	route.RegisterRoutes(r, appController)
	r.Run(serverConfig.Address + ":" + serverConfig.Port)
}
