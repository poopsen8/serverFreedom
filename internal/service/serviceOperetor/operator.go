package serviceOperetor

import (
	"userServer/internal/models/modelOperator"
)

type repository interface {
	Get(id int64) (*modelOperator.Operator, error)
	GetAll() ([]*modelOperator.Operator, error)
}

type PlanService struct {
	repo repository
}

func NewOperatorService(repo repository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) Get(id int64) (*modelOperator.Operator, error) {
	return s.repo.Get(id)
}

func (s *PlanService) GetAll() ([]*modelOperator.Operator, error) {
	return s.repo.GetAll()
}
