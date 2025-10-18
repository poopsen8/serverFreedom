package plan

type Model struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Duration  int64   `json:"duration"`
	Price     float64 `json:"price"`
	Discount  float64 `json:"discount"`
	IsPrivate bool    `json:"is_private"`
}
