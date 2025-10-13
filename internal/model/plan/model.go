package plan

import "time"

type Model struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
	Price     float64       `json:"price"`    // Исправлено на float64
	Discount  float64       `json:"discount"` // Исправлено на float64
	IsPrivate bool          `json:"is_private"`
}
