package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"essay/src/internal/config"
	"essay/src/internal/models"
	"essay/src/internal/services"
)

type EssayHandler struct {
	EssayService *services.EssayService
}

func NewEssayHandler(essayService *services.EssayService) *EssayHandler {
	return &EssayHandler{
		EssayService: essayService,
	}
}

func (h *EssayHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/essays", h.HandleEssaysRequests)
	mux.HandleFunc("/essays/", h.HandleEssayRequests)
	mux.HandleFunc("/essays/appeal", h.GetAppealEssays)
	mux.HandleFunc("/users/me/essays", h.GetUserEssays)
	mux.HandleFunc("/users/me/essays/", h.GetUserEssayByID)
}

// HandleEssaysRequests handles various methods on essays.
func (h *EssayHandler) HandleEssaysRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetPublishedEssays(w, r)
		return
	} else if r.Method == http.MethodPost {
		h.CreateEssay(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// HandleEssayRequests handles various methods on essay.
func (h *EssayHandler) HandleEssayRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetEssayByID(w, r)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		log.Println("Essay ID not provided in the request URL")
		http.Error(w, "Essay ID is required", http.StatusBadRequest)
		return
	} else if len(parts) > 4 {
		log.Println("404 page not found")
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPut && (strings.HasSuffix(r.URL.Path, "/save") || strings.HasSuffix(r.URL.Path, "/appeal") || strings.HasSuffix(r.URL.Path, "/publish")) {
		h.ChangeEssayStatus(w, r)
		return
	} else if r.Method == http.MethodPut {
		h.UpdateEssay(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// GetPublishedEssays handles GET /essays.
func (h *EssayHandler) GetPublishedEssays(w http.ResponseWriter, r *http.Request) {
	log.Print("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	essays, err := h.EssayService.GetPublishedEssays()
	if err != nil {
		log.Printf("Error retrieving essays: %v", err)
		http.Error(w, "Failed to retrieve essays", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essays)
}

// GetAppealEssays handles GET /essays/appeal.
func (h *EssayHandler) GetAppealEssays(w http.ResponseWriter, r *http.Request) {
	log.Print("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	isModerator, ok := session.Values["is_moderator"].(bool)
	if !ok || !isModerator {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	essays, err := h.EssayService.GetAppealEssays()
	if err != nil {
		log.Printf("Error retrieving essays: %v", err)
		http.Error(w, "Failed to retrieve essays", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essays)
}

// GetEssayByID handles GET /essays/:id.
func (h *EssayHandler) GetEssayByID(w http.ResponseWriter, r *http.Request) {
	log.Print("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Invalid essay ID: %v", err)
		http.Error(w, "Invalid essay ID", http.StatusBadRequest)
		return
	}

	essay, err := h.EssayService.GetDetailedEssayByID(uint8(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Essay not found", http.StatusNotFound)
		} else {
			log.Printf("Error GetDetailedEssayByID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// session, _ := config.SessionStore.Get(r, "session")
	// userIDInterface, ok := session.Values["user_id"]
	// if ok {
	// 	userID := userIDInterface.(uint8)
	// 	if essay.AuthorID != userID && !essay.IsPublished {
	// 		http.Error(w, "Forbidden", http.StatusForbidden)
	// 		return
	// 	}
	// }
	if !essay.IsPublished {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var detailedEssay models.DetailedEssay

	switch essay.Status {
	case "draft":
		http.Error(w, "Cannot access a draft", http.StatusBadRequest)
		return
	case "saved", "checked", "appeal", "appealed":
		detailedEssay = *essay
	}

	if essay.IsPublished {
		detailedEssay.Comments = essay.Comments
		detailedEssay.Likes = essay.Likes
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detailedEssay)
}

// GetUserEssays handles GET /users/me/essays.
func (h *EssayHandler) GetUserEssays(w http.ResponseWriter, r *http.Request) {
	log.Print("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	userIDInterface, ok := session.Values["user_id"]
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID := userIDInterface.(uint8)

	essays, err := h.EssayService.GetUserEssays(userID)
	if err != nil {
		log.Printf("Error retrieving user essays: %v", err)
		http.Error(w, "Failed to retrieve user essays", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essays)
}

// GetEssayByID handles GET /users/me/essays/:id.
func (h *EssayHandler) GetUserEssayByID(w http.ResponseWriter, r *http.Request) {
	log.Print("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Printf("Invalid essay ID: %v", err)
		http.Error(w, "Invalid essay ID", http.StatusBadRequest)
		return
	}

	essay, err := h.EssayService.GetDetailedEssayByID(uint8(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Essay not found", http.StatusNotFound)
		} else {
			log.Printf("Error GetDetailedEssayByID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	userIDInterface, ok := session.Values["user_id"]
	if ok {
		userID := userIDInterface.(uint8)
		if essay.AuthorID != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	var detailedEssay models.DetailedEssay

	switch essay.Status {
	case "draft":
		http.Error(w, "Cannot access a draft", http.StatusBadRequest)
		return
	case "saved", "checked", "appeal", "appealed":
		detailedEssay = *essay
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detailedEssay)
}

// CreateEssay handles POST /essays.
func (h *EssayHandler) CreateEssay(w http.ResponseWriter, r *http.Request) {
	log.Print("POST ", r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody struct {
		EssayText string `json:"essay_text"`
		VariantId uint8  `json:"variant_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	userIDInterface, ok := session.Values["user_id"]
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID := userIDInterface.(uint8)

	var essay models.Essay
	essay.Status = "draft"
	essay.EssayText = reqBody.EssayText
	essay.VariantID = reqBody.VariantId
	essay.UserID = userID

	if err := h.EssayService.CreateEssay(&essay); err != nil {
		log.Printf("Failed to create essay: %v", err)
		http.Error(w, "Failed to create essay", http.StatusInternalServerError)
		return
	}

	log.Printf("Essay created successfully: %+v", essay)
	w.WriteHeader(http.StatusCreated)
}

// UpdateEssay handles PUT /essays/:id.
func (h *EssayHandler) UpdateEssay(w http.ResponseWriter, r *http.Request) {
	log.Print("PUT ", r.URL.Path)

	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Invalid essay ID: %v", err)
		http.Error(w, "Invalid essay ID", http.StatusBadRequest)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	userIDInterface, ok := session.Values["user_id"]
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID := userIDInterface.(uint8)

	essay, err := h.EssayService.GetEssayByID(uint8(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to find essay with id %d: %v", id, err)
			http.Error(w, "Failed to find essay", http.StatusBadRequest)
			return
		}
		log.Printf("Failed to find essay with id %d: %v", id, err)
		http.Error(w, "Failed to find essay", http.StatusInternalServerError)
		return
	}
	if essay.UserID != userID {
		log.Printf("Failed to save essay with id %d: wrong user ID", id)
		http.Error(w, "Failed to save essay: wrong user ID", http.StatusForbidden)
		return
	}

	var reqBody struct {
		EssayText string `json:"essay_text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	var newEssay models.Essay
	newEssay.ID = uint8(id)
	newEssay.EssayText = reqBody.EssayText
	newEssay.UserID = userID

	log.Printf("Updating essay with id %d, text %s", newEssay.ID, newEssay.EssayText)
	if err := h.EssayService.UpdateEssay(&newEssay); err != nil {
		if errors.Is(err, services.ErrWrongID) {
			http.Error(w, "Wrong id", http.StatusBadRequest)
			return
		}
		log.Printf("Failed to update essay: %v", err)
		http.Error(w, "Failed to update essay", http.StatusInternalServerError)
		return
	}

	log.Print("Essay updated successfully")
	w.WriteHeader(http.StatusOK)
}

// ChangeEssayStatus handles PUT /essays/:id/<action>.
func (h *EssayHandler) ChangeEssayStatus(w http.ResponseWriter, r *http.Request) {
	log.Print("PUT ", r.URL.Path)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		log.Println("Essay ID or action not provided in the request URL")
		http.Error(w, "Essay ID and action is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Invalid essay ID: %v", err)
		http.Error(w, "Invalid essay ID", http.StatusBadRequest)
		return
	}
	action := parts[3]

	session, _ := config.SessionStore.Get(r, "session")
	userIDInterface, ok := session.Values["user_id"]
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID := userIDInterface.(uint8)

	essay, err := h.EssayService.GetEssayByID(uint8(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to find essay with id %d: %v", id, err)
			http.Error(w, "Failed to find essay", http.StatusBadRequest)
			return
		}
		log.Printf("Failed to find essay with id %d: %v", id, err)
		http.Error(w, "Failed to find essay", http.StatusInternalServerError)
		return
	}
	if essay.UserID != userID {
		log.Printf("Failed to set status %s to essay with id %d: wrong user ID", action, id)
		http.Error(w, "Failed to save essay: wrong user ID", http.StatusForbidden)
		return
	}

	var status string

	switch action {
	case "save":
		if essay.Status != "draft" {
			log.Printf("Failed to save essay with id %d: status should be draft", id)
			http.Error(w, "Failed to save essay: status should be draft", http.StatusBadRequest)
			return
		}
		status = "saved"
		// TODO: положить в очередь на проверку
	case "appeal":
		if essay.Status != "checked" {
			log.Printf("Failed to file appeal for essay with id %d: status should be checked", id)
			http.Error(w, "Failed to file appeal for essay: status should be checked", http.StatusBadRequest)
			return
		}
		status = "appeal"
		// TODO: добавить текст аппеляции
	case "publish":
		log.Printf("Publishing essay: ID %d", id)
		if err := h.EssayService.PublishEssay(uint8(id), userID); err != nil {
			log.Printf("Failed to publish essay: %v", err)
			http.Error(w, "Failed to publish essay", http.StatusInternalServerError)
			return
		}
		log.Printf("Essay published successfully: ID %d", id)
		w.WriteHeader(http.StatusOK)
		return
	default:
		log.Printf("Invalid action: %s", action)
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	log.Printf("Changing essay status to '%s' for essayID %d for userID %d", status, id, userID)
	if err := h.EssayService.ChangeEssayStatus(uint8(id), userID, status); err != nil {
		log.Printf("Failed to change essay status: %v", err)
		http.Error(w, "Failed to change essay status", http.StatusInternalServerError)
		return
	}

	log.Print("Essay status changed successfully")
	w.WriteHeader(http.StatusOK)
}
