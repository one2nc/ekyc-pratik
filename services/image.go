package service

import (
	"encoding/json"
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
	Data datatypes.JSON
}

type ImageService struct {
	imageRepository          repository.IImageRepository
	plansRepository          repository.IPlansRepository
	faceMatchScoreRepository repository.IFaceMatchScoreRepository
	ocrRepository            repository.IOCRRepository
	ocrService               *OCRService
	minioService             IMinioService
}

func newImageService(imageRepository repository.IImageRepository, plansRepository repository.IPlansRepository, faceMatchScoreRepository repository.IFaceMatchScoreRepository, ocrRepository repository.IOCRRepository, ocrService *OCRService, minioService IMinioService) *ImageService {
	return &ImageService{
		imageRepository:          imageRepository,
		plansRepository:          plansRepository,
		faceMatchScoreRepository: faceMatchScoreRepository,
		ocrRepository:            ocrRepository,
		ocrService:               ocrService,
		minioService:             minioService,
	}
}

func (i *ImageService) UploadImage(input UploadImageInput) (ImageUploadResult, error) {

	plan, err := i.plansRepository.FetchPlanById(input.Customer.PlanID)

	if err != nil {
		return ImageUploadResult{}, ErrPlanNotFound
	}

	file, fileInfo := input.File, input.FileInfo

	defer file.Close()

	filePath := fmt.Sprintf("images/%s/%s_%s", input.Customer.ID, fmt.Sprint(time.Now().Unix()), fileInfo.Filename)
	fileSizeBytes := fileInfo.Size
	fileName := fileInfo.Filename
	fileExtension := filepath.Ext(fileName)
	imageData := model.Image{
		CustomerID:    input.Customer.ID,
		FilePath:      filePath,
		FileExtension: fileExtension,
		FileSizeMB:    float64(fileSizeBytes) / 1000,
		ImageType:     input.ImageType,
	}

	connectionType := "application/" + fileExtension
	err = i.minioService.UploadFileToMinio(filePath, file, fileSizeBytes, connectionType)
	if err != nil {

		return ImageUploadResult{}, ErrUnknown

	}

	err = i.imageRepository.CreateImage(&imageData)

	if err != nil {

		return ImageUploadResult{}, ErrUnknown

	}

	imageStorageCharges := float64(plan.ImageUploadCost) * imageData.FileSizeMB

	imageUploadAPICall := model.ImageUploadAPICall{
		CustomerID:          input.Customer.ID,
		ImageID:             imageData.ID,
		ImageStorageCharges: imageStorageCharges,
	}
	err = i.imageRepository.CreateImageUploadRecord(&imageUploadAPICall)

	if err != nil {
		return ImageUploadResult{}, ErrUnknown
	}

	return ImageUploadResult{ImageId: imageData.ID}, nil
}

func (i *ImageService) FaceMatch(input FaceMatchInput) (FaceMatchResult, error) {
	plan, err := i.plansRepository.FetchPlanById(input.Customer.PlanID)
	if err != nil {
		return FaceMatchResult{}, ErrPlanNotFound
	}

	images, err := i.imageRepository.FindImagesByIdForCustomer([]string{input.ImageId1, input.ImageId2}, input.Customer.ID.String())

	if err != nil {
		return FaceMatchResult{}, err
	}
	if len(images) != 2 {
		return FaceMatchResult{}, ErrImageNotFound
	}

	if !helper.IsImagesComparable(images[0], images[1]) {
		return FaceMatchResult{}, ErrInvalidImageType
	}

	faceMatchScoreRecord, err := i.faceMatchScoreRepository.FetchScoreByImageAndCustomerId(input.ImageId1, input.ImageId2, input.Customer.ID.String())

	if err != nil {
		if gorm.ErrRecordNotFound.Error() != err.Error() {
			return FaceMatchResult{}, ErrUnknown
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
			return FaceMatchResult{}, ErrUnknown
		}
	}
	faceMatchApiData := model.FaceMatchAPICall{
		CustomerID:     input.Customer.ID,
		ScoreID:        faceMatchScoreRecord.ID,
		APICallCharges: plan.FaceMatchCost,
	}
	err = i.faceMatchScoreRepository.CreateFaceMatchScoreAPIRecord(&faceMatchApiData)
	if err != nil {
		return FaceMatchResult{}, ErrUnknown
	}
	return FaceMatchResult{Score: faceMatchScoreRecord.Score}, nil
}

func (i *ImageService) GetOCRData(input OCRInput) (OCRResult, error) {

	plan, err := i.plansRepository.FetchPlanById(input.Customer.PlanID)

	if err != nil {
		return OCRResult{}, ErrPlanNotFound
	}

	images, err := i.imageRepository.FindImagesByIdForCustomer([]string{input.ImageId}, input.Customer.ID.String())

	if err != nil {
		return OCRResult{}, ErrUnknown
	}
	if len(images) != 1 {
		return OCRResult{}, ErrImageNotFound
	}

	if images[0].ImageType != "id_card" {
		return OCRResult{}, ErrInvalidImageType
	}

	ocrData, err := i.ocrRepository.GetOCRDataForCustomerByImageId(input.ImageId, input.Customer.ID.String())

	if err != nil {
		if gorm.ErrRecordNotFound.Error() != err.Error() {
			return OCRResult{}, ErrUnknown
		}

	}
	if ocrData == nil {

		data, err := i.ocrService.OCRExtractData()
		if err != nil {
			return OCRResult{}, ErrUnknown

		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			return OCRResult{}, ErrUnknown
		}
		ocrData = &model.OCRData{
			CustomerID: input.Customer.ID,
			ImageID:    images[0].ID,
			OCRData:    datatypes.JSON(jsonData),
		}

		err = i.ocrRepository.CreateOCRData(ocrData)

		if err != nil {
			return OCRResult{}, ErrUnknown

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
		return OCRResult{}, ErrUnknown

	}
	return OCRResult{
		Data: ocrData.OCRData,
	}, nil
}

func (i *ImageService) GenerateFacteMatchScore() int {

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100) + 1

	return randomNumber

}
