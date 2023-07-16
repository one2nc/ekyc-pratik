package service

import (
	"errors"
	"go-ekyc/helper"
	"go-ekyc/model"
	"go-ekyc/repository"

	"gorm.io/gorm"
)

type ICustomerService interface {
}
type RegisterServiceInput struct {
	PlanName      string
	CustomerEmail string
	CustomerName  string
}

type RegisterCustomerResult struct {
	AccessKey string
	SecretKey string
}

type CustomerService struct {
	customerRepository repository.ICustomerRepository
	plansRepository    repository.IPlansRepository
}

func (c *CustomerService) RegisterCustomer(serviceInput RegisterServiceInput) (RegisterCustomerResult, error) {
	plan, err := c.plansRepository.FetchPlansByName(serviceInput.PlanName)

	if err != nil {

		return RegisterCustomerResult{}, errors.New("error while fetching plan")
	}

	customer, err := c.customerRepository.GetCustomerByEmail(serviceInput.CustomerEmail)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {

		return RegisterCustomerResult{}, errors.New("error while fetching customer")

	}
	if customer != (model.Customer{}) {

		return RegisterCustomerResult{}, errors.New("email is already registered")
	}
	accessKey := helper.GenerateRandomString(10)
	secretKey := helper.GenerateRandomString(20)

	customerData := model.Customer{
		Name:      serviceInput.PlanName,
		Email:     serviceInput.CustomerEmail,
		PlanID:    plan.ID,
		AccessKey: helper.GetMD5Hash(accessKey),
		SecretKey: helper.GetMD5Hash(secretKey),
	}
	err = c.customerRepository.CreateCustomer(&customerData)

	if err != nil {

		return RegisterCustomerResult{}, err
	}

	return RegisterCustomerResult{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, err
}

func (c *CustomerService) GetCustomerByCredendials(accessKey string, secretKey string) (model.Customer, error) {

	customer, err := c.customerRepository.GetCustomerByCredendials(accessKey, secretKey)
	if err != nil {
		return customer, errors.New("error while fetching customers")

	}
	return customer, nil
}
func newCustomerService(customerRepository repository.ICustomerRepository, plansRepository repository.IPlansRepository) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
		plansRepository:    plansRepository,
	}
}
