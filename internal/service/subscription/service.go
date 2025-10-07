package subscription

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"
	"userServer/internal/model/subscription"
	"userServer/internal/service/operetor"
	"userServer/internal/service/plan"
	"userServer/internal/service/user"
)

type repository interface {
	Subscription(id int64) (*subscription.FullModel, error)
	GetAll() ([]*subscription.Model, error)
	AddSubscription(subscription.Model) error
	UpdateKey(id int64, key string) error
	Delete(id int64) error
	GetSubscriptionsForCheck() ([]*subscription.Model, error)
}

type repositoryFile interface {
	AddKey(key string) error
	RemoveKey(key string) error
}

type SubscriptionService struct {
	repo repository
	pl   plan.PlanService
	js   repositoryFile
	us   user.UserService
	os   operetor.OperetorService
}

func NewSubscriptionService(repo repository, js repositoryFile, pl plan.PlanService, us user.UserService, os operetor.OperetorService) *SubscriptionService {
	return &SubscriptionService{repo: repo, pl: pl, js: js, us: us, os: os}
}

func (s *SubscriptionService) Subscription(id int64) (*subscription.FullModel, error) {
	sub, errS := s.repo.Subscription(id) // TODO
	if errS != nil {
		return nil, errS
	}

	var err error
	sub.Plan, err = s.pl.Plan(sub.Plan.ID)
	if err != nil {
		return nil, err //TODO
	}
	return sub, nil
}

func (s *SubscriptionService) BackgroundCheck() error {
	subscriptions, _ := s.repo.GetSubscriptionsForCheck() //TODO

	now := time.Now()
	expiredCount := 0

	for _, sub := range subscriptions {
		if sub.Expires_at.Before(now) {
			err := s.repo.Delete(sub.User_id)
			if err != nil {
				continue
			}

			expiredCount++
			log.Printf("Подписка пользователя %d истекла", sub.User_id) //TODO отпровлять что все пиздец конец подписьки
			s.js.RemoveKey(sub.Key)                                     // TODO

		}
	}
	return nil
}

func (s *SubscriptionService) GetAll() ([]*subscription.Model, error) {
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
	sub, _ := s.repo.Subscription(id)
	s.js.RemoveKey(sub.Key)
	s.js.AddKey(key)

	err := s.repo.UpdateKey(id, key)
	if err != nil {
		return "", nil
	}
	return key, nil
}

func (s *SubscriptionService) AddSubscription(u *subscription.Model) error {
	u.Key = s.newKey() //TODO нужно сделать проверку на то сущестует ли вообще такой планы

	if _, err := s.pl.Plan(int64(u.Plan_id)); err != nil {
		return err
	}

	if _, err := s.us.User(u.User_id); err != nil {
		return err
	}

	s.js.AddKey(u.Key)
	return s.repo.AddSubscription(*u)
}
