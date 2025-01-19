package handlers

import (
	"encoding/json"
	"errors"
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
	mux.HandleFunc("/users/", h.HandleUsers)
	mux.HandleFunc("/users", h.HandleCreateUser)
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Println("User ID not provided in the request URL")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Invalid user ID: %v\n", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching user with ID: %d\n", id)
	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		log.Printf("Error fetching user with ID %d: %v\n", id, err)
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		log.Printf("User with ID %d not found\n", id)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("Successfully fetched user with ID %d\n", id)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Failed to encode response for user with ID %d: %v\n", id, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("POST ", r.URL.Path)
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.UserService.CreateUser(&user)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmail) {
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}
