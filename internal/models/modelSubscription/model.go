package modelSubscription

import (
	"time"
	"userServer/internal/models/modelPlan"
)

type Subscription struct {
	User_id    int64     `json:"user_id"`
	Plan_id    int       `json:"plan_id"`
	CreateAt   time.Time `json:"create_at"`
	Expires_at time.Time `json:"expires_at"`
	Key        string    `json:"key"`
}

type FullSubscription struct {
	User_id    int64
	Plan       *modelPlan.Plan
	CreateAt   time.Time
	Expires_at time.Time
	Key        string
}
