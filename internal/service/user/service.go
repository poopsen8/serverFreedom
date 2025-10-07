package user

import (
	"userServer/internal/model/user"
	"userServer/internal/service/operetor"
)

type repository interface {
	Create(u user.Model) error
	User(id int64) (*user.FullModel, error)
	Update(u user.Model) error //TODO
}

type UserService struct {
	repo repository
	oper operetor.OperetorService
}

func NewUserService(repo repository, oper operetor.OperetorService) *UserService {
	return &UserService{repo: repo, oper: oper}
}

func (s *UserService) Create(u user.Model) error {
	return s.repo.Create(u)
}

func (s *UserService) User(id int64) (*user.FullModel, error) {
	user, err := s.repo.User(id)
	if err != nil {
		return nil, err
	}

	if user.MobileOperator.ID == 0 {
		return user, nil
	}
	operator, err := s.oper.Operator(user.MobileOperator.ID)
	if err != nil {
		return nil, err
	}

	user.MobileOperator = *operator
	return user, nil
}

func (s *UserService) Update(u user.Model) error {
	return s.repo.Update(u)
}
