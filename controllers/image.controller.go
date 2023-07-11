package controllers

import (
	"fmt"
	"go-ekyc/helper"
	"go-ekyc/model"
	service "go-ekyc/services"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ImageControllers struct {
	CustomerService       service.CustomerService
	PlansService          service.PlansService
	ImageService          service.ImageService
	MinioService          service.MinioService
	FaceMatchScoreService service.FaceMatchScoreService
}

func (i *ImageControllers) UplaodImage(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := model.Customer{}
	customerModel = customer.(model.Customer)

	// fetch plan for calculations
	plan, err := i.PlansService.FetchPlanById(customerModel.PlanID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch plan"})
		return
	}

	file, fileInfo, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve image"})
		return
	}
	defer file.Close()

	imageType := c.Request.PostFormValue("image_type")

	if imageType == "" || !helper.IsImageTypeValid(imageType) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": "Invalid image type. valid type are face or id_card",
		})
		return
	}

	filePath := fmt.Sprintf("images/%s/%s_%s", customerModel.ID, fmt.Sprint(time.Now().Unix()), fileInfo.Filename)
	fileSizeBytes := fileInfo.Size
	fileName := fileInfo.Filename
	fileExtension := filepath.Ext(fileName)
	imageData := model.Image{
		CustomerID:    customerModel.ID,
		FilePath:      filePath,
		FileExtension: fileExtension,
		FileSizeMB:    float64(fileSizeBytes) / 1000,
		ImageType:     imageType,
	}

	minioService, err := service.NewMinioService()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	bucketName := minioService.MinioConfig.ImageBucket
	connectionType := "application/" + fileExtension
	err = minioService.UploadFileToMinio(bucketName, filePath, file, fileSizeBytes, connectionType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	err = i.ImageService.CreateImage(&imageData)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"imageId": imageData.ID,
	})

	imageStorageCharges := float64(plan.ImageUploadCost) * imageData.FileSizeMB

	imageUploadAPICall := model.ImageUploadAPICall{
		CustomerID:          customerModel.ID,
		ImageID:             imageData.ID,
		ImageStorageCharges: imageStorageCharges,
	}
	_ = i.ImageService.CreateImageUploadAPICall(&imageUploadAPICall)

}
func (i *ImageControllers) FaceMatch(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := model.Customer{}
	customerModel = customer.(model.Customer)
	type FaceMatchRequest struct {
		ImageId1 string `json:"image_id_1" binding:"required,uuid"`
		ImageId2 string `json:"image_id_2" binding:"required,uuid"`
	}

	var faceMatchRequest FaceMatchRequest

	if err := c.ShouldBindJSON(&faceMatchRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	if faceMatchRequest.ImageId1 == faceMatchRequest.ImageId2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Cannot use same ids"})
		return
	}
	// fetch plan for calculations
	plan, err := i.PlansService.FetchPlanById(customerModel.PlanID)
	fmt.Println(plan.ID)
	fmt.Println(faceMatchRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Failed to fetch plan"})
		return
	}

	images, err := i.ImageService.FindImageForCustomer([]string{faceMatchRequest.ImageId1, faceMatchRequest.ImageId2}, customerModel.ID.String())

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	if len(images) != 2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid image ids"})
		return
	}

	if !helper.IsImagesComparable(images[0], images[1]) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Invalid image type"})
		return
	}

	faceMatchScoreRecord, err := i.FaceMatchScoreService.FetchScoreByImageAndCustomerId(faceMatchRequest.ImageId1, faceMatchRequest.ImageId2, customerModel.ID.String())

	if err != nil {
		if gorm.ErrRecordNotFound.Error() != err.Error() {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}

	}

	if faceMatchScoreRecord == nil {
		score := i.ImageService.GenerateFacteMatchScore()

		faceMatchScoreRecord = &model.FaceMatchScore{
			CustomerID: customerModel.ID,
			ImageID1:   images[0].ID,
			ImageID2:   images[1].ID,
			Score:      score,
		}

		err = i.FaceMatchScoreService.CreateFaceMatchScore(faceMatchScoreRecord)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
			return
		}
	}
	faceMatchApiData := model.FaceMatchAPICall{
		CustomerID:     customerModel.ID,
		ScoreID:        faceMatchScoreRecord.ID,
		APICallCharges: plan.FaceMatchCost,
	}
	err = i.FaceMatchScoreService.CreateFaceMatchScoreAPIRecord(&faceMatchApiData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{"faceMatchScore": faceMatchScoreRecord.Score})

}

func newImageController(customerService *service.CustomerService, plansService *service.PlansService, imageService *service.ImageService, faceMatchScoreService *service.FaceMatchScoreService) *ImageControllers {
	return &ImageControllers{
		CustomerService:       *customerService,
		PlansService:          *plansService,
		ImageService:          *imageService,
		FaceMatchScoreService: *faceMatchScoreService,
	}
}
