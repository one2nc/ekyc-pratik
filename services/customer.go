package service

import (
	"errors"
	"fmt"
	"go-ekyc/helper"
	"go-ekyc/model"
	"go-ekyc/repository"
	service_inputs "go-ekyc/services/service-input"
	service_results "go-ekyc/services/service-results"

	"gorm.io/gorm"
)

type CustomerService struct {
	customerRepository *repository.CustomerRepository
	plansRepository    *repository.PlansRepository
}

func (c *CustomerService) CreateCustomer(customer *model.Customer) error {
	err := c.customerRepository.CreateCustomer(customer)
	return err
}

func (c *CustomerService) RegisterCustomer(serviceInput service_inputs.RegisterServiceInput) (service_results.RegisterCustomerResult, error) {
	fmt.Println("fetching plan")
	plan, err := c.plansRepository.FetchPlansByName(serviceInput.PlanName)
	fmt.Println("fetched plan")

	if err != nil {

		return service_results.RegisterCustomerResult{}, err
	}

	customer, err := c.customerRepository.GetCustomerByEmail(serviceInput.CustomerEmail)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {

		return service_results.RegisterCustomerResult{}, err

	}
	if customer != (model.Customer{}) {

		return service_results.RegisterCustomerResult{}, errors.New("Email is already registered")
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

		return service_results.RegisterCustomerResult{}, err
	}

	return service_results.RegisterCustomerResult{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, err
}

func (c *CustomerService) GetCustomerByEmail(email string) (model.Customer, error) {

	customer, err := c.customerRepository.GetCustomerByEmail(email)
	if err != nil {
		return customer, err

	}
	return customer, nil
}
func (c *CustomerService) GetCustomerByCredendials(accessKey string, secretKey string) (model.Customer, error) {

	customer, err := c.customerRepository.GetCustomerByCredendials(accessKey, secretKey)
	if err != nil {
		return customer, err

	}
	return customer, nil
}
func newCustomerService(customerRepository *repository.CustomerRepository, plansRepository *repository.PlansRepository) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
		plansRepository: plansRepository,
	}
}
