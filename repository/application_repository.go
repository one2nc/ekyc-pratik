package repository

import (
	"go-ekyc/db"
)

type ApplicationRepository struct {
	CustomerRepository       *CustomerRepository
	PlansRepository          *PlansRepository
	ImageRepository          *ImageRepository
	FaceMatchScoreRepository *FaceMatchScoreRepository
}

func NewApplicationRepository() (*ApplicationRepository, error) {
	db, err := db.InitiateDB()
	if err != nil {
		return nil, err
	}

	return &ApplicationRepository{
		CustomerRepository:       newCustomerRepository(db),
		PlansRepository:          newPlansRepository(db),
		ImageRepository:          newImageRepository(db),
		FaceMatchScoreRepository: newFaceMatchScoreRepository(db),
	}, nil
}
