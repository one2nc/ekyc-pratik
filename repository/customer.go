package repository

import (
	"go-ekyc/model"
	"time"

	"gorm.io/gorm"
)

type ICustomerRepository interface {
	CreateCustomer(customer *model.Customer) error
	GetCustomerByEmail(email string) (model.Customer, error)
	GetCustomerByCredendials(accessKey string, secretKey string) (model.Customer, error)
	GetCustomersWithPlans(limit int, offset int, customerBeforeDate time.Time) ([]model.Customer, error)
}

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
func (c *CustomerRepository) GetCustomerByCredendials(accessKey string, secretKey string) (model.Customer, error) {
	customer := model.Customer{}
	result := c.dbInstance.Where("access_key = ? and secret_key = ? ", accessKey, secretKey).First(&customer)
	return customer, result.Error
}

func (c *CustomerRepository) GetCustomersWithPlans(limit int, offset int, customerBeforeDate time.Time) ([]model.Customer, error) {
	var customers []model.Customer
	result := c.dbInstance.Preload("Plan").Where("customers.created_at < ?", customerBeforeDate).Offset(offset).Limit(limit).Order("id").Find(&customers)

	if result.Error != nil {
		return customers, result.Error
	}
	return customers, result.Error
}
func newCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{
		dbInstance: db,
	}
}
