package modelUser

import (
	"time"
	"userServer/internal/models/modelOperator"
)

type User struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	CreateAt         time.Time `json:"create_at"`
	MobileOperatorID int       `json:"operator_id"`
	TotalSum         int       `json:"total_sum"`
	IsTrial          bool      `json:"is_trial"`
}

type FullUser struct {
	ID             int64
	Username       string
	CreateAt       time.Time
	MobileOperator modelOperator.Operator
	TotalSum       int
	IsTrial        bool
}
