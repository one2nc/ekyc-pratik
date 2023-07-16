package service

import (
	"errors"
	"go-ekyc/model"
	"go-ekyc/repository"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCustomerService_RegisterCustomer(t *testing.T) {
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

	var createCustomerData = model.Customer{
		ID:        uuid.New(),
		Email:     "createCustomer@gmail.com",
		PlanID:    MockedPlans[0].ID,
		Name:      "create customer",
		AccessKey: "access-key-for-created-customer",
		SecretKey: "access-key-for-created-customer",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	appMockRepository := repository.NewApplicationMockRepository(MockedCustomers, MockedPlans)
	customerService := newCustomerService(appMockRepository.CustomerRepository, appMockRepository.PlansRepository)

	testCases := []struct {
		serviceInput    RegisterServiceInput
		expectedOutput  RegisterCustomerResult
		isErrorExpected bool
		failedMessage   string
		decription      string
		expectedError   error
	}{
		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: MockedCustomers[0].Email,
				PlanName:      MockedPlans[0].PlanName,
			},
			decription:      "Customer exists test case",
			isErrorExpected: true,
			failedMessage:   "Failed for Customer exists test case",
			expectedError:   errors.New("email is already registered"),
		},
		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: MockedCustomers[0].Email,
				PlanName:     "invalid plan name",
			},
			decription:      "Plan not found",
			isErrorExpected: true,
			failedMessage:   "Failed for Plan Not found test case",
			expectedError:   errors.New("error while fetching plan"),
		},
		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: createCustomerData.Email,
				PlanName:      MockedPlans[0].PlanName,
				CustomerName:  createCustomerData.Name,
			},
			decription:      "Customer registerd successfully",
			isErrorExpected: false,
			expectedOutput:  RegisterCustomerResult{AccessKey: createCustomerData.AccessKey, SecretKey: createCustomerData.SecretKey},
			failedMessage:   "Failed for Customer registerd successfully test case",
		},
	}

	for _, testCase := range testCases {
		result, err := customerService.RegisterCustomer(testCase.serviceInput)

		if testCase.isErrorExpected {

			if err == nil {
				t.Fatal(testCase.failedMessage)
			} else {
				if !reflect.DeepEqual(testCase.expectedError, err) {
					t.Fatal(testCase.failedMessage)
				}
			}
		} else {
			if err != nil {
				t.Fatal(testCase.failedMessage)
			} else {
				if result.AccessKey == "" || result.SecretKey == "" {
					t.Fatal(testCase.failedMessage)

				}
			}
		}

	}
}
func TestCustomerService_GetCustomerByCredendials(t *testing.T) {
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

	appMockRepository := repository.NewApplicationMockRepository(MockedCustomers, MockedPlans)
	customerService := newCustomerService(appMockRepository.CustomerRepository, appMockRepository.PlansRepository)

	testCases := []struct {
		accessKey       string
		secretKey       string
		expectedOutput  model.Customer
		isErrorExpected bool
		expectedError error
		failedMessage   string
		decription      string
	}{
		{

			decription:      "Customer found",
			isErrorExpected: false,
			failedMessage:   "Failed for Customer found test case",
			accessKey:       MockedCustomers[0].AccessKey,
			secretKey:       MockedCustomers[0].SecretKey,
			expectedOutput:  MockedCustomers[0],
		},
		{

			decription:      "Customer not found",
			isErrorExpected: true,
			failedMessage:   "Failed for Customer not found test case",
			accessKey:       "wrong-key",
			secretKey:       "wrong-key",
			expectedError: errors.New("error while fetching customers"),
		},
	}

	for _, testCase := range testCases {
		result, err := customerService.GetCustomerByCredendials(testCase.accessKey,testCase.secretKey)

		if testCase.isErrorExpected {

			if err == nil {
				t.Fatal(testCase.failedMessage)
			} else {
				if !reflect.DeepEqual(testCase.expectedError, err) {
					t.Fatal(testCase.failedMessage)
				}
			}
		} else {
			if err != nil {
				t.Fatal(testCase.failedMessage)
			} else {
				if result.AccessKey == "" || result.SecretKey == "" {
					t.Fatal(testCase.failedMessage)

				}
			}
		}

	}
}
