package subscription

import (
	"time"
	"userServer/internal/model/plan"
)

type Model struct {
	User_id    int64     `json:"user_id"`
	Plan_id    int       `json:"plan_id"`
	CreateAt   time.Time `json:"create_at"`
	Expires_at time.Time `json:"expires_at"`
	Key        string    `json:"key"`
}

type FullModel struct {
	User_id    int64       `json:"user_id"`
	Plan       *plan.Model `json:"plan_id"`
	CreateAt   time.Time   `json:"create_at"`
	Expires_at time.Time   `json:"expires_at"`
	Key        string      `json:"key"`
}
