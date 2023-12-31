package repository

import (
	"go-ekyc/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FaceMatchAPIReport struct {
	CustomerId     uuid.UUID
	TotalApiCharge float64
	TotalApiCount  int
}
type IFaceMatchScoreRepository interface{
	CreateFaceMatchScore(faceScoreData *model.FaceMatchScore) error
	CreateFaceMatchScoreAPIRecord(faceMatchApiData *model.FaceMatchAPICall) error
	FetchScoreByImageAndCustomerId(imageId1 string, imageId2 string, customerId string) (*model.FaceMatchScore, error)
	GetFaceMatchAPIReport(startDate time.Time, endDate time.Time, customerIds []uuid.UUID) (map[uuid.UUID]FaceMatchAPIReport, error)
}
type FaceMatchScoreRepository struct {
	dbInstance *gorm.DB
}

func newFaceMatchScoreRepository(db *gorm.DB) IFaceMatchScoreRepository {
	return &FaceMatchScoreRepository{
		dbInstance: db,
	}
}

func (f *FaceMatchScoreRepository) CreateFaceMatchScore(faceScoreData *model.FaceMatchScore) error {
	result := f.dbInstance.Create(&faceScoreData)

	return result.Error
}
func (f *FaceMatchScoreRepository) CreateFaceMatchScoreAPIRecord(faceMatchApiData *model.FaceMatchAPICall) error {
	result := f.dbInstance.Create(&faceMatchApiData)

	return result.Error
}
func (f *FaceMatchScoreRepository) FetchScoreByImageAndCustomerId(imageId1 string, imageId2 string, customerId string) (*model.FaceMatchScore, error) {

	faceImageScore := model.FaceMatchScore{}
	result := f.dbInstance.Where("(image_id_1 = ? OR image_id_1 = ?) AND (image_id_2 = ? OR image_id_2 = ?) AND customer_id = ? ", imageId1, imageId2, imageId1, imageId2, customerId).First(&faceImageScore)

	if result.Error != nil {
		return nil, result.Error
	}
	return &faceImageScore, nil
}

func (i *FaceMatchScoreRepository) GetFaceMatchAPIReport(startDate time.Time, endDate time.Time, customerIds []uuid.UUID)  (map[uuid.UUID]FaceMatchAPIReport, error) {

	rows, err := i.dbInstance.Table("ekyc_schema.face_match_api_calls").
		Select("face_match_api_calls.customer_id, SUM(api_call_charges) as total_api_charge, COUNT(*) as total_api_count").
		Where("face_match_api_calls.created_at BETWEEN ? AND ? AND customer_id IN (?)", startDate, endDate,customerIds).
		Group("face_match_api_calls.customer_id").
		Rows()

		results := map[uuid.UUID]FaceMatchAPIReport{}
	if err != nil {
		return results,err
	}
	defer rows.Close()

	for rows.Next() {
		result := FaceMatchAPIReport{}
		err := i.dbInstance.ScanRows(rows, &result)
		if err != nil {
			return results, err
		}
		results[result.CustomerId] = result
	}
	return results, nil
}
