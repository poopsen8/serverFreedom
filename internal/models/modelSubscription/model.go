package modelsubscription

import "time"

type Subscription struct {
	User_id    int64     `json:"user_id"`
	Plan       int       `json:"plan"`
	CreateAt   time.Time `json:"create_at"`
	Expires_at string    `json:"expires_at"`
	Key        string    `json:"key"`
}
