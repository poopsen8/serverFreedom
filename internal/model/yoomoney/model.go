package yoomoney

type Notification struct {
	NotificationType string `json:"notification_type"`
	OperationID      string `json:"operation_id"`
	Amount           string `json:"amount"` // TODO: save сумма
	WithdrawAmount   string `json:"withdraw_amount"`
	Currency         string `json:"currency"`
	DateTime         string `json:"datetime"`
	Sender           string `json:"sender"`
	Codepro          bool   `json:"codepro"`
	Label            string `json:"label"`      // TODO: label индификатор
	IsLimited        bool   `json:"unaccepted"` //TODO много денег пиздец будет
}
