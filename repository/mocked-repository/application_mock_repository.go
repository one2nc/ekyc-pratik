package mockedrepository

import (
	"go-ekyc/model"
	"go-ekyc/repository"
)

type ApplicationMockRepository struct {
	CustomerRepository repository.ICustomerRepository
	PlansRepository repository.IPlansRepository
	ImageRepository repository.IImageRepository
	FaceMatchScoreRepository repository.IFaceMatchScoreRepository
	OCRRepository repository.IOCRRepository
	DailyReportsRepository repository.IDailyReportsRepository
	RedisRepository repository.RedisRepository
}

func NewApplicationMockRepository(customerData []model.Customer,plans []model.Plan,images []model.Image,imageUploadApiCalls []model.ImageUploadAPICall) (*ApplicationMockRepository) {

	return &ApplicationMockRepository{
		CustomerRepository: newCustomerMockRepository(customerData),
		PlansRepository: newPlansMockRepository(plans),
		ImageRepository: newImageMockRepository(images,imageUploadApiCalls),
	}
}
