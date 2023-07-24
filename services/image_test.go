package service

import (
	"encoding/json"
	"errors"
	"go-ekyc/helper"
	"go-ekyc/model"
	"mime/multipart"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func TestImageService_UploadImage(t *testing.T) {
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

	file, err := os.Open("./testdata/test_image_1.png")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		t.Fatal(err.Error())
	}

	// Create a *multipart.FileHeader using the file name and size
	fileHeader := &multipart.FileHeader{
		Filename: fileInfo.Name(),
		Size:     fileInfo.Size(),
	}

	testCases := []struct {
		serviceInput    UploadImageInput
		expectedOutput  ImageUploadResult
		isErrorExpected bool
		failedMessage   string
		testName        string
		expectedError   error
	}{
		{
			serviceInput: UploadImageInput{
				Customer:  testCustomer,
				File:      file,
				FileInfo:  fileHeader,
				ImageType: "face",
			},
			testName:        "Image Created",
			isErrorExpected: false,
			failedMessage:   "Failed for Image Created test case",
		},
		{
			serviceInput: UploadImageInput{
				Customer: model.Customer{
					PlanID: uuid.New(),
				},
				File:      file,
				FileInfo:  fileHeader,
				ImageType: "face",
			},
			testName:        "Invalid Plan",
			isErrorExpected: true,
			failedMessage:   "Failed for Invalid Plan test case",
			expectedError:   ErrPlanNotFound,
		},
	}

	for _, testCase := range testCases {
		result, err := appService.ImageService.UploadImage(testCase.serviceInput)

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
				if result.ImageId == uuid.Nil {
					t.Fatal(testCase.failedMessage)

				}
			}
		}

	}
}

