package handlers

import (
	"go-ekyc/handlers/requests"
	"go-ekyc/helper"
	"go-ekyc/model"
	service "go-ekyc/services"
	service_inputs "go-ekyc/services/service-input"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageControllers struct {
	CustomerService service.CustomerService
	ImageService    service.ImageService
}

func (i *ImageControllers) UplaodImage(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := customer.(model.Customer)

	// fetch plan for calculations
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

	result, err := i.ImageService.UploadImage(service_inputs.UploadImageInput{
		Customer:  customerModel,
		File:      file,
		FileInfo:  fileInfo,
		ImageType: imageType,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"image_id": result.ImageId.String(),
	})
}
func (i *ImageControllers) FaceMatch(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := customer.(model.Customer)

	var faceMatchRequest requests.FaceMatchRequest

	if err := c.ShouldBindJSON(&faceMatchRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	if faceMatchRequest.ImageId1 == faceMatchRequest.ImageId2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": "Cannot use same ids"})
		return
	}
	results, err := i.ImageService.FaceMatch(service_inputs.FaceMatchInput{
		Customer: customerModel,
		ImageId1: faceMatchRequest.ImageId1,
		ImageId2: faceMatchRequest.ImageId2,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{"faceMatchScore": results.Score})

}
func (i *ImageControllers) GetOcrData(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := customer.(model.Customer)

	var ocrRequest requests.OCRRequest

	if err := c.ShouldBindJSON(&ocrRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessage": err.Error()})
		return
	}

	// fetch plan for calculations
	result, err := i.ImageService.GetOCRData(service_inputs.OCRInput{Customer: customerModel, ImageId: ocrRequest.ImageId1})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessage": err,
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"data": result.Data,
	})

}

func newImageController(customerService *service.CustomerService, imageService *service.ImageService) *ImageControllers {
	return &ImageControllers{
		CustomerService: *customerService,

		ImageService: *imageService,
	}
}
