package repository

import (
	"go-ekyc/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlansRepository struct {
	dbInstance *gorm.DB
}

func (c *PlansRepository) FetchAllPlans() ([]model.Plan, error) {
	var plans []model.Plan
	result := c.dbInstance.Find(&plans)

	return plans, result.Error
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

func newPlansRepository(db *gorm.DB) *PlansRepository {
	return &PlansRepository{
		dbInstance: db,
	}
}
