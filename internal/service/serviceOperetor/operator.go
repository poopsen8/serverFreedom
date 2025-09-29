package serviceOperetor

import (
	"userServer/internal/models/modelOperator"
)

type repository interface {
	Get(id int64) (*modelOperator.Operator, error)
	GetAll() ([]*modelOperator.Operator, error)
}

type OperetorService struct {
	repo repository
}

func NewOperatorService(repo repository) *OperetorService {
	return &OperetorService{repo: repo}
}

func (s *OperetorService) Get(id int64) (*modelOperator.Operator, error) {
	return s.repo.Get(id)
}

func (s *OperetorService) GetAll() ([]*modelOperator.Operator, error) {
	return s.repo.GetAll()
}
