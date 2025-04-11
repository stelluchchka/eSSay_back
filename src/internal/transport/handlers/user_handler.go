package handlers

import (
	"essay/src/internal/services"
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
	// user
	mux.HandleFunc("/users/nickname", h.GetNickname)
	mux.HandleFunc("/users/login", h.HandleLogin)
	mux.HandleFunc("/users/logout", h.HandleLogout)
	mux.HandleFunc("/users/info", h.HandleUserInfo)
	mux.HandleFunc("/users", h.HandleUser)

	// content
	mux.HandleFunc("/counts/", h.GetCounts)
	mux.HandleFunc("/likes/is_liked/", h.HandleIsLiked)
	mux.HandleFunc("/likes/", h.HandleLikes)
	mux.HandleFunc("/comments/", h.HandleComments)
	mux.HandleFunc("/variants", h.CreateVariant)
	mux.HandleFunc("/variants/", h.GetVariant)
	mux.HandleFunc("/criteria", h.GetCriteria)
	// result
	mux.HandleFunc("/result/", h.CreateResult)

	// essay
	mux.HandleFunc("/essays", h.HandleEssaysRequests)
	mux.HandleFunc("/essays/", h.HandleEssayRequests)
	mux.HandleFunc("/essays/appeal", h.GetAppealEssays)
	mux.HandleFunc("/users/me/essays", h.GetUserEssays)
	mux.HandleFunc("/users/me/essays/", h.GetUserEssayByID)
}
