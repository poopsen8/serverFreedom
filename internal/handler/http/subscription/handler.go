package subscription

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/subscription"
	y "userServer/internal/model/yoomoney"
	"userServer/internal/pkg/payment/yoomoney"
)

type service interface {
	Subscription(id int64) (*subscription.FullModel, error)
	Subscriptions() ([]*subscription.Model, error)
	AddSubscription(user_id int64, plan_id int) (*subscription.FullModel, error)
	UpdateKey(id int64) (string, error)
	AddPayment(id int64, label string, price int) error
	CheckPayment(n *y.Notification) (int64, error)
}

type SubscriptionHandler struct {
	serv service
	rCfg yaml.Config
}

func NewSubscriptionHandler(s service, rCfg yaml.Config) *SubscriptionHandler {
	return &SubscriptionHandler{serv: s, rCfg: rCfg}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func (h *SubscriptionHandler) Validator(n *y.Notification) {
	user_id, err := h.serv.CheckPayment(n)
	if err != nil {
		log.Printf("Ошибка проверки платежа: %v", err)
		return
	}

	values, err := url.ParseQuery(n.Label)
	if err != nil {
		log.Printf("Ошибка парсинга Label: %v", err)
		return
	}

	planIDStr := values.Get("plan_id")
	if planIDStr == "" {
		log.Println("plan_id не найден в метке")
		return
	}

	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		log.Printf("Некорректный plan_id: %v", err)
		return
	}

	// Добавляем подписку
	h.AddSubscription(user_id, planID)
}

func (h *SubscriptionHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	var sub subscription.FullModel
	fmt.Printf("r.Body: %v\n", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		writeJSONError(w, httperr.ErrInvalidJSON.StatusRequest, httperr.ErrInvalidJSON.Err.Error())
		return
	}

	p := yoomoney.NewPayment(h.rCfg.Yoomoney)
	l := fmt.Sprintf("%suser_id=%d&plan_id=%d",
		time.Now().Format("2006-01-02 15:04:05"),
		sub.User_id,
		sub.Plan.ID)
	price := sub.Plan.Price - ((sub.Plan.Price / 100) * sub.Plan.Discount)
	url, _ := p.Build(l, int(price))

	h.serv.AddPayment(sub.User_id, l, int(price))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"url": url,
	})
}

func (h *SubscriptionHandler) AddSubscription(user_id int64, plan_id int) {
	// Вспомогательная функция для отправки ошибок
	sendError := func(statusCode int, message string) {
		errorResponse := map[string]interface{}{
			"error":   true,
			"status":  statusCode,
			"message": message,
		}
		jsonData, _ := json.Marshal(errorResponse)
		http.Post(h.rCfg.PathSend.Base_url, "application/json", bytes.NewBuffer(jsonData))
	}

	if errid := httperr.ValidateUserID(user_id); errid.StatusRequest != 0 {
		sendError(errid.StatusRequest, errid.Err.Error())
		return
	}

	sub, err := h.serv.AddSubscription(user_id, plan_id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") && strings.Contains(err.Error(), "user") {
			sendError(http.StatusBadRequest, fmt.Sprintf("%d user not found", user_id))
			return
		}
		if strings.Contains(err.Error(), "no rows in result set") && strings.Contains(err.Error(), "plan") {
			sendError(http.StatusBadRequest, fmt.Sprintf("%d plan not found", plan_id))
			return
		}

		if strings.Contains(err.Error(), "operator not found") {
			sendError(http.StatusBadRequest, "operator not found")
			return
		}

		if strings.Contains(err.Error(), "error adding subscriber key") {
			sendError(http.StatusInternalServerError, err.Error())
			return
		}
		if strings.Contains(err.Error(), "user permission denied") {
			sendError(http.StatusConflict, fmt.Sprintf(err.Error()+"user_id: %d ", user_id))
			return
		}

		fmt.Printf("err.Error(): %v\n", err.Error())
		sendError(http.StatusInternalServerError, err.Error())
		return
	}

	sub.Key = h.rCfg.Link.Left_part + sub.Key + h.rCfg.Link.Right_part
	// Отправка успешного результата
	u, err := json.Marshal(sub)
	if err != nil {
		sendError(http.StatusInternalServerError, "Failed to marshal subscription data")
		return
	}

	resp, err := http.Post(h.rCfg.PathSend.Base_url, "application/json", bytes.NewBuffer(u))
	if err != nil {
		fmt.Printf("Failed to send subscription data: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func (h *SubscriptionHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing 'id' query parameter")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	_, e := h.serv.UpdateKey(id)
	if e != nil {
		writeJSONError(w, 500, httperr.ErrIDNotFound.Err.Error())
		return
	}

	sub, er := h.serv.Subscription(id)
	if er != nil {
		writeJSONError(w, 500, httperr.ErrIDNotFound.Err.Error())
		return
	}

	sub.Key = h.rCfg.Link.Left_part + sub.Key + h.rCfg.Link.Right_part
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Subscription(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing 'id' query parameter")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	sub, err := h.serv.Subscription(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			writeJSONError(w, http.StatusOK, fmt.Sprintf("%d subscription not found", id))
			return
		}
		writeJSONError(w, 500, httperr.ErrIDNotFound.Err.Error())
		return
	}

	sub.Key = h.rCfg.Link.Left_part + sub.Key + h.rCfg.Link.Right_part
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Subscriptions(w http.ResponseWriter, r *http.Request) {
	plan, err := h.serv.Subscriptions()
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
