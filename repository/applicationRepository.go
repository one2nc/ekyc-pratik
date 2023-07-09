package repository

import (
	"go-ekyc/db"
)

type ApplicationRepository struct {
	CustomerRepository *CustomerRepository
	PlansRepository    *PlansRepository
}

func NewApplicationRepository() (*ApplicationRepository, error) {
	db, err := db.InitiateDB()
	if err != nil {
		return nil, err
	}

	return &ApplicationRepository{
		CustomerRepository: newCustomerRepository(db),
		PlansRepository:    newPlansRepository(db),
	}, nil
}
