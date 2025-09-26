package handlerSubscription

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"userServer/internal/models/modelSubscription"
)

//  TODO весь нахуй файл нуно снасить нахуй и переделовать нормально

type service interface {
	Get(id int64) (*modelSubscription.FullSubscription, error)
	GetAll() ([]*modelSubscription.Subscription, error)
	AddSubscription(*modelSubscription.Subscription) error
	UpdateKey(id int64) (string, error)
}

type SubscriptionHandler struct {
	serv service
}

func NewSubscriptionHandler(s service) *SubscriptionHandler {
	return &SubscriptionHandler{serv: s}
}

func (h *SubscriptionHandler) AddSubscription(w http.ResponseWriter, r *http.Request) {
	var sub modelSubscription.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if sub.User_id == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	err := h.serv.AddSubscription(&sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": sub.User_id,
	})
}

func (h *SubscriptionHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/update-key-subscription/") //TODO
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	key, exists := h.serv.UpdateKey(id)
	if exists != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	type sub struct {
		User_id int64  `json:"user_id"`
		Key     string `json:"key"`
	}
	var s sub
	s.Key = key
	s.User_id = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/subscription/") //TODO
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sub, exists := h.serv.Get(id)
	if exists != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	plan, exists := h.serv.GetAll()
	if exists != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
