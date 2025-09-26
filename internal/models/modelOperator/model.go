package modelOperator

type Operator struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Is_active bool   `json:"is_active"`
}
