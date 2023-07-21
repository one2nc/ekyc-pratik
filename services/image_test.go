package service

import (
	"errors"
	"fmt"
	mockedrepository "go-ekyc/repository/mocked-repository"
	"mime/multipart"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestImageService_UploadImage(t *testing.T) {

	file, err := os.Open("./../testdata/test_image_1.png")
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

	appMockRepository := mockedrepository.NewApplicationMockRepository(mockedrepository.MockedCustomers, mockedrepository.MockedPlans, mockedrepository.MockedImageData, mockedrepository.MockedImageUploadApiCalls)
	ocrService := newOCRService(appMockRepository.OCRRepository)
	minioMockService, _ := NewMinioMockService()
	imageService := newImageService(appMockRepository.ImageRepository, appMockRepository.PlansRepository, appMockRepository.FaceMatchScoreRepository, appMockRepository.OCRRepository, ocrService, minioMockService)

	testCases := []struct {
		serviceInput    UploadImageInput
		expectedOutput  ImageUploadResult
		isErrorExpected bool
		failedMessage   string
		decription      string
		expectedError   error
	}{
		{
			serviceInput: UploadImageInput{
				Customer:  mockedrepository.MockedCustomers[0],
				File:      file,
				FileInfo:  fileHeader,
				ImageType: "face",
			},
			decription:      "Image Created",
			isErrorExpected: false,
			failedMessage:   "Failed for Image Created test case",
		},
		{
			serviceInput: UploadImageInput{
				Customer:  mockedrepository.MockedCustomerWithInvalidPlan,
				File:      file,
				FileInfo:  fileHeader,
				ImageType: "face",
			},
			decription:      "Invalid Plan",
			isErrorExpected: true,
			failedMessage:   "Failed for Invalid Plan test case",
			expectedError: errors.New("Plan not found"),
		},
	}

	for _, testCase := range testCases {
		result, err := imageService.UploadImage(testCase.serviceInput)
		fmt.Println(result, err)
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
				if result.ImageId == uuid.Nil {
					t.Fatal(testCase.failedMessage)

				}
			}
		}

	}
}
