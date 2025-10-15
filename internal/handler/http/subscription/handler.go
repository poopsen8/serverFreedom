package subscription

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/subscription"
)

type service interface {
	Subscription(id int64) (*subscription.FullModel, error)
	Subscriptions() ([]*subscription.Model, error)
	AddSubscription(*subscription.FullModel) (*subscription.FullModel, error)
	UpdateKey(id int64) (string, error)
}

type SubscriptionHandler struct {
	serv service
	rCfg yaml.RouteConfig
}

func NewSubscriptionHandler(s service, rCfg yaml.RouteConfig) *SubscriptionHandler {
	return &SubscriptionHandler{serv: s, rCfg: rCfg}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func (h *SubscriptionHandler) AddSubscription(w http.ResponseWriter, r *http.Request) {
	var sub subscription.FullModel
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		writeJSONError(w, httperr.ErrInvalidJSON.StatusRequest, httperr.ErrInvalidJSON.Err.Error())
		return
	}

	if errid := httperr.ValidateUserID(sub.User_id); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	if _, err := h.serv.Subscription(sub.User_id); err != nil {
		writeJSONError(w, http.StatusConflict, fmt.Sprintf("user_id: %s ", sub.User_id))
		return
	}

	u, err := h.serv.AddSubscription(&sub)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") && strings.Contains(err.Error(), "user") {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("%d user not found", sub.User_id))
			return
		}
		if strings.Contains(err.Error(), "no rows in result set") && strings.Contains(err.Error(), "plan") {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("%d plan not found", sub.Plan.ID))
			return
		}

		if strings.Contains(err.Error(), "operator not found") {
			writeJSONError(w, http.StatusBadRequest, "operator not found")
			return
		}

		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "subscriptions_pkey") {
			writeJSONError(w, http.StatusConflict, fmt.Sprintf("Subscription %d already created", sub.User_id))
			return
		}

		if strings.Contains(err.Error(), "error adding subscriber key") {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if strings.Contains(err.Error(), "user permission denied") {
			writeJSONError(w, http.StatusConflict, fmt.Sprintf(err.Error()+"user_id: %s ", sub.User_id))
			return
		}

		writeJSONError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":  u.User_id,
		"key": u.Key,
	})
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

	sub.Key = fmt.Sprintf(
		"trojan://4cgbmbBKQO@178.17.62.19:443?type=tcp&security=reality&pbk=rmn4ibN95iKGTKXyDpGvEWXYcl1SAe28rt644tBV6VE&fp=chrome&sni=st.max.ru&%s&spx=%%2F#mgl5gsxv",
		sub.Key,
	)
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

	sub.Key = fmt.Sprintf(
		"trojan://4cgbmbBKQO@178.17.62.19:443?type=tcp&security=reality&pbk=rmn4ibN95iKGTKXyDpGvEWXYcl1SAe28rt644tBV6VE&fp=chrome&sni=st.max.ru&%s&spx=%%2F#mgl5gsxv",
		sub.Key,
	)
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
