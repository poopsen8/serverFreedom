package operetor

import (
	"userServer/internal/model/operator"
)

type repository interface {
	Operator(id int64) (*operator.Model, error)
	Operators() ([]*operator.Model, error)
}

type OperetorService struct {
	repo repository
}

func NewOperatorService(repo repository) *OperetorService {
	return &OperetorService{repo: repo}
}

func (s *OperetorService) Operator(id int64) (*operator.Model, error) {
	return s.repo.Operator(id)
}

func (s *OperetorService) Operators() ([]*operator.Model, error) {
	return s.repo.Operators()
}
