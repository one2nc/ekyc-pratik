package repository

import "go-ekyc/model"

type ApplicationMockRepository struct {
	CustomerRepository ICustomerRepository
	PlansRepository IPlansRepository
	ImageRepository IImageRepository
	FaceMatchScoreRepository IFaceMatchScoreRepository
	OCRRepository IOCRRepository
	DailyReportsRepository IDailyReportsRepository
	RedisRepository RedisRepository
}

func NewApplicationMockRepository(customerData []model.Customer,plans []model.Plan) (*ApplicationMockRepository) {

	return &ApplicationMockRepository{
		CustomerRepository: newCustomerMockRepository(customerData),
		PlansRepository: newPlansMockRepository(plans),
	}
}
