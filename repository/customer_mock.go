package repository

import (
	"go-ekyc/model"

	"gorm.io/gorm"
)

type CustomerMockRepository struct {
	customerList []model.Customer
}

func newCustomerMockRepository(customer []model.Customer) ICustomerRepository {

	return &CustomerMockRepository{
		customerList: customer,
	}
}

func (c CustomerMockRepository) CreateCustomer(customer *model.Customer) error {
	c.customerList = append(c.customerList, *customer)

	return nil
}

func (c CustomerMockRepository) GetCustomerByEmail(email string) (model.Customer, error) {

	for _, customer := range c.customerList {
		if customer.Email == email {
			return customer, nil
		}
	}
	return model.Customer{}, gorm.ErrRecordNotFound
}
func (c CustomerMockRepository) GetCustomerByCredendials(accessKey string, secretKey string) (model.Customer, error) {

	for _, customer := range c.customerList {
		if customer.AccessKey == accessKey && customer.SecretKey == secretKey {
			return customer, nil
		}
	}
	return model.Customer{}, gorm.ErrRecordNotFound
}
func (c CustomerMockRepository) GetCustomersWithPlans() ([]model.Customer, error) {

	return []model.Customer{}, gorm.ErrRecordNotFound
}

