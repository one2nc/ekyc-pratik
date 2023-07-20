package service

import (
	"errors"
	"fmt"
	"go-ekyc/model"
	mockedrepository "go-ekyc/repository/mocked-repository"
	"reflect"
	"testing"
)

func TestCustomerService_RegisterCustomer(t *testing.T) {

	appMockRepository := mockedrepository.NewApplicationMockRepository(mockedrepository.MockedCustomers, mockedrepository.MockedPlans)
	customerService := newCustomerService(appMockRepository.CustomerRepository, appMockRepository.PlansRepository, appMockRepository.ImageRepository, appMockRepository.FaceMatchScoreRepository, appMockRepository.OCRRepository, appMockRepository.DailyReportsRepository, appMockRepository.RedisRepository)

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
				CustomerEmail: mockedrepository.MockedCustomers[0].Email,
				PlanName:      mockedrepository.MockedPlans[0].PlanName,
			},
			decription:      "Customer exists test case",
			isErrorExpected: true,
			failedMessage:   "Failed for Customer exists test case",
			expectedError:   errors.New("email is already registered"),
		},
		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: mockedrepository.MockedCustomers[0].Email,
				PlanName:      "invalid plan name",
			},
			decription:      "Plan not found",
			isErrorExpected: true,
			failedMessage:   "Failed for Plan Not found test case",
			expectedError:   errors.New("error while fetching plan"),
		},
		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: mockedrepository.CreateCustomerData.Email,
				PlanName:      mockedrepository.MockedPlans[0].PlanName,
				CustomerName:  mockedrepository.CreateCustomerData.Name,
			},
			decription:      "Customer registerd successfully",
			isErrorExpected: false,
			expectedOutput:  RegisterCustomerResult{AccessKey: mockedrepository.CreateCustomerData.AccessKey, SecretKey: mockedrepository.CreateCustomerData.SecretKey},
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

	appMockRepository := mockedrepository.NewApplicationMockRepository(mockedrepository.MockedCustomers, mockedrepository.MockedPlans)
	customerService := newCustomerService(appMockRepository.CustomerRepository, appMockRepository.PlansRepository, appMockRepository.ImageRepository, appMockRepository.FaceMatchScoreRepository, appMockRepository.OCRRepository, appMockRepository.DailyReportsRepository, appMockRepository.RedisRepository)

	testCases := []struct {
		accessKey       string
		secretKey       string
		expectedOutput  model.Customer
		isErrorExpected bool
		expectedError   error
		failedMessage   string
		decription      string
	}{
		{

			decription:      "Customer found",
			isErrorExpected: false,
			failedMessage:   "Failed for Customer found test case",
			accessKey:       mockedrepository.MockedCustomers[0].AccessKey,
			secretKey:       mockedrepository.MockedCustomers[0].SecretKey,
			expectedOutput:  mockedrepository.MockedCustomers[0],
		},
		{

			decription:      "Customer not found",
			isErrorExpected: true,
			failedMessage:   "Failed for Customer not found test case",
			accessKey:       "wrong-key",
			secretKey:       "wrong-key",
			expectedError:   errors.New("error while fetching customer"),
		},
	}

	for _, testCase := range testCases {
		result, err := customerService.GetCustomerByCredendials(testCase.accessKey, testCase.secretKey)

		if testCase.isErrorExpected {

			if err == nil {
				t.Fatal(testCase.failedMessage)
			} else {
				fmt.Println(testCase.expectedError.Error())
				fmt.Println(err.Error())
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
