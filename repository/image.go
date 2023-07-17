package repository

import (
	"go-ekyc/model"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IImageRepository interface {
	CreateImage(image *model.Image) error
	CreateImageUploadRecord(imageUploadData *model.ImageUploadAPICall) error
	FindImagesByIdForCustomer(imageIds []string, customerId string) ([]model.Image, error)
	GetImageUploadAPIReport(startDate time.Time, endDate time.Time) (map[uuid.UUID]ImageUploadAPIReport, error)
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

type ImageUploadAPIReport struct {
	CustomerId         uuid.UUID
	TotalUploadCharges float64
	TotalImageSize     float64
	TotalApiCount      int32
}

func (i *ImageRepository) GetImageUploadAPIReport(startDate time.Time, endDate time.Time) (map[uuid.UUID]ImageUploadAPIReport, error) {

	rows, err := i.dbInstance.Table("ekyc_schema.image_upload_api_calls").
		Select("image_upload_api_calls.customer_id, SUM(images.file_size_mb) AS total_image_size, SUM(image_storage_charges) as total_upload_charges, COUNT(*) as total_api_count").
		Joins("JOIN ekyc_schema.images ON image_upload_api_calls.image_id=images.id").
		Where("image_upload_api_calls.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("image_upload_api_calls.customer_id").
		Rows()

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and retrieve the results
	results := map[uuid.UUID]ImageUploadAPIReport{}
	for rows.Next() {
		result := ImageUploadAPIReport{}
		err := i.dbInstance.ScanRows(rows, &result)
		if err != nil {
			return results, err
		}
		results[result.CustomerId] = result
	}
	return results, nil
}
func (i *ImageRepository) GetOcrUploadAPIReport(startDate time.Time, endDate time.Time) (map[uuid.UUID]ImageUploadAPIReport, error) {

	rows, err := i.dbInstance.Table("ekyc_schema.image_upload_api_calls").
		Select("image_upload_api_calls.customer_id, SUM(images.file_size_mb) AS total_image_size, SUM(image_storage_charges) as total_upload_charges, COUNT(*) as total_api_count").
		Joins("JOIN ekyc_schema.images ON image_upload_api_calls.image_id=images.id").
		Where("image_upload_api_calls.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("image_upload_api_calls.customer_id").
		Rows()

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and retrieve the results
	results := map[uuid.UUID]ImageUploadAPIReport{}
	for rows.Next() {
		result := ImageUploadAPIReport{}
		err := i.dbInstance.ScanRows(rows, &result)
		if err != nil {
			return results, err
		}
		results[result.CustomerId] = result
	}
	return results, nil
}

func newImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{
		dbInstance: db,
	}
}
