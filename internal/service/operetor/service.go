package operetor

import (
	"userServer/internal/model/operator"
)

type repository interface {
	Operator(id int64) (*operator.Model, error)
	GetAll() ([]*operator.Model, error)
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

func (s *OperetorService) GetAll() ([]*operator.Model, error) {
	return s.repo.GetAll()
}
