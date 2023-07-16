package repository

import (
	"go-ekyc/model"

	"gorm.io/gorm"
)

type IImageRepository interface{
	CreateImage(image *model.Image) error
	CreateImageUploadRecord(imageUploadData *model.ImageUploadAPICall) error
	FindImagesByIdForCustomer(imageIds []string, customerId string) ([]model.Image, error)
}

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
func (i *ImageRepository) FindImagesByIdForCustomer(imageIds []string, customerId string) ([]model.Image, error) {
	
	images := []model.Image{}

	result := i.dbInstance.Where("id IN (?) AND customer_id = ? ", imageIds, customerId).Find(&images)
	return images, result.Error
}


func newImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		dbInstance: db,
	}
}
