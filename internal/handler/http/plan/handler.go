package plan

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/plan"
)

type service interface {
	Plan(id int64) (*plan.Model, error)
	GetAll() ([]*plan.Model, error)
}

type PlanHandler struct {
	serv service
	rCfg yaml.RouteConfig
}

func NewPlanHandler(s service, rCfg yaml.RouteConfig) *PlanHandler {
	return &PlanHandler{serv: s, rCfg: rCfg}
}

// универсальная функция для JSON-ошибок
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func (p *PlanHandler) Plan(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/plan/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	pl, err := p.serv.Plan(id)
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pl)
}

func (p *PlanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	plans, err := p.serv.GetAll()
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}
