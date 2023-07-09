package service

import (
	"go-ekyc/model"
	"go-ekyc/repository"
)

type PlansService struct {
	plansRepository *repository.PlansRepository
}

func (p *PlansService) FetchAllPlans() ([]model.Plan,error) {
	plans,err := p.plansRepository.FetchAllPlans()
	return plans ,err
}
func (p *PlansService) FetchPlansByName(name string) (model.Plan,error) {
	plan,err := p.plansRepository.FetchPlansByName(name)
	
	return plan ,err
}


func newPlansService(plansRepository *repository.PlansRepository) *PlansService {
	return &PlansService{
		plansRepository: plansRepository,
	}
}
