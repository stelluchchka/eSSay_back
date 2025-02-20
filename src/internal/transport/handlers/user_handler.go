package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"essay/src/internal/config"
	"essay/src/internal/models"
	"essay/src/internal/services"
	"fmt"
	"log"
	"net/http"
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
	mux.HandleFunc("/users/nickname", h.GetNickname)
	mux.HandleFunc("/users/login", h.HandleLogin)
	mux.HandleFunc("/users/logout", h.HandleLogout)
	mux.HandleFunc("/users/info", h.HandleUserInfo)
	mux.HandleFunc("/users", h.HandleUser)
	mux.HandleFunc("/users/count", h.GetUsersCount)
}

func (h *UserHandler) GetNickname(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	mail := r.URL.Query().Get("mail")
	if mail == "" {
		log.Print("Mail parameter is missing")
		http.Error(w, "Mail parameter is missing", http.StatusBadRequest)
		return
	}

	nickname, err := h.UserService.GetNickname(mail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Print("User is not registered")
			http.Error(w, "User is not registered", http.StatusNotFound)
			return
		}
		log.Print("Error getting nickname: ", err)
		http.Error(w, "Error getting nickname", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"nickname": nickname,
	})
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
	log.Print("session: ", session)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged in successfully")
}

func (h *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)

	session, _ := config.SessionStore.Get(r, "session")
	delete(session.Values, "user_id")
	delete(session.Values, "is_moderator")
	session.Options.MaxAge = -1
	session.Save(r, w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

func (h *UserHandler) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)

	session, _ := config.SessionStore.Get(r, "session")
	id, ok := session.Values["user_id"].(uint64)
	log.Print("session: ", session)
	if !ok {
		log.Printf("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Fetching user with ID: %d\n", id)
	user, err := h.UserService.GetUserInfoByID(id)
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

func (h *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.CreateUser(w, r)
	} else if r.Method == http.MethodPut {
		h.UpdateUser(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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
			log.Print("Email already in use")
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}
		log.Print("Error creating user: ", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	created_user, err := h.UserService.Authenticate(user.Mail, user.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	session.Values["user_id"] = created_user.ID
	session.Values["is_moderator"] = created_user.IsModerator
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Print("session: ", session)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT ", r.URL.Path)

	var userData models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userData)
	if err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	id, ok := session.Values["user_id"].(uint64)
	if !ok {
		log.Printf("Unauthorized")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.UserService.UpdateUser(userData.Mail, userData.Nickname, id)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmail) {
			log.Print("Email already in use")
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}
		log.Print("Error changing info: ", err)
		http.Error(w, "Error changing info", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Info changed successfully")
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
		http.Error(w, "Error getting users count", http.StatusInternalServerError)
		return
	}
	log.Print("Users count: ", count)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"users_count": count,
	})
}
