package subscription

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/subscription"

	"github.com/gorilla/mux"
)

//  TODO весь нахуй файл нуно снасить нахуй и переделовать нормально

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

func (s *SubscriptionHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(s.rCfg.SR.AddSubscription, s.AddSubscription).Methods("POST") // TODO принимаент user_id ,plan_id, create_at, expires_at
	router.HandleFunc(s.rCfg.SR.Get+"{id}", s.Subscription).Methods("GET")          // TODO принимает user_id
	router.HandleFunc(s.rCfg.SR.UpdateKey+"{id}", s.UpdateKey).Methods("PUT")       // TODO принимает user_id

}

func NewSubscriptionHandler(s service, rCfg yaml.RouteConfig) *SubscriptionHandler {
	return &SubscriptionHandler{serv: s, rCfg: rCfg}
}

func (h *SubscriptionHandler) AddSubscription(w http.ResponseWriter, r *http.Request) {
	var sub subscription.Model
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, httperr.ErrInvalidJSON.Err.Error(), httperr.ErrInvalidJSON.StatusRequest)
		return
	}

	if errid := httperr.ValidateUserID(sub.User_id); errid.StatusRequest != 0 {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	err := h.serv.AddSubscription(&sub)
	if err != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest) //TODO
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": sub.User_id,
	})
}

func (h *SubscriptionHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, h.rCfg.SR.UpdateKey) //TODO
	id, err := strconv.ParseInt(path, 10, 64)

	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 || err != nil {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	key, exists := h.serv.UpdateKey(id)
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest) //TODO
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

func (h *SubscriptionHandler) Subscription(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, h.rCfg.SR.Get) //TODO
	id, err := strconv.ParseInt(path, 10, 64)
	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 || err != nil {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	sub, exists := h.serv.Subscription(id)
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest) //TODO
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	plan, exists := h.serv.GetAll()
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest) //TODO
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
