package repository

import (
	"go-ekyc/model"

	"gorm.io/gorm"
)

type OCRRepository struct {
	dbInstance *gorm.DB
}

func (o *OCRRepository) CreateOCRData(ocrData *model.OCRData) error {

	result := o.dbInstance.Create(&ocrData)
	return result.Error
}
func (o *OCRRepository) GetOCRDataForCustomerByImageId(imageId string, customerId string) (*model.OCRData, error) {

	ocrRecord := &model.OCRData{}
	result := o.dbInstance.Where("image_id = ?  and customer_id = ?",imageId,customerId).First(ocrRecord)
	if result.Error != nil {
		return nil, result.Error
	}

	return ocrRecord, result.Error
}

func (o *OCRRepository) CreateOcrAPICall(ocrDataModel *model.OCRAPICalls) error {

	result := o.dbInstance.Create(&ocrDataModel)

	return result.Error
}

func newOCRRepositoty(db *gorm.DB) *OCRRepository {

	return &OCRRepository{
		dbInstance: db,
	}
}
