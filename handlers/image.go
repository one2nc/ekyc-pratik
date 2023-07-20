package handlers

import (
	"go-ekyc/handlers/requests"
	"go-ekyc/helper"
	"go-ekyc/model"
	service "go-ekyc/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageHandlers struct {
	CustomerService service.CustomerService
	ImageService    service.ImageService
}

func (i *ImageHandlers) UplaodImage(c *gin.Context) {
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
			"errorMessages": []string{"Invalid image type. valid type are face or id_card"},
		})
		return
	}

	result, err := i.ImageService.UploadImage(service.UploadImageInput{
		Customer:  customerModel,
		File:      file,
		FileInfo:  fileInfo,
		ImageType: imageType,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessages": []string{err.Error()},
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"image_id": result.ImageId.String(),
	})
}
func (i *ImageHandlers) FaceMatch(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := customer.(model.Customer)

	var faceMatchRequest requests.FaceMatchRequest

	if err := c.ShouldBindJSON(&faceMatchRequest); err != nil {

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": helper.ErrorParser(err)})
		return
	}

	if faceMatchRequest.ImageId1 == faceMatchRequest.ImageId2 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{"Cannot use same ids"}})
		return
	}
	results, err := i.ImageService.FaceMatch(service.FaceMatchInput{
		Customer: customerModel,
		ImageId1: faceMatchRequest.ImageId1,
		ImageId2: faceMatchRequest.ImageId2,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": []string{err.Error()}})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{"faceMatchScore": results.Score})

}
func (i *ImageHandlers) GetOcrData(c *gin.Context) {
	customer, _ := c.Get("customer")
	customerModel := customer.(model.Customer)

	var ocrRequest requests.OCRRequest

	if err := c.ShouldBindJSON(&ocrRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errorMessages": helper.ErrorParser(err)})
		return
	}

	// fetch plan for calculations
	result, err := i.ImageService.GetOCRData(service.OCRInput{Customer: customerModel, ImageId: ocrRequest.ImageId1})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"errorMessages": []string{err.Error()},
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusCreated, gin.H{
		"data": result.Data,
	})

}

func newImageHandler(customerService *service.CustomerService, imageService *service.ImageService) *ImageHandlers {
	return &ImageHandlers{
		CustomerService: *customerService,

		ImageService: *imageService,
	}
}
