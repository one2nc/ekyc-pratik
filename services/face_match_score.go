package service

import (
	"go-ekyc/model"
	"go-ekyc/repository"
)

type FaceMatchScoreService struct {
	faceMatchScoreRepository *repository.FaceMatchScoreRepository
}

func (i *FaceMatchScoreService) CreateFaceMatchScore(faceScoreData *model.FaceMatchScore) error {
	err := i.faceMatchScoreRepository.CreateFaceMatchScore(faceScoreData)
	return err
}

func (f *FaceMatchScoreService) FetchScoreByImageAndCustomerId(imageId1 string, imageId2 string, customerId string) (*model.FaceMatchScore, error) {

	faceImageScore, err := f.faceMatchScoreRepository.FetchScoreByImageAndCustomerId(imageId1, imageId2, customerId)
	return faceImageScore, err
}


func (f *FaceMatchScoreService) CreateFaceMatchScoreAPIRecord(faceMatchApiData *model.FaceMatchAPICall) error {
	err := f.faceMatchScoreRepository.CreateFaceMatchScoreAPIRecord(faceMatchApiData)

	return err
}

func newFaceMatchScoreService(faceMatchScoreRepository *repository.FaceMatchScoreRepository) *FaceMatchScoreService {
	return &FaceMatchScoreService{
		faceMatchScoreRepository: faceMatchScoreRepository,
	}
}
