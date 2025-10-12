package user

import (
	"time"
	"userServer/internal/model/operator"
)

type Model struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	CreateAt         time.Time `json:"create_at"`
	MobileOperatorID int64     `json:"operator_id"`
	TotalSum         int       `json:"total_sum"`
	IsTrial          bool      `json:"is_trial"`
}

type FullModel struct {
	ID             int64
	Username       string
	CreateAt       time.Time
	MobileOperator operator.Model
	TotalSum       int
	IsTrial        bool
}
