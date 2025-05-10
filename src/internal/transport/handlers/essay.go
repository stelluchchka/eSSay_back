package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"essay/src/internal/config"
	"essay/src/internal/models"
	"essay/src/internal/services"
)

// HandleEssaysRequests handles various methods on essays.
func (h *UserHandler) HandleEssaysRequests(w http.ResponseWriter, r *http.Request) {
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
func (h *UserHandler) HandleEssayRequests(w http.ResponseWriter, r *http.Request) {
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
func (h *UserHandler) GetPublishedEssays(w http.ResponseWriter, r *http.Request) {
	log.Print("GET ", r.URL.Path)

	essays, err := h.UserService.GetPublishedEssays()
	if err != nil {
		log.Printf("Error retrieving essays: %v", err)
		http.Error(w, "Failed to retrieve essays", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essays)
}

// GetAppealEssays handles GET /essays/appeal.
func (h *UserHandler) GetAppealEssays(w http.ResponseWriter, r *http.Request) {
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

	essays, err := h.UserService.GetAppealEssays()
	if err != nil {
		log.Printf("Error retrieving essays: %v", err)
		http.Error(w, "Failed to retrieve essays", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essays)
}

// GetEssayByID handles GET /essays/:id.
func (h *UserHandler) GetEssayByID(w http.ResponseWriter, r *http.Request) {
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

	essay, err := h.UserService.GetDetailedEssayByID(uint64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Essay not found", http.StatusNotFound)
		} else {
			log.Printf("Error GetDetailedEssayByID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

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
func (h *UserHandler) GetUserEssays(w http.ResponseWriter, r *http.Request) {
	log.Print("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	userID, ok := session.Values["user_id"].(uint64)
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	essays, err := h.UserService.GetUserEssays(userID)
	if err != nil {
		log.Printf("Error retrieving user essays: %v", err)
		http.Error(w, "Failed to retrieve user essays", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essays)
}

// GetEssayByID handles GET /users/me/essays/:id.
func (h *UserHandler) GetUserEssayByID(w http.ResponseWriter, r *http.Request) {
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

	essay, err := h.UserService.GetDetailedEssayByID(uint64(id))
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
		userID := userIDInterface.(uint64)
		if essay.AuthorID != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	// var detailedEssay models.DetailedEssay

	// detailedEssay = *essay

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(essay)
}

// CreateEssay handles POST /essays.
func (h *UserHandler) CreateEssay(w http.ResponseWriter, r *http.Request) {
	log.Print("POST ", r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody struct {
		EssayText string `json:"essay_text"`
		VariantId uint64 `json:"variant_id"`
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
	userID := userIDInterface.(uint64)

	var essay models.Essay
	essay.Status = "draft"
	essay.EssayText = reqBody.EssayText
	essay.VariantID = reqBody.VariantId
	essay.UserID = userID

	essayId, err := h.UserService.CreateEssay(&essay)
	if err != nil {
		log.Printf("Failed to create essay: %v", err)
		http.Error(w, "Failed to create essay", http.StatusInternalServerError)
		return
	}

	log.Printf("Essay created successfully: %+v", essay)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"essay_id": essayId,
	})
}

// UpdateEssay handles PUT /essays/:id.
func (h *UserHandler) UpdateEssay(w http.ResponseWriter, r *http.Request) {
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
	userID := userIDInterface.(uint64)

	essay, err := h.UserService.GetEssayByID(uint64(id))
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
	newEssay.ID = uint64(id)
	newEssay.EssayText = reqBody.EssayText
	newEssay.UserID = userID

	log.Printf("Updating essay with id %d, text %s", newEssay.ID, newEssay.EssayText)
	if err := h.UserService.UpdateEssay(&newEssay); err != nil {
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
func (h *UserHandler) ChangeEssayStatus(w http.ResponseWriter, r *http.Request) {
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
	userID := userIDInterface.(uint64)

	essay, err := h.UserService.GetEssayByID(uint64(id))
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
			log.Printf("Failed to save essay with id %d: status should be draft but it is %s", id, essay.Status)
			http.Error(w, "Failed to save essay: status should be draft", http.StatusBadRequest)
			return
		}

		// check count verification
		err = h.UserService.DecreaseCheckCount(userID)
		if err != nil {
			if errors.Is(err, services.ErrNoChecksLeft) {
				log.Printf("Failed to save essay with id %d: no checks left", id)
				http.Error(w, "No checks left", http.StatusNotFound)
				return
			}
			log.Printf("Failed to decrease check count: %v", err)
			http.Error(w, "Failed to decrease check count", http.StatusInternalServerError)
			return
		}

		status = "saved"
		log.Printf("Essay ID %d enqueued for checking", id)
		log.Printf("Changing essay status to '%s' for essayID %d", status, id)
		if err := h.UserService.ChangeEssayStatus(uint64(id), status); err != nil {
			log.Printf("Failed to change essay status: %v", err)
			http.Error(w, "Failed to change essay status", http.StatusInternalServerError)
			return
		}

		vaiant, err := h.UserService.GetVariantByID(essay.VariantID)
		if err != nil {
			log.Printf("Failed to get variant in ChangeEssayStatus: %v", err)
			http.Error(w, "Failed to get variant in ChangeEssayStatus", http.StatusInternalServerError)
			return
		}
		requestBody, err := json.Marshal(models.EssayRequest{
			EssayID:        essay.ID,
			EssayText:      essay.EssayText,
			VariantText:    vaiant.VariantText,
			AuthorPosition: vaiant.AuthorPosition,
		})
		if err != nil {
			log.Printf("Failed to marshal JSON in ChangeEssayStatus: %v", err)
			http.Error(w, "Failed to marshal JSON in ChangeEssayStatus", http.StatusInternalServerError)
			return
		}

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Post(config.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Failed to send request to Python service: %v", err)
			http.Error(w, "Failed to send request to Python service", http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		log.Printf("Send to Python service")

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error from Python service: %s", resp.Status)
			http.Error(w, "Error from Python service", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		return
	case "appeal":
		if essay.Status != "checked" {
			log.Printf("Failed to file appeal for essay with id %d: status should be checked but it is %s", id, essay.Status)
			http.Error(w, "Failed to file appeal for essay: status should be checked", http.StatusBadRequest)
			return
		}
		status = "appeal"

		var reqBody struct {
			AppealText string `json:"appeal_text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Printf("Invalid request body for appeal: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := h.UserService.UnpublishEssay(uint64(id), userID); err != nil {
			log.Printf("Failed to unpublish essay: %v", err)
			http.Error(w, "Failed to unpublish essay", http.StatusInternalServerError)
			return
		}
		if err := h.UserService.SetAppealText(uint64(id), reqBody.AppealText); err != nil {
			log.Printf("Failed to set appeal text: %v", err)
			http.Error(w, "Failed to set appeal text", http.StatusInternalServerError)
			return
		}
	case "publish":
		log.Printf("Publishing essay: ID %d", id)
		if err := h.UserService.PublishEssay(uint64(id), userID); err != nil {
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

	log.Printf("Changing essay status to '%s' for essayID %d", status, id)
	if err := h.UserService.ChangeEssayStatus(uint64(id), status); err != nil {
		log.Printf("Failed to change essay status: %v", err)
		http.Error(w, "Failed to change essay status", http.StatusInternalServerError)
		return
	}

	log.Print("Essay status changed successfully")
	w.WriteHeader(http.StatusOK)
}
