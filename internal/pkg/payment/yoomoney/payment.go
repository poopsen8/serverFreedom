package yoomoney

import (
	"fmt"
	"net/url"
	config "userServer/internal/model/config/YAML"
)

const (
	quickpayFormValue   = "button"
	quickpayPaymentType = "AC" // PC — оплата из кошелька ЮMoney; AC — с банковской карты.
	quickpayCommission  = 0.03
)

type Payment struct {
	*url.URL

	account string
}

func NewPayment(cfg *config.Yoomoney) *Payment {
	return &Payment{
		URL: &url.URL{
			Scheme: "https",
			Host:   cfg.BaseURL,
			Path:   cfg.BasePath,
		},
		account: cfg.Receiver.Account,
	}
}

func (p *Payment) Build(label string, sum float64) (string, error) {
	if len(label) <= 0 {
		return "", fmt.Errorf("empty label")
	}
	if sum <= 0 {
		return "", fmt.Errorf("empty sum")
	}

	v := url.Values{}
	v.Set("receiver", p.account)
	v.Set("quickpay-form", quickpayFormValue)
	v.Set("paymentType", quickpayPaymentType)

	amount_due := fmt.Sprintf("%.2f", calculateCommission(sum))

	v.Set("sum", amount_due)
	v.Set("label", label)

	return p.String() + "?" + v.Encode(), nil
}

func calculateCommission(sum float64) float64 {
	return float64(sum) * (1 + quickpayCommission)
}
