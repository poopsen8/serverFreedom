package plan

import (
	"userServer/internal/model/plan"
)

type Repository interface {
	Plan(id int64) (*plan.Model, error)
	GetAll() ([]*plan.Model, error)
}

type PlanService struct {
	repo Repository
}

func NewPlanService(repo Repository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) Plan(id int64) (*plan.Model, error) {
	return s.repo.Plan(id)
}

func (s *PlanService) GetAll() ([]*plan.Model, error) {
	return s.repo.GetAll()
}
