package main

import (
	"go-ekyc/controllers"
	"go-ekyc/repository"
	"go-ekyc/server"
	service "go-ekyc/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	applicationRepository, err := repository.NewApplicationRepository()
	if err != nil {
		log.Fatal(err.Error())
	}

	appService := service.NewApplicationService(applicationRepository)

	appController := controllers.NewApplicationController(appService)

	serverConfig := server.ServerConfig{
		Port:    os.Getenv("SERVER_PORT"),
		Address: os.Getenv("SERVER_HOST"),
	}
	server.InitiateServer(serverConfig, appController)
}
