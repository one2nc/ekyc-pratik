package mockedrepository

// import (
// 	"go-ekyc/model"
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type FaceMatchScoreMockRepository struct {
// 	FaceMatchScores []model.FaceMatchScore
// }

// func newFaceMatchScoreMockRepository(db *gorm.DB) IFaceMatchScoreRepository {
// 	return &FaceMatchScoreRepository{
// 		dbInstance: db,
// 	}
// }

// func (f *FaceMatchScoreMockRepository) CreateFaceMatchScore(faceScoreData *model.FaceMatchScore) error {
// 	faceScoreData.ID = uuid.New()
// 	return nil
// }
// func (f *FaceMatchScoreMockRepository) CreateFaceMatchScoreAPIRecord(faceMatchApiData *model.FaceMatchAPICall) error {

// 	faceMatchApiData.ID = uuid.New()
// 	return nil
// }

// func (f *FaceMatchScoreMockRepository) FetchScoreByImageAndCustomerId(imageId1 string, imageId2 string, customerId string) (*model.FaceMatchScore, error) {

// 	for _, faceMatchScore := range f.FaceMatchScores {
// 		if faceMatchScore.CustomerID.String() == customerId {
// 			cond1 := faceMatchScore.ImageID1.String() == imageId1 || faceMatchScore.ImageID1.String() == imageId2
// 			cond2 := faceMatchScore.ImageID2.String() == imageId1 || faceMatchScore.ImageID2.String() == imageId2
// 			if cond1 && cond2 {
// 				return &faceMatchScore, nil
// 			}
// 		}
// 	}

// 	return nil, gorm.ErrRecordNotFound
// }

// func (f *FaceMatchScoreMockRepository) GetFaceMatchAPIReport(startDate time.Time, endDate time.Time) (map[uuid.UUID]FaceMatchAPIReport, error) {

// 	rows, err := i.dbInstance.Table("ekyc_schema.face_match_api_calls").
// 		Select("face_match_api_calls.customer_id, SUM(api_call_charges) as total_api_charge, COUNT(*) as total_api_count").
// 		Where("face_match_api_calls.created_at BETWEEN ? AND ?", startDate, endDate).
// 		Group("face_match_api_calls.customer_id").
// 		Rows()

// 	results := map[uuid.UUID]FaceMatchAPIReport{}
// 	if err != nil {
// 		return results, err
// 	}
// 	defer rows.Close()

// 	// Iterate over the rows and retrieve the results
// 	for rows.Next() {
// 		result := FaceMatchAPIReport{}
// 		err := i.dbInstance.ScanRows(rows, &result)
// 		if err != nil {
// 			return results, err
// 		}
// 		results[result.CustomerId] = result
// 	}
// 	return results, nil
// }