func TestImageService_FaceMatch(t *testing.T) {
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

	testImages := []model.Image{
		{
			CustomerID:    testCustomer.ID,
			FilePath:      "test-path/image-1",
			FileExtension: ".jpeg",
			FileSizeMB:    1000,
			ImageType:     "face",
		},
		{
			CustomerID:    testCustomer.ID,
			FilePath:      "test-path/image-2",
			FileExtension: ".jpeg",
			FileSizeMB:    1000,
			ImageType:     "id_card",
		},
		{
			CustomerID:    testCustomer.ID,
			FilePath:      "test-path/image-3",
			FileExtension: ".jpeg",
			FileSizeMB:    1000,
			ImageType:     "id_card",
		},
	}

	err = appRepository.ImageRepository.CreateBulkImage(&testImages)
	if err != nil {
		t.Fatal(err.Error())
	}
	testFaceMatchScore := model.FaceMatchScore{
		CustomerID: testCustomer.ID,
		ImageID1:   testImages[0].ID,
		ImageID2:   testImages[1].ID,
		Score:      50,
	}

	err = appRepository.FaceMatchScoreRepository.CreateFaceMatchScore(&testFaceMatchScore)
	if err != nil {
		t.Fatal(err.Error())
	}

	testCases := []struct {
		serviceInput      FaceMatchInput
		expectedOutput    FaceMatchResult
		isErrorExpected   bool
		failedMessage     string
		testName          string
		expectedError     error
		shouldMatchOutput bool
	}{
		{
			serviceInput: FaceMatchInput{
				Customer: testCustomer,
				ImageId1: testImages[0].ID.String(),
				ImageId2: testImages[2].ID.String(),
			},

			testName:        "New Face match score generated",
			isErrorExpected: false,
			failedMessage:   "Failed for New Face match score generated test case",
		},
		{
			serviceInput: FaceMatchInput{
				Customer: testCustomer,
				ImageId1: testImages[0].ID.String(),
				ImageId2: testImages[1].ID.String(),
			},

			shouldMatchOutput: true,

			expectedOutput:  FaceMatchResult{Score: testFaceMatchScore.Score},
			testName:        "Old Face match score generated",
			isErrorExpected: false,
			failedMessage:   "Failed for Old Face match score generated test case",
		},
		{
			serviceInput: FaceMatchInput{
				Customer: testCustomer,
				ImageId1: uuid.New().String(), // incorrect uuid
				ImageId2: testImages[1].ID.String(),
			},

			expectedError:   ErrImageNotFound,
			testName:        "Image not found",
			isErrorExpected: true,
			failedMessage:   "Failed for Image not found test case",
		},
		{
			serviceInput: FaceMatchInput{
				Customer: testCustomer,
				ImageId1: testImages[1].ID.String(),
				ImageId2: testImages[2].ID.String(),
			},

			expectedError:   ErrInvalidImageType,
			testName:        "Invalid image type",
			isErrorExpected: true,
			failedMessage:   "Failed for Invalid image type test case",
		},
		{
			serviceInput: FaceMatchInput{
				Customer: model.Customer{
					PlanID: uuid.New(),
				},
			},

			expectedError:   ErrPlanNotFound,
			testName:        "plan not fond",
			isErrorExpected: true,
			failedMessage:   "Failed for plan not found test case",
		},
	}

	for _, testCase := range testCases {
		result, err := appService.ImageService.FaceMatch(testCase.serviceInput)

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

				if testCase.shouldMatchOutput {

					if !reflect.DeepEqual(result, testCase.expectedOutput) {
						t.Fatal(testCase.failedMessage)

					}
				} else {
					if result.Score == 0 {
						t.Fatal(testCase.failedMessage)

					}
				}
			}
		}

	}
}
func TestImageService_GetOCRData(t *testing.T) {
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

	testImages := []model.Image{
		{
			CustomerID:    testCustomer.ID,
			FilePath:      "test-path/image-1",
			FileExtension: ".jpeg",
			FileSizeMB:    1000,
			ImageType:     "face",
		},
		{
			CustomerID:    testCustomer.ID,
			FilePath:      "test-path/image-2",
			FileExtension: ".jpeg",
			FileSizeMB:    1000,
			ImageType:     "id_card",
		},
		{
			CustomerID:    testCustomer.ID,
			FilePath:      "test-path/image-3",
			FileExtension: ".jpeg",
			FileSizeMB:    1000,
			ImageType:     "id_card",
		},
	}

	err = appRepository.ImageRepository.CreateBulkImage(&testImages)
	if err != nil {
		t.Fatal(err.Error())
	}

	ocrTestData, err := appService.OCRService.OCRExtractData()
	if err != nil {
		t.Fatal(err.Error())
	}
	jsonData, err := json.Marshal(ocrTestData)
	if err != nil {
		t.Fatal(err.Error())
	}
	ocrData := &model.OCRData{
		CustomerID: testCustomer.ID,
		ImageID:    testImages[2].ID,
		OCRData:    datatypes.JSON(jsonData),
	}

	err = appRepository.OCRRepository.CreateOCRData(ocrData)
	if err != nil {
		t.Fatal(err.Error())
	}

	testCases := []struct {
		serviceInput      OCRInput
		expectedOutput    OCRResult
		isErrorExpected   bool
		failedMessage     string
		testName          string
		expectedError     error
		shouldMatchOutput bool
	}{
		{
			serviceInput: OCRInput{
				Customer: testCustomer,
				ImageId:  testImages[1].ID.String(),
			},

			testName:          "New OCR data generated",
			isErrorExpected:   false,
			failedMessage:     "Failed for new ocr data generated test case",
			shouldMatchOutput: false,
		},
		{
			serviceInput: OCRInput{
				Customer: testCustomer,
				ImageId:  testImages[2].ID.String(),
			},

			testName:          "Old OCR data fetched",
			isErrorExpected:   false,
			failedMessage:     "Failed for Old OCR data fetched test case",
			shouldMatchOutput: true,
			expectedOutput: OCRResult{
				Data: ocrData.OCRData,
			},
		},
		{
			serviceInput: OCRInput{
				Customer: testCustomer,
				ImageId:  testImages[0].ID.String(),
			},

			testName:        "Incorrect image type",
			isErrorExpected: true,
			failedMessage:   "Failed for incorrect image type test case",
			expectedError:   ErrInvalidImageType,
		},
	}

	for _, testCase := range testCases {
		result, err := appService.ImageService.GetOCRData(testCase.serviceInput)

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

					var data1 map[string]interface{}
					var data2 map[string]interface{}

					if err := json.Unmarshal([]byte(testCase.expectedOutput.Data), &data1); err != nil {
						t.Fatal(err.Error())
					}

					if err := json.Unmarshal([]byte(result.Data), &data2); err != nil {
						t.Fatal(err.Error())
					}

					if !reflect.DeepEqual(data1, data2) {
						t.Fatal(testCase.failedMessage)
					}
				} else {
					if result.Data == nil {
						t.Fatal(testCase.failedMessage)

					}
				}
			}
		}

	}
}
