package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-ekyc/helper"
	"go-ekyc/model"
	"go-ekyc/repository"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UploadImageInput struct {
	Customer  model.Customer
	File      multipart.File
	FileInfo  *multipart.FileHeader
	ImageType string
}
type FaceMatchInput struct {
	Customer model.Customer
	ImageId1 string
	ImageId2 string
}
type OCRInput struct {
	Customer model.Customer
	ImageId  string
}
type ImageUploadResult struct {
	ImageId uuid.UUID
}
type FaceMatchResult struct {
	Score int
}
type OCRResult struct {
	Data interface{}
}

type ImageService struct {
	imageRepository          repository.IImageRepository
	plansRepository          repository.IPlansRepository
	faceMatchScoreRepository repository.IFaceMatchScoreRepository
	ocrRepository            repository.IOCRRepository
	ocrService               *OCRService
}

func newImageService(imageRepository repository.IImageRepository, plansRepository repository.IPlansRepository, faceMatchScoreRepository repository.IFaceMatchScoreRepository, ocrRepository repository.IOCRRepository, ocrService *OCRService) *ImageService {
	return &ImageService{
		imageRepository:          imageRepository,
		plansRepository:          plansRepository,
		faceMatchScoreRepository: faceMatchScoreRepository,
		ocrRepository:            ocrRepository,
		ocrService:               ocrService,
	}
}


func (i *ImageService) CreateImage(image *model.Image) error {
	err := i.imageRepository.CreateImage(image)
	return err
}
func (i *ImageService) UploadImage(input UploadImageInput) (ImageUploadResult, error) {

	plan, err := i.plansRepository.FetchPlanById(input.Customer.PlanID)

	if err != nil {
		return ImageUploadResult{}, err
	}

	file, fileInfo := input.File, input.FileInfo

	defer file.Close()

	imageType := input.ImageType

	if imageType == "" || !helper.IsImageTypeValid(imageType) {

		return ImageUploadResult{}, errors.New("Invalid image type. valid type are face or id_card")

	}

	filePath := fmt.Sprintf("images/%s/%s_%s", input.Customer.ID, fmt.Sprint(time.Now().Unix()), fileInfo.Filename)
	fileSizeBytes := fileInfo.Size
	fileName := fileInfo.Filename
	fileExtension := filepath.Ext(fileName)
	imageData := model.Image{
		CustomerID:    input.Customer.ID,
		FilePath:      filePath,
		FileExtension: fileExtension,
		FileSizeMB:    float64(fileSizeBytes) / 1000,
		ImageType:     imageType,
	}

	minioService, err := NewMinioService()

	if err != nil {

		return ImageUploadResult{}, err

	}
	bucketName := minioService.MinioConfig.ImageBucket
	connectionType := "application/" + fileExtension
	err = minioService.UploadFileToMinio(bucketName, filePath, file, fileSizeBytes, connectionType)
	if err != nil {

		return ImageUploadResult{}, err

	}

	err = i.imageRepository.CreateImage(&imageData)

	if err != nil {

		return ImageUploadResult{}, err

	}

	imageStorageCharges := float64(plan.ImageUploadCost) * imageData.FileSizeMB

	imageUploadAPICall := model.ImageUploadAPICall{
		CustomerID:          input.Customer.ID,
		ImageID:             imageData.ID,
		ImageStorageCharges: imageStorageCharges,
	}
	err = i.imageRepository.CreateImageUploadRecord(&imageUploadAPICall)

	if err != nil {
		return ImageUploadResult{}, err
	}

	return ImageUploadResult{ImageId: imageData.ID}, nil
}

