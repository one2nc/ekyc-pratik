package repository

import (
	"go-ekyc/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type IPlansRepository interface {
	FetchPlansByName(name string) (model.Plan, error)
	FetchPlanById(id uuid.UUID) (model.Plan, error)
}

type PlansRepository struct {
	dbInstance *gorm.DB
}


func (c *PlansRepository) FetchPlansByName(name string) (model.Plan, error) {

	var plan model.Plan

	result := c.dbInstance.Where("plan_name = ?", name).First(&plan)
	return plan, result.Error
}
func (c *PlansRepository) FetchPlanById(id uuid.UUID) (model.Plan, error) {
	var plan model.Plan
	result := c.dbInstance.Where("id = ?", id).First(&plan)

	return plan, result.Error
}

func newPlansRepository(db *gorm.DB) IPlansRepository {
	return &PlansRepository{
		dbInstance: db,
	}
}
