package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"userServer/internal/models"
)

type service interface {
	CreateUser(u models.User) error
	GetUser(id int64) (models.User, error)
}

type UserHandler struct {
	serv service
}

func NewUserHandler(s service) *UserHandler {
	return &UserHandler{serv: s}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверяем что ID передан
	if user.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Создаем пользователя с указанным ID
	err := h.serv.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, exists := h.serv.GetUser(id)
	if exists != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
