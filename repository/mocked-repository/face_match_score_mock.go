package mockedrepository

import (
	"go-ekyc/model"
	"go-ekyc/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FaceMatchScoreMockRepository struct {
	FaceMatchScores []model.FaceMatchScore
}

func newFaceMatchScoreMockRepository(faceMatchScores []model.FaceMatchScore) repository.IFaceMatchScoreRepository {
	return &FaceMatchScoreMockRepository{
		FaceMatchScores: faceMatchScores,
	}
}

func (f *FaceMatchScoreMockRepository) CreateFaceMatchScore(faceScoreData *model.FaceMatchScore) error {
	faceScoreData.ID = uuid.New()
	return nil
}
func (f *FaceMatchScoreMockRepository) CreateFaceMatchScoreAPIRecord(faceMatchApiData *model.FaceMatchAPICall) error {

	faceMatchApiData.ID = uuid.New()
	return nil
}

func (f *FaceMatchScoreMockRepository) FetchScoreByImageAndCustomerId(imageId1 string, imageId2 string, customerId string) (*model.FaceMatchScore, error) {

	for _, faceMatchScore := range f.FaceMatchScores {
		if faceMatchScore.CustomerID.String() == customerId {
			cond1 := faceMatchScore.ImageID1.String() == imageId1 || faceMatchScore.ImageID1.String() == imageId2
			cond2 := faceMatchScore.ImageID2.String() == imageId1 || faceMatchScore.ImageID2.String() == imageId2
			if cond1 && cond2 {

				return &faceMatchScore, nil
			}
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (f *FaceMatchScoreMockRepository) GetFaceMatchAPIReport(startDate time.Time, endDate time.Time) (map[uuid.UUID]repository.FaceMatchAPIReport, error) {

	return map[uuid.UUID]repository.FaceMatchAPIReport{}, nil
}
