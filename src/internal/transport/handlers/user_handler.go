package handlers

import (
	"encoding/json"
	"errors"
	"essay/src/internal/config"
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
	mux.HandleFunc("/users/login", h.HandleLogin)
	mux.HandleFunc("/users/logout", h.HandleLogout)
	mux.HandleFunc("/users/", h.HandleUsers)
	mux.HandleFunc("/users", h.HandleCreateUser)
	mux.HandleFunc("/users/count", h.GetUsersCount)
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("POST ", r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("Error decoding login request: %v\n", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Authenticate(credentials.Mail, credentials.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["is_moderator"] = user.IsModerator
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged in successfully")
}

func (h *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)

	session, _ := config.SessionStore.Get(r, "session")
	delete(session.Values, "user_id")
	delete(session.Values, "is_moderator")
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
		log.Print("Error creating user: ", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}

func (h *UserHandler) GetUsersCount(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	count, err := h.UserService.GetUsersCount()
	if err != nil {
		log.Print("Error getting users count: ", err)
		http.Error(w, "0", http.StatusInternalServerError)
		return
	}
	log.Print("Get users count: ", count)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}
