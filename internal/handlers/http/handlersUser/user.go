package handlersUser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"userServer/internal/models/modelUser"
)

//  TODO весь нахуй файл нуно снасить нахуй и переделовать нормально

type service interface {
	Create(u modelUser.User) error
	Get(id int64) (*modelUser.FullUser, error)
	Update(user modelUser.User) error
}

type UserHandler struct {
	serv service
}

func NewUserHandler(s service) *UserHandler {
	return &UserHandler{serv: s}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

	var user modelUser.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	fmt.Println(user)
	err := h.serv.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": user.ID,
	})
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/get-user/") //TODO
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, exists := h.serv.Get(id)
	if exists != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {

	var user modelUser.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if err := h.serv.Update(user); err != nil {
		http.Error(w, "ID is f", http.StatusBadRequest) //TODO
	}

	w.WriteHeader(http.StatusCreated)

}
