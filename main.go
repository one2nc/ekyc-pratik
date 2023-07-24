package main

import (
	"go-ekyc/config"
	"go-ekyc/crons"
	"go-ekyc/db"
	"go-ekyc/handlers"
	"go-ekyc/repository"
	"go-ekyc/server"
	service "go-ekyc/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load enviroment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initiate and Migrate DB
	db, err := db.InitiateDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize application repository
	applicationRepository, err := repository.NewApplicationRepository(db)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize application service
	minioConfig := config.NewMinioConfig()
	minioService, err := service.NewMinioService(minioConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	appService := service.NewApplicationService(applicationRepository, minioService)

	// Register cron jobs
	cronConfig := config.NewCronConfig()
	crons.RegisterCron(appService, cronConfig)

	// Register handlers
	appHandler := handlers.NewApplicationHandler(appService)

	// Initiate server
	serverConfig := server.ServerConfig{
		Port:    os.Getenv("SERVER_PORT"),
		Address: os.Getenv("SERVER_HOST"),
	}
	server.InitiateServer(serverConfig, appHandler)
}
