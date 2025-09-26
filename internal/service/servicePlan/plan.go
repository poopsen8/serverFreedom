package servicePlan

import (
	"userServer/internal/models/modelPlan"
)

type Repository interface {
	Get(id int64) (*modelPlan.Plan, error)
	GetAll() ([]*modelPlan.Plan, error)
}

type PlanService struct {
	repo Repository
}

func NewPlanService(repo Repository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) Get(id int64) (*modelPlan.Plan, error) {
	return s.repo.Get(id)
}

func (s *PlanService) GetAll() ([]*modelPlan.Plan, error) {
	return s.repo.GetAll()
}
