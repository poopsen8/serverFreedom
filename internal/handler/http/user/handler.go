package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/user"

	"github.com/gorilla/mux"
)

type service interface {
	Create(u user.Model) error
	User(id int64) (*user.FullModel, error)
	Update(user user.Model) error
}

type UserHandler struct {
	serv service
	rCfg yaml.RouteConfig
}

func (u *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(u.rCfg.UR.Register, u.Create).Methods("POST")    //   принимает id, username,
	router.HandleFunc(u.rCfg.UR.Update, u.Update).Methods("PUT")       //  принимает одно значения на изменения и id пользователя - MobileOperatorID, IsTrial
	router.HandleFunc((u.rCfg.UR.Get + "{id}"), u.User).Methods("GET") // TODO принимает id
}

func NewUserHandler(s service, rCfg yaml.RouteConfig) *UserHandler {
	return &UserHandler{serv: s, rCfg: rCfg}
}

func (u *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user user.Model
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, httperr.ErrInvalidJSON.Err.Error(), httperr.ErrInvalidJSON.StatusRequest)
		return
	}

	if errid := httperr.ValidateUserID(user.ID); errid.StatusRequest != 0 {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	err := u.serv.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") //TODO
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": user.ID,
	})
}

func (u *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, u.rCfg.UR.Get)
	id, err := strconv.ParseInt(path, 10, 64)

	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 || err != nil {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	user, exists := u.serv.User(id)
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json") //TODO
	json.NewEncoder(w).Encode(user)
}

func (u *UserHandler) Update(w http.ResponseWriter, r *http.Request) {

	var user user.Model
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, httperr.ErrInvalidJSON.Err.Error(), httperr.ErrInvalidJSON.StatusRequest)
		return
	}

	if errid := httperr.ValidateUserID(user.ID); errid.StatusRequest != 0 {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	if err := u.serv.Update(user); err != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest) //TODO + нужно обрабоать ошибку мол нет такого поля или еще чего
	}

	w.WriteHeader(http.StatusCreated)

}
