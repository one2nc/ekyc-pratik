package service

import (
	"go-ekyc/repository"

	"github.com/go-faker/faker/v4"
)

type OCRService struct {
	OCRRepository repository.IOCRRepository
}

func (o *OCRService) OCRExtractData() (map[string]interface{}, error) {

	user := struct {
		Address     faker.RealAddress `faker:"real_address"`
		Name        string            `faker:"name"`
		DateOfBirth string            `faker:"date"`
		IdNumber    string            `faker:"uuid_digit"`
		Gender      string            `faker:"gender"`
	}{}
	err := faker.FakeData(&user)

	userDetails := map[string]interface{}{}

	if err == nil {
		userDetails["name"] = user.Name
		userDetails["gender"] = user.Gender
		userDetails["dob"] = user.DateOfBirth
		userDetails["idNumber"] = user.IdNumber
		userDetails["address"] = user.Address.Address
		userDetails["pincode"] = user.Address.PostalCode

	}
	return userDetails, err
}

func newOCRService(ocrRepository repository.IOCRRepository) *OCRService {
	return &OCRService{
		OCRRepository: ocrRepository,
	}
}
