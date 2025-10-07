package operator

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/operator"

	"github.com/gorilla/mux"
)

//  TODO весь нахуй файл нуно снасить нахуй и переделовать нормально

type service interface {
	Operator(id int64) (*operator.Model, error)
	GetAll() ([]*operator.Model, error)
}

type OperatorHandler struct {
	serv service
	rCfg yaml.RouteConfig
}

func (o *OperatorHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(o.rCfg.OP.Get+"{id}", o.Operator).Methods("GET") // TODO принимает id
	router.HandleFunc(o.rCfg.OP.GetAll, o.GetAll).Methods("GET")       // TODO
}

func NewOperatorHandler(s service, rCfg yaml.RouteConfig) *OperatorHandler {
	return &OperatorHandler{serv: s, rCfg: rCfg}
}

func (o *OperatorHandler) Operator(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, o.rCfg.OP.Get) //TODO
	id, err := strconv.ParseInt(path, 10, 64)
	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 || err != nil {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	plan, exists := o.serv.Operator(id)
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

func (o *OperatorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	plan, exists := o.serv.GetAll()
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
