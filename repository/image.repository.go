package repository

import (
	"go-ekyc/model"

	"gorm.io/gorm"
)

type ImageRepository struct {
	dbInstance *gorm.DB
}

func (i *ImageRepository) CreateImage(image *model.Image) error {

	result := i.dbInstance.Create(&image)

	return result.Error

}

func (i *ImageRepository) CreateImageUploadRecord(imageUploadData *model.ImageUploadAPICall) error {

	result := i.dbInstance.Create(&imageUploadData)

	return result.Error
}

func newImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		dbInstance: db,
	}
}
