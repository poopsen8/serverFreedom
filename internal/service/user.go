package service

import (
	"userServer/internal/models"
)

type repository interface {
	CreateUser(u models.User) error
	GetUser(id int64) (models.User, error)
}

type UserService struct {
	repo repository
}

func NewUserService(repo repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(u models.User) error {
	return s.repo.CreateUser(u)
}

func (s *UserService) GetUser(id int64) (models.User, error) {
	return s.repo.GetUser(id)
}
