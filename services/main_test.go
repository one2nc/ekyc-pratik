package service

import (
	"go-ekyc/db"
	"go-ekyc/repository"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var appRepository *repository.ApplicationRepository
var appService *ApplicationService

func TestMain(m *testing.M) {
	os.Chdir("./..")
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
	repo, err := repository.NewApplicationRepository(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	appRepository = repo
	minioService, err := NewMinioMockService()

	if err != nil {
		log.Fatal(err.Error())
	}
	// Initialize application service
	service := NewApplicationService(appRepository, minioService)

	appService = service
	m.Run()
}
