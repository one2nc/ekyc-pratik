package repository

import (
	"fmt"
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
	fmt.Println("email", email)
	customer := model.Customer{}
	result := c.dbInstance.Where("email = ?", email).First(&customer)
	return customer, result.Error
}

func newCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		dbInstance: db,
	}
}
