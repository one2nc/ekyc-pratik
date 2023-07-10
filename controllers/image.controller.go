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
)

type ImageControllers struct {
	CustomerService service.CustomerService
	PlansService    service.PlansService
	ImageService    service.ImageService
	MinioService    service.MinioService
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

func newImageController(customerService *service.CustomerService, plansService *service.PlansService, imageService *service.ImageService) *ImageControllers {
	return &ImageControllers{
		CustomerService: *customerService,
		PlansService:    *plansService,
		ImageService:    *imageService,
	}
}
