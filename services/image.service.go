package service

import (
	"fmt"
	"go-ekyc/model"
	"go-ekyc/repository"
	"math/rand"
	"time"
)

type ImageService struct {
	imageRepository *repository.ImageRepository
}

func (c *ImageService) CreateImage(image *model.Image) error {
	err := c.imageRepository.CreateImage(image)
	return err
}
func (c *ImageService) CreateImageUploadAPICall(image *model.ImageUploadAPICall) error {
	err := c.imageRepository.CreateImageUploadRecord(image)
	return err
}
func (c *ImageService) FindImageForCustomer(imageIds []string, customerId string) ([]model.Image, error) {
	fmt.Println(imageIds)
	images, err := c.imageRepository.FindImagesByIdForCustomer(imageIds, customerId)

	return images, err

}

func (c *ImageService) GenerateFacteMatchScore() int {

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(101)

	return randomNumber

}

func newImageService(imageRepository *repository.ImageRepository) *ImageService {
	return &ImageService{
		imageRepository: imageRepository,
	}
}
