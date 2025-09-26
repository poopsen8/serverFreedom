package servicePlan

import (
	"userServer/internal/models/modelPlan"
)

type repository interface {
	Get(id int64) (*modelPlan.Plan, error)
	GetAll() ([]*modelPlan.Plan, error)
}

type PlanService struct {
	repo repository
}

func NewPlanService(repo repository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) Get(id int64) (*modelPlan.Plan, error) {
	return s.repo.Get(id)
}

func (s *PlanService) GetAll() ([]*modelPlan.Plan, error) {
	return s.repo.GetAll()
}
