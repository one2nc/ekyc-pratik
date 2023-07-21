package mockedrepository

import (
	"go-ekyc/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var MockedPlans = []model.Plan{
	{
		ID:              uuid.UUID{},
		PlanName:        "basic",
		ImageUploadCost: 0.1,
		FaceMatchCost:   0.1,
		OCRCost:         0.1,
		DailyBaseCost:   10,
		IsActive:        true,
	},
	{
		ID:              uuid.UUID{},
		PlanName:        "advanced",
		ImageUploadCost: 0.01,
		FaceMatchCost:   0.01,
		OCRCost:         0.01,
		DailyBaseCost:   15,
		IsActive:        true,
	},
	{
		ID:              uuid.UUID{},
		PlanName:        "enterprise",
		ImageUploadCost: 0.05,
		FaceMatchCost:   0.05,
		OCRCost:         0.05,
		DailyBaseCost:   20,
		IsActive:        true,
	},
}
var MockedCustomers = []model.Customer{
	{
		ID:        uuid.New(),
		Email:     "customer1@gmail.com",
		PlanID:    MockedPlans[0].ID,
		Plan:      MockedPlans[0],
		Name:      "customer 1",
		AccessKey: "access-key-1",
		SecretKey: "Secret-key-1",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	{
		ID:        uuid.New(),
		Email:     "customer2@gmail.com",
		PlanID:    MockedPlans[0].ID,
		Plan:      MockedPlans[0],
		Name:      "customer 2",
		AccessKey: "access-key-2",
		SecretKey: "Secret-key-2",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
}

var MockedCustomerWithInvalidPlan = model.Customer{
	ID:        uuid.New(),
	Email:     "customer@gmail.com",
	PlanID:    uuid.New(),
	Name:      "customer",
	AccessKey: "access-key-1",
	SecretKey: "Secret-key-1",
	CreatedAt: time.Time{},
	UpdatedAt: time.Time{},
}

var CreateCustomerData = model.Customer{
	ID:        uuid.New(),
	Email:     "createCustomer@gmail.com",
	PlanID:    MockedPlans[0].ID,
	Name:      "create customer",
	AccessKey: "access-key-for-created-customer",
	SecretKey: "access-key-for-created-customer",
	CreatedAt: time.Time{},
	UpdatedAt: time.Time{},
}

var MockedImageData = []model.Image{
	{
		ID:            uuid.New(),
		CustomerID:    MockedCustomers[0].ID,
		FilePath:      "test/file/path.jpeg",
		FileExtension: "jpeg",
		FileSizeMB:    1000,
		ImageType:     "face",
	},
	{
		ID:            uuid.New(),
		CustomerID:    MockedCustomers[0].ID,
		FilePath:      "test/file/path.jpeg",
		FileExtension: "jpeg",
		FileSizeMB:    1000,
		ImageType:     "id_card",
	},
	{
		ID:            uuid.New(),
		CustomerID:    MockedCustomers[1].ID,
		FilePath:      "test/file/path.jpeg",
		FileExtension: "jpeg",
		FileSizeMB:    1000,
		ImageType:     "face",
	},
}

var MockedImageUploadApiCalls = []model.ImageUploadAPICall{
	{
		ID:                  uuid.New(),
		CustomerID:          MockedImageData[0].ID,
		ImageID:             MockedImageData[0].ID,
		ImageStorageCharges: MockedCustomers[0].Plan.ImageUploadCost * MockedImageData[0].FileSizeMB,
	},
	{
		ID:                  uuid.New(),
		CustomerID:          MockedImageData[1].ID,
		ImageID:             MockedImageData[1].ID,
		ImageStorageCharges: MockedCustomers[1].Plan.ImageUploadCost * MockedImageData[0].FileSizeMB,
	},
	{
		ID:                  uuid.New(),
		CustomerID:          MockedImageData[2].ID,
		ImageID:             MockedImageData[2].ID,
		ImageStorageCharges: MockedCustomers[1].Plan.ImageUploadCost * MockedImageData[0].FileSizeMB,
	},
}

var MockedOCRData = []model.OCRData{
	{
		ID:         uuid.New(),
		CustomerID: MockedImageData[1].CustomerID,
		ImageID:    MockedImageData[1].ID,
		OCRData:    datatypes.JSON(""),
	},
}
var MockedOCRAPICalls = []model.OCRAPICalls{
	{
		ID:             uuid.New(),
		CustomerID:     MockedOCRData[0].CustomerID,
		ImageID:        MockedOCRData[0].ImageID,
		OCRID:          MockedOCRData[0].ID,
		APICallCharges: MockedPlans[0].OCRCost,
	},
	{
		ID:             uuid.New(),
		CustomerID:     MockedOCRData[0].CustomerID,
		ImageID:        MockedOCRData[0].ImageID,
		OCRID:          MockedOCRData[0].ID,
		APICallCharges: MockedPlans[0].OCRCost,
	},
}

var MockedFaceMatchScoreData = []model.FaceMatchScore{
	{
		ID:         uuid.New(),
		CustomerID: MockedImageData[0].CustomerID,
		ImageID1:   MockedImageData[0].ID,
		ImageID2:   MockedImageData[1].ID,
		Score:      70,
	},
}

var MockedFaceMatchAPICalls = []model.FaceMatchAPICall{
	{
		ID: uuid.New(),
		CustomerID: MockedFaceMatchScoreData[0].CustomerID,
		ScoreID: MockedFaceMatchScoreData[0].ID,
		APICallCharges: MockedPlans[0].FaceMatchCost ,
	},
	{
		ID: uuid.New(),
		CustomerID: MockedFaceMatchScoreData[0].CustomerID,
		ScoreID: MockedFaceMatchScoreData[0].ID,
		APICallCharges: MockedPlans[0].FaceMatchCost ,
	},
}