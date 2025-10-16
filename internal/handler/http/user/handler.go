package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/user"
)

type service interface {
	Create(u user.Model) error
	User(id int64) (*user.FullModel, error)
	Update(user user.FullModel) error
	Users() ([]*user.Model, error)
}

type UserHandler struct {
	serv service
	rCfg yaml.Config
}

func NewUserHandler(s service, rCfg yaml.Config) *UserHandler {
	return &UserHandler{serv: s, rCfg: rCfg}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var usr user.Model
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		writeJSONError(w, httperr.ErrInvalidJSON.StatusRequest, httperr.ErrInvalidJSON.Err.Error())
		return
	}

	if errid := httperr.ValidateUserID(usr.ID); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	if err := u.serv.Create(usr); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "users_pkey") {
			writeJSONError(w, http.StatusOK, fmt.Sprintf("%d user already exists", usr.ID))
			return
		}

		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": usr.ID})
}

func (u *UserHandler) User(w http.ResponseWriter, r *http.Request) {
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

	usr, err := u.serv.User(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			writeJSONError(w, http.StatusOK, fmt.Sprintf("%d user not found", id))
			return
		}

		writeJSONError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

func (u *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	plans, err := u.serv.Users()
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}

func (u *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var usr user.FullModel
	if err := json.NewDecoder(r.Body).Decode(&usr); err != nil {
		writeJSONError(w, httperr.ErrInvalidJSON.StatusRequest, httperr.ErrInvalidJSON.Err.Error())
		return
	}

	if errid := httperr.ValidateUserID(usr.ID); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	if err := u.serv.Update(usr); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") && strings.Contains(err.Error(), "operator") {
			writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("%d operator not found", usr.MobileOperator.ID))
			return
		}

		writeJSONError(w, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"updated": true,
		"id":      usr.ID,
	})
}
