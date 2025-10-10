package operator

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/operator"
)

type service interface {
	Operator(id int64) (*operator.Model, error)
	GetAll() ([]*operator.Model, error)
}

type OperatorHandler struct {
	serv service
	rCfg yaml.RouteConfig
}

func NewOperatorHandler(s service, rCfg yaml.RouteConfig) *OperatorHandler {
	return &OperatorHandler{serv: s, rCfg: rCfg}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func (o *OperatorHandler) Operator(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/operator/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid operator ID")
		return
	}

	if errid := httperr.ValidateUserID(id); errid.StatusRequest != 0 {
		writeJSONError(w, errid.StatusRequest, errid.Err.Error())
		return
	}

	op, err := o.serv.Operator(id)
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(op)
}

func (o *OperatorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ops, err := o.serv.GetAll()
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ops)
}
