package subscription

import (
	"database/sql"
	"userServer/internal/model/plan"
	"userServer/internal/model/subscription"

	_ "github.com/lib/pq"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) AddSubscription(u subscription.FullModel) error {
	query := `INSERT INTO subscriptions (user_id, plan_id, create_at, expires_at, key ) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	return r.db.QueryRow(query, u.User_id, u.Plan.ID, u.CreateAt, u.Expires_at, u.Key).Scan(&u.User_id)
}

func (r *SubscriptionRepository) Subscription(id int64) (*subscription.FullModel, error) {
	query := `SELECT user_id, plan_id, create_at, expires_at, key FROM subscriptions WHERE user_id = $1`
	u := &subscription.FullModel{}
	u.Plan = &plan.Model{}

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

func (r *SubscriptionRepository) GetSubscriptionsForCheck() ([]*subscription.Model, error) {
	query := `SELECT user_id, expires_at, key FROM subscriptions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*subscription.Model
	for rows.Next() {
		sub := &subscription.Model{}
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

func (r *SubscriptionRepository) Subscriptions() ([]*subscription.Model, error) {
	query := `SELECT user_id, plan_id, create_at, expires_at, key FROM subscriptions`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*subscription.Model

	for rows.Next() {
		var sub subscription.Model

		err := rows.Scan(
			&sub.User_id,
			&sub.Plan_id,
			&sub.CreateAt,
			&sub.Expires_at,
			&sub.Key,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, &sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}
