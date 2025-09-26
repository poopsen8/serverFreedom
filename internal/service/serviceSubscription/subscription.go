package serviceSubscription

import (
	"crypto/rand"
	"encoding/hex"
	"userServer/internal/models/modelSubscription"
	"userServer/internal/service/servicePlan"
)

type repository interface {
	Get(id int64) (*modelSubscription.FullSubscription, error)
	GetAll() ([]*modelSubscription.Subscription, error)
	AddSubscription(modelSubscription.Subscription) error
	UpdateKey(id int64, key string) error
}

type SubscriptionService struct {
	repo repository
	pl   servicePlan.Repository
}

func NewSubscriptionService(repo repository, pl servicePlan.Repository) *SubscriptionService {
	return &SubscriptionService{repo: repo, pl: pl}
}

func (s *SubscriptionService) Get(id int64) (*modelSubscription.FullSubscription, error) {
	sub, _ := s.repo.Get(id) // TODO
	var err error
	sub.Plan, err = s.pl.Get(sub.Plan.ID)
	if err != nil {
		return nil, nil //TODO
	}
	return sub, nil
}

func (s *SubscriptionService) GetAll() ([]*modelSubscription.Subscription, error) {
	return s.repo.GetAll()
}

func (s *SubscriptionService) newKey() string {
	keySize := 16 //TODO добавить в конфиг

	keyBytes := make([]byte, keySize)
	_, err := rand.Read(keyBytes)
	if err != nil {
		panic(err)
	}
	keyHex := hex.EncodeToString(keyBytes) //TODO добавить проверку на совподения

	return keyHex
}

func (s *SubscriptionService) UpdateKey(id int64) (string, error) {
	key := s.newKey()

	err := s.repo.UpdateKey(id, key)
	if err != nil {
		return "", nil
	}
	return key, nil
}

func (s *SubscriptionService) AddSubscription(u *modelSubscription.Subscription) error {
	u.Key = s.newKey() //TODO нужно сделать проверку на то сущестует ли вообще такой планы
	return s.repo.AddSubscription(*u)
}
