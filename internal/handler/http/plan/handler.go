package plan

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/plan"

	"github.com/gorilla/mux"
)

//  TODO весь нахуй файл нуно снасить нахуй и переделовать нормально

type service interface {
	Plan(id int64) (*plan.Model, error)
	GetAll() ([]*plan.Model, error)
}

type PlanHandler struct {
	serv service
	rCfg yaml.RouteConfig
}

func (p *PlanHandler) RegisterRoutes(router *mux.Router) {

	router.HandleFunc(p.rCfg.PR.Get+"{id}", p.Plan).Methods("GET") // TODO принимает id
	router.HandleFunc(p.rCfg.PR.GetAll, p.GetAll).Methods("GET")   // TODO
}

func NewPlanHandler(s service, rCfg yaml.RouteConfig) *PlanHandler {
	return &PlanHandler{serv: s, rCfg: rCfg}
}

func (p *PlanHandler) Plan(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, p.rCfg.PR.Get) //TODO
	id, err := strconv.ParseInt(path, 10, 64)
	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 || err != nil {
		http.Error(w, errid.Err.Error(), errid.StatusRequest)
		return
	}

	plan, exists := p.serv.Plan(id)
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

func (p *PlanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	plan, exists := p.serv.GetAll()
	if exists != nil {
		http.Error(w, httperr.ErrIDNotFound.Err.Error(), httperr.ErrIDNotFound.StatusRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
