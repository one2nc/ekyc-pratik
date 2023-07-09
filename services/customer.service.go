package service

import (
	"go-ekyc/model"
	"go-ekyc/repository"

)

type CustomerService struct {
	customerRepository *repository.CustomerRepository
}

func (c *CustomerService) CreateCustomer(customer *model.Customer) error {
	err := c.customerRepository.CreateCustomer(customer)
	return err
}

func (c *CustomerService) GetCustomerByEmail(email string) (model.Customer,error) {

	customer,err := c.customerRepository.GetCustomerByEmail(email)
	if err != nil {
		return customer,err

	}
	return customer , nil
}
func newCustomerService(customerRepository *repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
	}
}
