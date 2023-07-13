package repository

import (
	"go-ekyc/model"

	"gorm.io/gorm"
)

type FaceMatchScoreRepository struct {
	dbInstance *gorm.DB
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
	result := f.dbInstance.Where("(image_id_1 = ? OR image_id_1 = ?) AND (image_id_2 = ? OR image_id_2 = ?) AND customer_id = ? ", imageId1,imageId2,imageId1,imageId2,customerId).First(&faceImageScore)

	if result.Error != nil {
		return nil, result.Error
	}
	return &faceImageScore, nil
}

func newFaceMatchScoreRepository(db *gorm.DB) *FaceMatchScoreRepository {
	return &FaceMatchScoreRepository{
		dbInstance: db,
	}
}
