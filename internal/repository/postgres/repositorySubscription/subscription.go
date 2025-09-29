package repositorySubscription

import (
	"database/sql"
	"userServer/internal/models/modelPlan"
	"userServer/internal/models/modelSubscription"

	_ "github.com/lib/pq"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository() *SubscriptionRepository {
	connStr := "host=localhost port=5432 user=postgres password=1234  dbname=postgres sslmode=disable" //TODO ПИЗДЕЦ
	db, _ := sql.Open("postgres", connStr)
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) AddSubscription(u modelSubscription.Subscription) error {
	query := `INSERT INTO subscriptions (user_id, plan_id, create_at, expires_at, key ) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	return r.db.QueryRow(query, u.User_id, u.Plan_id, u.CreateAt, u.Expires_at, u.Key).Scan(&u.User_id)
}

func (r *SubscriptionRepository) Get(id int64) (*modelSubscription.FullSubscription, error) {
	query := `SELECT user_id, plan_id, create_at, expires_at, key FROM subscriptions WHERE user_id = $1`
	u := &modelSubscription.FullSubscription{}
	u.Plan = &modelPlan.Plan{}

	err := r.db.QueryRow(query, id).Scan(&u.User_id, &u.Plan.ID, &u.CreateAt, &u.Expires_at, &u.Key)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *SubscriptionRepository) Delete(id int64) error {
	query := `DELETE FROM subscriptions WHERE user_id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SubscriptionRepository) GetSubscriptionsForCheck() ([]*modelSubscription.Subscription, error) {
	query := `SELECT user_id, expires_at, key FROM subscriptions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*modelSubscription.Subscription
	for rows.Next() {
		sub := &modelSubscription.Subscription{}
		err := rows.Scan(&sub.User_id, &sub.Expires_at, &sub.Key)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) UpdateKey(id int64, key string) error {
	query := `UPDATE subscriptions SET key = $1 WHERE user_id = $2`
	_, err := r.db.Exec(query, key, id)
	return err
}

func (r *SubscriptionRepository) GetAll() ([]*modelSubscription.Subscription, error) {

	return nil, nil
}
