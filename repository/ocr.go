package repository

import (
	"go-ekyc/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OCRRepository struct {
	dbInstance *gorm.DB
}

type OCRAPIReport struct {
	CustomerId     uuid.UUID
	TotalApiCharge float64
	TotalApiCount  int
}

func newOCRRepositoty(db *gorm.DB) *OCRRepository {

	return &OCRRepository{
		dbInstance: db,
	}
}

func (o *OCRRepository) CreateOCRData(ocrData *model.OCRData) error {

	result := o.dbInstance.Create(&ocrData)
	return result.Error
}
func (o *OCRRepository) GetOCRDataForCustomerByImageId(imageId string, customerId string) (*model.OCRData, error) {

	ocrRecord := &model.OCRData{}
	result := o.dbInstance.Where("image_id = ?  and customer_id = ?", imageId, customerId).First(ocrRecord)
	if result.Error != nil {
		return nil, result.Error
	}

	return ocrRecord, result.Error
}

func (o *OCRRepository) CreateOcrAPICall(ocrDataModel *model.OCRAPICalls) error {

	result := o.dbInstance.Create(&ocrDataModel)

	return result.Error
}

func (i *OCRRepository) GetOCRAPIReport(startDate time.Time, endDate time.Time) (map[uuid.UUID]OCRAPIReport, error) {

	rows, err := i.dbInstance.Table("ekyc_schema.ocr_api_calls").
		Select("ocr_api_calls.customer_id, SUM(api_call_charges) as total_api_charge, COUNT(*) as total_api_count").
		Where("ocr_api_calls.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("ocr_api_calls.customer_id").
		Rows()

	results := map[uuid.UUID]OCRAPIReport{}
	if err != nil {
		return results, err
	}
	defer rows.Close()

	// Iterate over the rows and retrieve the results
	for rows.Next() {
		result := OCRAPIReport{}
		err := i.dbInstance.ScanRows(rows, &result)
		if err != nil {
			return results, err
		}
		results[result.CustomerId] = result
	}
	return results, nil
}
