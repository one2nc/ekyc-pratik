package service

import (
	"errors"
	"go-ekyc/helper"
	"go-ekyc/model"
	"go-ekyc/repository"
	"mime/multipart"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCustomerService_RegisterCustomer(t *testing.T) {
	defer TearDownTables(t)

	testEmail := "customer" + helper.GenerateRandomString(5) + "@gmail.com"

	testCases := []struct {
		serviceInput    RegisterServiceInput
		expectedOutput  RegisterCustomerResult
		isErrorExpected bool
		failedMessage   string
		testName        string
		expectedError   error
	}{

		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: "customer@gmail.com",
				PlanName:      "invalid plan name",
			},
			testName:        "Plan not found",
			isErrorExpected: true,
			failedMessage:   "Failed for Plan Not found test case",
			expectedError:   ErrPlanNotFound,
		},
		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: testEmail,
				PlanName:      "basic",
				CustomerName:  "test customer",
			},
			testName:        "Customer registerd successfully",
			isErrorExpected: false,
			failedMessage:   "Failed for Customer registerd successfully test case",
		},
		{
			serviceInput: RegisterServiceInput{
				CustomerEmail: testEmail,
				PlanName:      "basic",
			},
			testName:        "Customer exists test case",
			isErrorExpected: true,
			failedMessage:   "Failed for Customer exists test case",
			expectedError:   ErrEmailExists,
		},
	}

	for _, testCase := range testCases {
		result, err := appService.CustomerService.RegisterCustomer(testCase.serviceInput)

		if testCase.isErrorExpected {

			if err == nil {
				t.Fatal(testCase.failedMessage)
			} else {

				if !errors.Is(testCase.expectedError, err) {
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
	defer TearDownTables(t)
	testEmail := "customer" + helper.GenerateRandomString(5) + "@gmail.com"
	plan, err := appRepository.PlansRepository.FetchPlansByName("basic")
	if err != nil {
		t.Fatal(err.Error())
	}
	testCustomer := model.Customer{
		Name:      "test-customer",
		Email:     testEmail,
		PlanID:    plan.ID,
		AccessKey: helper.GenerateRandomString(10),
		SecretKey: helper.GenerateRandomString(10),
	}

	err = appRepository.CustomerRepository.CreateCustomer(&testCustomer)
	if err != nil {
		t.Fatal(err.Error())
	}
	testCases := []struct {
		accessKey       string
		secretKey       string
		expectedOutput  model.Customer
		isErrorExpected bool
		expectedError   error
		failedMessage   string
		testName        string
	}{
		{

			testName:        "Customer found",
			isErrorExpected: false,
			failedMessage:   "Failed for Customer found test case",
			accessKey:       testCustomer.AccessKey,
			secretKey:       testCustomer.SecretKey,
		},
		{

			testName:        "Customer not found",
			isErrorExpected: true,
			failedMessage:   "Failed for Customer not found test case",
			accessKey:       "wrong-key",
			secretKey:       "wrong-key",
			expectedError:   ErrCustomerNotFound,
		},
	}

	for _, testCase := range testCases {
		result, err := appService.CustomerService.GetCustomerByCredendials(testCase.accessKey, testCase.secretKey)

		if testCase.isErrorExpected {

			if err == nil {
				t.Fatal(testCase.failedMessage)
			} else {
				if !errors.Is(testCase.expectedError, err) {
					t.Fatal(testCase.failedMessage)
				}
			}
		} else {
			if err != nil {
				t.Fatal(testCase.failedMessage)
			} else {

				if result.ID.String() != testCustomer.ID.String() {
					t.Fatal(testCase.failedMessage)

				}
			}
		}

	}
}

func TestCustomerService_GetAggregateReportForCustomer(t *testing.T) {

	defer TearDownTables(t)
	testEmail := "customer" + helper.GenerateRandomString(5) + "@gmail.com"

	// fetching plan
	plan, err := appRepository.PlansRepository.FetchPlansByName("basic")
	if err != nil {
		t.Fatal(err.Error())
	}
	testCustomer := model.Customer{
		Name:      "test-customer",
		Email:     testEmail,
		PlanID:    plan.ID,
		AccessKey: helper.GenerateRandomString(10),
		SecretKey: helper.GenerateRandomString(10),
	}
	// creating test customer
	err = appRepository.CustomerRepository.CreateCustomer(&testCustomer)
	if err != nil {
		t.Fatal(err.Error())
	}

	startTime := testCustomer.CreatedAt.Add(time.Millisecond * 5)
	endTime := startTime.Add(time.Hour)
	// preparing test file data to create image record
	file, err := os.Open("./testdata/test_image_1.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err.Error())
	}

	fileHeader := &multipart.FileHeader{
		Filename: fileInfo.Name(),
		Size:     fileInfo.Size(),
	}
	testImages := []struct {
		Customer  model.Customer
		ImageType string
	}{
		{
			Customer:  testCustomer,
			ImageType: "face",
		},
		{
			Customer:  testCustomer,
			ImageType: "id_card",
		},
	}

	uploadedImageResluts := []ImageUploadResult{}
	// creating test image record
	for _, image := range testImages {
		result, err := appService.ImageService.UploadImage(UploadImageInput{
			Customer:  image.Customer,
			File:      file,
			FileInfo:  fileHeader,
			ImageType: image.ImageType,
		})

		if err != nil {
			t.Fatal(err.Error())
		}

		uploadedImageResluts = append(uploadedImageResluts, result)
	}

	// performing face match on test record to generate data for testing
	_, err = appService.ImageService.FaceMatch(FaceMatchInput{
		Customer: testCustomer,
		ImageId1: uploadedImageResluts[0].ImageId.String(),
		ImageId2: uploadedImageResluts[1].ImageId.String(),
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	// performing ocr on test record to generate data for testing

	_, err = appService.ImageService.GetOCRData(OCRInput{
		Customer: testCustomer,
		ImageId:  uploadedImageResluts[1].ImageId.String(),
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	// generating daily report
	err = appService.CustomerService.CreateCustomerReports(startTime, endTime, 1, 0)

	if err != nil {
		t.Fatal(err.Error())
	}
	// calculating report to compare it with resulte returned by service function
	totalImageUploadCost := float64(fileHeader.Size) / 1000 * 2 * plan.ImageUploadCost
	totalFaceMatchCost := 1 * plan.FaceMatchCost
	totalOCRCost := 1 * plan.OCRCost
	totalApiCharge := totalImageUploadCost + totalFaceMatchCost + totalOCRCost
	customerReports := repository.CustomerAggregatedReport{
		StartDate:           startTime,
		EndDate:             endTime,
		CustomerID:          testCustomer.ID,
		TotalBaseCharge:     plan.DailyBaseCost,
		TotalFaceMatchCount: 1,
		TotalFaceMatchCost:  totalFaceMatchCost,
		TotalOCRCount:       1,
		TotalOCRCost:        totalOCRCost,

		TotalImageStorageSizeMb: (float64(fileHeader.Size) * 2) / 1000,
		TotalImageStorageCost:   totalImageUploadCost,
		TotalInvoiceAmount:      plan.DailyBaseCost + totalApiCharge,
	}

	testCases := []struct {
		customerID        []uuid.UUID
		startDate         time.Time
		endDate           time.Time
		expectedOutput    repository.CustomerAggregatedReport
		isErrorExpected   bool
		failedMessage     string
		testName          string
		expectedError     error
		shouldMatchOutput bool
	}{
		{
			customerID: []uuid.UUID{
				testCustomer.ID,
			},
			expectedOutput:  customerReports,
			startDate:       startTime,
			endDate:         endTime,
			testName:        "Invoice number matched",
			isErrorExpected: false,
			failedMessage:   "Failed for Invoice number matched test case",
		},
	}

	for _, testCase := range testCases {
		result, err := appService.CustomerService.GetAggregateReportForCustomer(startTime, endTime, testCase.customerID)
		if testCase.isErrorExpected {

			if err == nil {
				t.Fatal(testCase.failedMessage)
			} else {
				if !errors.Is(testCase.expectedError, err) {
					t.Fatal(testCase.failedMessage)
				}
			}
		} else {
			if err != nil {
				t.Fatal(testCase.failedMessage)
			} else {

				if testCase.shouldMatchOutput {

					if !reflect.DeepEqual(result, testCase.expectedOutput) {
						t.Fatal(testCase.failedMessage)

					}
				} else {
					if result[0].TotalInvoiceAmount != testCase.expectedOutput.TotalInvoiceAmount {
						t.Fatal(testCase.failedMessage)

					}
				}
			}
		}

	}
}
