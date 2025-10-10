package subscription

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/subscription"
)

type service interface {
	Subscription(id int64) (*subscription.FullModel, error)
	GetAll() ([]*subscription.Model, error)
	AddSubscription(*subscription.Model) error
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
	var sub subscription.Model
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		writeJSONError(w, httperr.ErrInvalidJSON.StatusRequest, httperr.ErrInvalidJSON.Err.Error())
		return
	}

	if errid := httperr.ValidateUserID(sub.User_id); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	if err := h.serv.AddSubscription(&sub); err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": sub.User_id,
	})
}

func (h *SubscriptionHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/update-key-subscription/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	key, err := h.serv.UpdateKey(id)
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	type sub struct {
		User_id int64  `json:"user_id"`
		Key     string `json:"key"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub{
		User_id: id,
		Key:     key,
	})
}

func (h *SubscriptionHandler) Subscription(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/subscription/")
	id, err := strconv.ParseInt(path, 10, 64)
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
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	plan, err := h.serv.GetAll()
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
