package postgres

import (
	"userServer/internal/models"
	"userServer/internal/repository"
)

type MemoryRepo struct {
	users map[int64]models.User
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		users: make(map[int64]models.User),
	}
}

func (r *MemoryRepo) CreateUser(user models.User) error {

	if _, exists := r.users[user.ID]; exists {
		return repository.Err
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryRepo) GetUser(id int64) (models.User, error) {

	user := r.users[id]
	return user, nil
}
