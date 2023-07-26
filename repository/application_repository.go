package repository

import (
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	CustomerRepository       ICustomerRepository
	PlansRepository          IPlansRepository
	ImageRepository          IImageRepository
	FaceMatchScoreRepository IFaceMatchScoreRepository
	OCRRepository            IOCRRepository
	DailyReportsRepository   IDailyReportsRepository
	RedisRepository          RedisRepository
	CronRegistryRepository   ICronRegistryRepository
}

func NewApplicationRepository(db *gorm.DB) (*ApplicationRepository, error) {

	// redisConfig := config.NewRedisConfig()
	// redisRepository, err := newRedisRepository(redisConfig)

	// if err != nil {
	// 	return nil, err
	// }
	return &ApplicationRepository{
		CustomerRepository:       newCustomerRepository(db),
		PlansRepository:          newPlansRepository(db),
		ImageRepository:          newImageRepository(db),
		FaceMatchScoreRepository: newFaceMatchScoreRepository(db),
		OCRRepository:            newOCRRepositoty(db),
		DailyReportsRepository:   newDailyReportsRepository(db),
		CronRegistryRepository:   newCronRegistryRepository(db),
		// RedisRepository:          *redisRepository,
	}, nil
}
