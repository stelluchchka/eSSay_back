package rest

import (
	"encoding/json"
	"essay/src/internal/models"
	"essay/src/internal/services"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/", h.HandleUsers)     // /users/{id}
	mux.HandleFunc("/users", h.HandleCreateUser) // /users (POST)
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("HandleCreateUser %s", err)
		log.Printf("%s", decoder.Decode(&user))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.UserService.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}