func (i *ImageService) FaceMatch(input FaceMatchInput) (FaceMatchResult, error) {
	plan, err := i.plansRepository.FetchPlanById(input.Customer.PlanID)
	if err != nil {
		return FaceMatchResult{}, errors.New("Failed to fetch plan")
	}

	images, err := i.imageRepository.FindImagesByIdForCustomer([]string{input.ImageId1, input.ImageId2}, input.Customer.ID.String())

	if err != nil {
		return FaceMatchResult{}, err
	}
	if len(images) != 2 {
		return FaceMatchResult{}, errors.New("invalid image ids")
	}

	if !helper.IsImagesComparable(images[0], images[1]) {
		return FaceMatchResult{}, errors.New("Invalid image type")
	}

	faceMatchScoreRecord, err := i.faceMatchScoreRepository.FetchScoreByImageAndCustomerId(input.ImageId1, input.ImageId2, input.Customer.ID.String())

	if err != nil {
		if gorm.ErrRecordNotFound.Error() != err.Error() {
			return FaceMatchResult{}, err
		}

	}

	if faceMatchScoreRecord == nil {
		score := i.GenerateFacteMatchScore()

		faceMatchScoreRecord = &model.FaceMatchScore{
			CustomerID: input.Customer.ID,
			ImageID1:   images[0].ID,
			ImageID2:   images[1].ID,
			Score:      score,
		}

		err = i.faceMatchScoreRepository.CreateFaceMatchScore(faceMatchScoreRecord)
		if err != nil {
			return FaceMatchResult{}, err
		}
	}
	faceMatchApiData := model.FaceMatchAPICall{
		CustomerID:     input.Customer.ID,
		ScoreID:        faceMatchScoreRecord.ID,
		APICallCharges: plan.FaceMatchCost,
	}
	err = i.faceMatchScoreRepository.CreateFaceMatchScoreAPIRecord(&faceMatchApiData)
	if err != nil {
		return FaceMatchResult{}, err
	}
	return FaceMatchResult{Score: faceMatchScoreRecord.Score}, err
}

func (i *ImageService) GetOCRData(input OCRInput) (OCRResult, error) {

	plan, err := i.plansRepository.FetchPlanById(input.Customer.PlanID)

	if err != nil {
		return OCRResult{}, errors.New("Failed to fetch plan")
	}

	images, err := i.imageRepository.FindImagesByIdForCustomer([]string{input.ImageId}, input.Customer.ID.String())

	if err != nil {
		return OCRResult{}, err
	}
	if len(images) != 1 {
		return OCRResult{}, errors.New("Invalid image id")
	}

	if images[0].ImageType != "id_card" {
		return OCRResult{}, errors.New("OCR can be done only on image of type ID")
	}

	ocrData, err := i.ocrRepository.GetOCRDataForCustomerByImageId(input.ImageId, input.Customer.ID.String())

	if err != nil {
		if gorm.ErrRecordNotFound.Error() != err.Error() {
			return OCRResult{}, err
		}

	}
	if ocrData == nil {

		data, err := i.ocrService.OCRExtractData()
		if err != nil {
			return OCRResult{}, errors.New("Error while extracting data")

		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return OCRResult{}, errors.New("Error while extracting data")
		}
		ocrData = &model.OCRData{
			CustomerID: input.Customer.ID,
			ImageID1:   images[0].ID,
			OCRData:    datatypes.JSON(jsonData),
		}

		err = i.ocrRepository.CreateOCRData(ocrData)

		if err != nil {
			return OCRResult{}, err

		}
	}

	ocrAPICallData := &model.OCRAPICalls{
		CustomerID:     input.Customer.ID,
		ImageID:        images[0].ID,
		OCRID:          ocrData.ID,
		APICallCharges: plan.OCRCost,
	}

	err = i.ocrRepository.CreateOcrAPICall(ocrAPICallData)
	if err != nil {
		return OCRResult{}, err

	}
	return OCRResult{
		Data: ocrData.OCRData,
	}, nil
}

func (i *ImageService) CreateImageUploadAPICall(image *model.ImageUploadAPICall) error {
	err := i.imageRepository.CreateImageUploadRecord(image)
	return err
}
func (i *ImageService) FindImageForCustomer(imageIds []string, customerId string) ([]model.Image, error) {
	images, err := i.imageRepository.FindImagesByIdForCustomer(imageIds, customerId)

	return images, err

}

func (i *ImageService) GenerateFacteMatchScore() int {

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(101)

	return randomNumber

}

