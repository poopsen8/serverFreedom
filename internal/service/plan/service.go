package plan

import (
	"sort"
	"userServer/internal/model/plan"
)

type Repository interface {
	Plan(id int64) (*plan.Model, error)
	Plans() ([]*plan.Model, error)
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

func (s *PlanService) Plans() ([]*plan.Model, error) {
	plans, err := s.repo.Plans()
	if err != nil {
		return nil, err
	}

	// Сортировка по возрастанию Duration
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].Duration < plans[j].Duration
	})

	return plans, nil
}
