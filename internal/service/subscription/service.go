package subscription

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
	"userServer/internal/model/subscription"

	"userServer/internal/service/operetor"
	"userServer/internal/service/plan"
	"userServer/internal/service/user"
)

type repository interface {
	Subscription(id int64) (*subscription.FullModel, error)
	Subscriptions() ([]*subscription.Model, error)
	AddSubscription(subscription.FullModel) error
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

func (s *SubscriptionService) deleteSubscription(id int64, key string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	s.js.RemoveKey(key) //TODO
	return nil

}

func (s *SubscriptionService) BackgroundCheck() error {
	subscriptions, _ := s.repo.GetSubscriptionsForCheck() //TODO

	now := time.Now()
	expiredCount := 0

	for _, sub := range subscriptions {
		if sub.Expires_at.Before(now) {
			err := s.deleteSubscription(sub.User_id, sub.Key)
			if err != nil {
				continue
			}

			expiredCount++
			log.Printf("Подписка пользователя %d истекла", sub.User_id) //TODO отпровлять что все пиздец конец подписьки

		}
	}
	return nil
}

func (s *SubscriptionService) Subscriptions() ([]*subscription.Model, error) {
	return s.repo.Subscriptions()
}

func (s *SubscriptionService) newKey() string {
	keySize := 4 //TODO добавить в конфиг

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

func (s *SubscriptionService) determineTermination(CreatedAt time.Time, Duration int64) time.Time {
	Expires_at := CreatedAt
	duration := time.Duration(Duration) * time.Minute
	Expires_at = CreatedAt.Add(duration)
	return Expires_at
}

func (s *SubscriptionService) AddSubscription(u *subscription.FullModel) (*subscription.FullModel, error) {
	u.Key = s.newKey()

	pl, err := s.pl.Plan(int64(u.Plan.ID))
	if err != nil {
		return nil, errors.New(err.Error() + "plan")
	}

	usr, err := s.us.User(u.User_id)
	if err != nil {
		return nil, errors.New(err.Error() + "user")
	}

	if usr.MobileOperator.ID == 0 {
		return nil, errors.New("operator not found")
	}

	usr.TotalSum = usr.TotalSum + int(pl.Price)
	if err := s.us.Update(*usr); err != nil {
		return nil, err
	}

	if err := s.js.AddKey(u.Key); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return nil, errors.New("error adding subscriber key")
	}

	if u.Plan.ID == 15 && !usr.IsTrial { //TODO
		return nil, errors.New("user permission denied")
	}

	if sub, err := s.Subscription(usr.ID); err == nil {
		s.deleteSubscription(sub.User_id, sub.Key)
	}

	u.Expires_at = s.determineTermination(u.CreateAt, pl.Duration)
	if err := s.repo.AddSubscription(*u); err != nil {
		return u, err
	}

	if pl.ID == 15 { //TODO
		usr.IsTrial = false
		return u, s.us.Update(*usr)
	}

	return u, nil
}
