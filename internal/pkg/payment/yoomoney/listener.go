package yoomoney

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	config "userServer/internal/model/config/YAML"
	"userServer/internal/model/yoomoney"
)

func PaymentHandler(cfg *config.Yoomoney, validator func(n *yoomoney.Notification)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		fmt.Println(cfg.Receiver.Account)
		fmt.Println(cfg.Receiver.NotifSecret)

		isValid, err := VerifySignature(r, cfg.Receiver.NotifSecret)
		if err != nil {
			http.Error(w, "Error verifying signature", http.StatusInternalServerError)
			return
		}
		if !isValid {
			http.Error(w, "Invalid signature", http.StatusForbidden)
			return
		}

		notif := yoomoney.Notification{
			NotificationType: r.PostFormValue("notification_type"),
			OperationID:      r.PostFormValue("operation_id"),
			Amount:           r.PostFormValue("amount"),
			WithdrawAmount:   r.PostFormValue("withdraw_amount"),
			Currency:         r.PostFormValue("currency"),
			DateTime:         r.PostFormValue("datetime"),
			Sender:           r.PostFormValue("sender"),
			Codepro:          r.PostFormValue("codepro") == "true",
			Label:            r.PostFormValue("label"),
			IsLimited:        r.PostFormValue("unaccepted") == "true",
		}

		// FIXME
		validator(&notif)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func VerifySignature(r *http.Request, secret string) (bool, error) {
	params := []string{
		r.PostFormValue("notification_type"),
		r.PostFormValue("operation_id"),
		r.PostFormValue("amount"),
		r.PostFormValue("currency"),
		r.PostFormValue("datetime"),
		r.PostFormValue("sender"),
		r.PostFormValue("codepro"),
		secret,
		r.PostFormValue("label"),
	}

	signString := strings.Join(params, "&")

	hash := sha1.Sum([]byte(signString))
	calculatedSignature := hex.EncodeToString(hash[:])

	receivedSignature := r.PostFormValue("sha1_hash")

	return subtle.ConstantTimeCompare([]byte(calculatedSignature), []byte(receivedSignature)) == 1, nil
}
