package serviceUser

import (
	"userServer/internal/models/modelUser"
	"userServer/internal/service/serviceOperetor"
)

type repository interface {
	Create(u modelUser.User) error
	Get(id int64) (*modelUser.FullUser, error)
	Update(u modelUser.User) error //TODO
}

type UserService struct {
	repo repository
	oper serviceOperetor.OperetorService
}

func NewUserService(repo repository, oper serviceOperetor.OperetorService) *UserService {
	return &UserService{repo: repo, oper: oper}
}

func (s *UserService) Create(u modelUser.User) error {
	return s.repo.Create(u)
}

func (s *UserService) Get(id int64) (*modelUser.FullUser, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if user.MobileOperator.ID == 0 {
		return user, nil
	}
	operator, err := s.oper.Get(user.MobileOperator.ID)
	if err != nil {
		return nil, err
	}

	user.MobileOperator = *operator
	return user, nil
}

func (s *UserService) Update(u modelUser.User) error {
	return s.repo.Update(u)
}
