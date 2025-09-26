package modelUser

import "time"

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	CreateAt       time.Time `json:"create_at"`
	MobileOperator string    `json:"mobile_operator"`
	TotalSum       int       `json:"total_sum"`
	IsTrial        bool      `json:"is_trial"`
}
