package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"essay/src/internal/config"
	"essay/src/internal/services"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ContentHandler struct {
	ContentService *services.ContentService
	EssayService   *services.EssayService
}

func NewContentHandler(contentService *services.ContentService, essayService *services.EssayService) *ContentHandler {
	return &ContentHandler{
		ContentService: contentService,
		EssayService:   essayService,
	}
}

func (h *ContentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/counts", h.GetCounts)
	mux.HandleFunc("/likes/is_liked/", h.HandleIsLiked)
	mux.HandleFunc("/likes/", h.HandleLikes)
	mux.HandleFunc("/comments/", h.HandleComments)
	mux.HandleFunc("/variants/", h.GetVariant)
}

func (h *ContentHandler) HandleIsLiked(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)

	essayId, err := strconv.Atoi(r.URL.Path[len("/likes/is_liked/"):])
	if err != nil {
		http.Error(w, "Invalid essay ID", http.StatusBadRequest)
		return
	}

	_, err = h.EssayService.GetEssayByID(uint64(essayId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to find essay with id %d: %v", essayId, err)
			http.Error(w, "Essay not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching essay with id %d: %v", essayId, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session, _ := config.SessionStore.Get(r, "session")
	userID, ok := session.Values["user_id"].(uint64)
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	isLiked, err := h.ContentService.IsLiked(userID, uint64(essayId))
	if err != nil {
		log.Printf("Error checking if liked: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Print("isLiked", isLiked)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"is_liked": isLiked})
}

func (h *ContentHandler) HandleLikes(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	id, err := strconv.Atoi(r.URL.Path[len("/likes/"):])
	if err != nil {
		http.Error(w, "Invalid essay ID", http.StatusBadRequest)
		return
	}

	_, err = h.EssayService.GetEssayByID(uint64(id))
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

	switch r.Method {
	case http.MethodGet:
		count, err := h.ContentService.GetLikesCount(uint64(id))
		if err != nil {
			log.Printf("Error fetching likes count: %s", err)
			http.Error(w, "Error fetching likes count", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"likes": count})

	case http.MethodPut:
		session, _ := config.SessionStore.Get(r, "session")
		userID, ok := session.Values["user_id"].(uint64)
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if isLiked, err := h.ContentService.IsLiked(userID, uint64(id)); err != nil {
			log.Print("Error with like: ", err)
			http.Error(w, "Error with like", http.StatusInternalServerError)
			return
		} else if !isLiked {
			if err := h.ContentService.AddLike(userID, uint64(id)); err != nil {
				if errors.Is(err, services.ErrLikeAlreadyExists) {
					log.Print("Error adding like: ", err)
					http.Error(w, "Error adding like", http.StatusBadRequest)
					return
				}
				http.Error(w, "Error adding like", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Like added successfully"))
		} else {
			if err := h.ContentService.DeleteLike(userID, uint64(id)); err != nil {
				if errors.Is(err, services.ErrLikeAlreadyExists) {
					log.Print("Error removing like: ", err)
					http.Error(w, "Error removing like", http.StatusBadRequest)
					return
				}
				http.Error(w, "Error removing like", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Like removed successfully"))
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ContentHandler) HandleComments(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	id, err := strconv.Atoi(r.URL.Path[len("/comments/"):])
	if err != nil {
		http.Error(w, "Invalid essay ID", http.StatusBadRequest)
		return
	}

	_, err = h.EssayService.GetEssayByID(uint64(id))
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

	switch r.Method {
	case http.MethodGet:
		comments, err := h.ContentService.GetComments(uint64(id))
		if err != nil {
			log.Print("Error fetching comments: ", err)
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comments)

	case http.MethodPost:
		session, _ := config.SessionStore.Get(r, "session")
		userIDInterface, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		userID := userIDInterface.(uint64)

		var comment struct {
			CommentText string `json:"comment_text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, "Invalid comment data", http.StatusBadRequest)
			return
		}
		if err := h.ContentService.AddComment(userID, uint64(id), comment.CommentText); err != nil {
			http.Error(w, "Error adding comment", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Comment added successfully"))

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GetVariant handles GET /variants/id
func (h *ContentHandler) GetVariant(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Printf("Invalid variant ID: %v", err)
		http.Error(w, "Invalid variant ID", http.StatusBadRequest)
		return
	}

	variant, err := h.ContentService.GetVariantByID(uint64(id))
	if err != nil {
		log.Print("Error getting variant: ", err)
		http.Error(w, "Error getting variant:", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variant)
}

// GetCounts handles GET /counts
func (h *ContentHandler) GetCounts(w http.ResponseWriter, r *http.Request) {
	log.Println("GET ", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	variants_count, essays_count, users_count, err := h.ContentService.GetCounts()
	if err != nil {
		log.Print("Error getting variants count: ", err)
		http.Error(w, "Error getting variants count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"variants_count": variants_count,
		"essays_count":   essays_count,
		"users_count":    users_count,
	})
}
