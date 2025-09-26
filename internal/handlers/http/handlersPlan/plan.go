package handlersPlan

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"userServer/internal/models/modelPlan"
)

//  TODO весь нахуй файл нуно снасить нахуй и переделовать нормально

type service interface {
	Get(id int64) (*modelPlan.Plan, error)
	GetAll() ([]*modelPlan.Plan, error)
}

type PlanHandler struct {
	serv service
}

func NewPlanHandler(s service) *PlanHandler {
	return &PlanHandler{serv: s}
}

func (h *PlanHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/get-plan/") //TODO
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	plan, exists := h.serv.Get(id)
	if exists != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

func (h *PlanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	plan, exists := h.serv.GetAll()
	if exists != nil {
		http.Error(w, "plan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
