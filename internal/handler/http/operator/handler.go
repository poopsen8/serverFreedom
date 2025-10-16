package operator

import (
	"encoding/json"
	"net/http"
	"strconv"

	httperr "userServer/internal/handler/http"
	yaml "userServer/internal/model/config/YAML"
	"userServer/internal/model/operator"
)

type service interface {
	Operator(id int64) (*operator.Model, error)
	Operators() ([]*operator.Model, error)
}

type OperatorHandler struct {
	serv service
	rCfg yaml.Config
}

func NewOperatorHandler(s service, rCfg yaml.Config) *OperatorHandler {
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
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing 'id' query parameter")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid Operator ID")
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

func (o *OperatorHandler) Operators(w http.ResponseWriter, r *http.Request) {
	ops, err := o.serv.Operators()
	if err != nil {
		writeJSONError(w, httperr.ErrIDNotFound.StatusRequest, httperr.ErrIDNotFound.Err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ops)
}
