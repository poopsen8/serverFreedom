package subscription

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
	pl "userServer/internal/model/plan"
	"userServer/internal/model/subscription"
	"userServer/internal/model/yoomoney"

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
	AddPayment(id int64, label string, price int, date_time time.Time, expires_at time.Time) error
	CheckPayment(n *yoomoney.Notification) (int64, error)
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

func (s *SubscriptionService) AddPayment(id int64, label string, price int) error {

	date_time := time.Now()
	_, expires_at := s.determineTermination(15)
	err := s.repo.AddPayment(id, label, price, date_time, expires_at)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) CheckPayment(n *yoomoney.Notification) (int64, error) {
	return s.repo.CheckPayment(n)
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

func (s *SubscriptionService) determineTermination(Duration int64) (time.Time, time.Time) {
	CreatedAt := time.Now()
	duration := time.Duration(Duration) * time.Minute
	Expires_at := CreatedAt.Add(duration)
	return CreatedAt, Expires_at
}

func (s *SubscriptionService) AddSubscription(user_id int64, plan_id int) (*subscription.FullModel, error) {
	var sub subscription.FullModel
	sub.User_id = user_id
	sub.Plan = &pl.Model{ID: int64(plan_id)}
	sub.Key = s.newKey()

	pl, err := s.pl.Plan(int64(sub.Plan.ID))
	if err != nil {
		return nil, errors.New(err.Error() + "plan")
	}

	usr, err := s.us.User(sub.User_id)
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

	if sub.Plan.ID == 0 && !usr.IsTrial { //TODO
		return nil, errors.New("user permission denied")
	}

	if sub, err := s.Subscription(usr.ID); err == nil {
		s.deleteSubscription(sub.User_id, sub.Key)
	}

	sub.CreateAt, sub.Expires_at = s.determineTermination(pl.Duration)
	if err := s.repo.AddSubscription(sub); err != nil {
		return &sub, err
	}

	if err := s.js.AddKey(sub.Key); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return nil, errors.New("error adding subscriber key")
	}
	if pl.ID == 0 { //TODO
		usr.IsTrial = false
		return &sub, s.us.Update(*usr)
	}

	return &sub, nil
}
