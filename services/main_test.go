package service

import (
	"go-ekyc/db"
	"go-ekyc/model"
	"go-ekyc/repository"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var appRepository *repository.ApplicationRepository
var appService *ApplicationService
var dbInstance *gorm.DB

func TestMain(m *testing.M) {
	os.Chdir("./..")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initiate and Migrate DB
	db, err := db.InitiateDB()
	dbInstance = db
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

func TearDownTables(t *testing.T) {
	dbInstance.Where("1=1").Delete(&model.DailyReport{})
	dbInstance.Where("1=1").Delete(&model.FaceMatchAPICall{})
	dbInstance.Where("1=1").Delete(&model.FaceMatchScore{})
	dbInstance.Where("1=1").Delete(&model.OCRAPICalls{})
	dbInstance.Where("1=1").Delete(&model.OCRData{})
	dbInstance.Where("1=1").Delete(&model.ImageUploadAPICall{})
	dbInstance.Where("1=1").Delete(&model.Image{})
	dbInstance.Where("1=1").Delete(&model.Customer{})
}
