package repository

import (
	"go-ekyc/model"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	dbInstance *gorm.DB
}

func (c *CustomerRepository) CreateCustomer(customer *model.Customer) error {
	result := c.dbInstance.Create(customer)

	return result.Error
}

func (c *CustomerRepository) GetCustomerByEmail(email string) (model.Customer, error) {
	customer := model.Customer{}
	result := c.dbInstance.Where("email = ?", email).First(&customer)
	return customer, result.Error
}
func (c *CustomerRepository) GetCustomerByCredendials(accessKey string,secretKey string) (model.Customer, error) {
	customer := model.Customer{}
	result := c.dbInstance.Where("access_key = ? and secret_key = ? ", accessKey,secretKey).First(&customer)
	return customer, result.Error
}
func newCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		dbInstance: db,
	}
}
