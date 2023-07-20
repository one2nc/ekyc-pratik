package mockedrepository

import (
	"go-ekyc/model"
	"time"

	"github.com/google/uuid"
)

var MockedPlans = []model.Plan{
	{
		ID:       uuid.UUID{},
		PlanName: "basic",
	},
	{
		ID:       uuid.UUID{},
		PlanName: "advanced",
	},
	{
		ID:       uuid.UUID{},
		PlanName: "enterprise",
	},
}
var MockedCustomers = []model.Customer{
	{
		ID:        uuid.New(),
		Email:     "customer1@gmail.com",
		PlanID:    MockedPlans[0].ID,
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
		Name:      "customer 2",
		AccessKey: "access-key-2",
		SecretKey: "Secret-key-2",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
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
