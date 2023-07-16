package repository

import (
	"errors"
	"go-ekyc/model"

	"github.com/google/uuid"
)

type PlansMockRepository struct {
	plans []model.Plan
}

func (c *PlansMockRepository) FetchPlansByName(name string) (model.Plan, error) {

	for _, plan := range c.plans {
		if plan.PlanName == name {
			return plan, nil

		}
	}
	return model.Plan{}, errors.New("Plan not found")
}

func (c *PlansMockRepository) FetchPlanById(id uuid.UUID) (model.Plan, error) {
	for _, plan := range c.plans {
		if plan.ID == id {
			return plan, nil

		}
	}
	return model.Plan{}, errors.New("Plan not found")
}

func newPlansMockRepository(plans []model.Plan) IPlansRepository {
	return &PlansMockRepository{
		plans: plans,
	}
}
