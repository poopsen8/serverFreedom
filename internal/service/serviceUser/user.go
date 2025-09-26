package serviceUser

import (
	"userServer/internal/models/modelUser"
)

type repository interface {
	Create(u modelUser.User) error
	Get(id int64) (*modelUser.User, error)
	Update(u modelUser.User) error //TODO
}

type UserService struct {
	repo repository
}

func NewUserService(repo repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(u modelUser.User) error {

	return s.repo.Create(u)
}

func (s *UserService) Get(id int64) (*modelUser.User, error) {

	return s.repo.Get(id)
}

func (s *UserService) Update(u modelUser.User) error {
	return s.repo.Update(u)
}
