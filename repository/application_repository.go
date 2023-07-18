package repository

import (
	"go-ekyc/config"
	"go-ekyc/db"
)

type ApplicationRepository struct {
	CustomerRepository       ICustomerRepository
	PlansRepository          IPlansRepository
	ImageRepository          IImageRepository
	FaceMatchScoreRepository IFaceMatchScoreRepository
	OCRRepository            IOCRRepository
	DailyReportsRepository   IDailyReportsRepository
	RedisRepository          RedisRepository
}

func NewApplicationRepository() (*ApplicationRepository, error) {
	db, err := db.InitiateDB()
	if err != nil {
		return nil, err
	}

	redisConfig := config.NewRedisConfig()
	redisRepository, err := newRedisRepository(redisConfig)

	if err != nil {
		return nil, err
	}
	return &ApplicationRepository{
		CustomerRepository:       newCustomerRepository(db),
		PlansRepository:          newPlansRepository(db),
		ImageRepository:          newImageRepository(db),
		FaceMatchScoreRepository: newFaceMatchScoreRepository(db),
		OCRRepository:            newOCRRepositoty(db),
		DailyReportsRepository:   newDailyReportsRepository(db),
		RedisRepository:          *redisRepository,
	}, nil
}
