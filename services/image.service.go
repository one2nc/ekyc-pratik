package service

import (
	"go-ekyc/model"
	"go-ekyc/repository"

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

func newImageService(imageRepository *repository.ImageRepository) *ImageService {
	return &ImageService{
		imageRepository: imageRepository,
	}
}
