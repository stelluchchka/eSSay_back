package app

import (
	"essay/src/internal/services"
	"essay/src/internal/transport/rest"
	"log"
	"net/http"
)

type App struct {
	UserService *services.UserService
	UserHandler *rest.UserHandler
}

func NewApp() *App {
	userService := services.NewUserService()

	userHandler := rest.NewUserHandler(userService)

	return &App{
		UserService: userService,
		UserHandler: userHandler,
	}
}

func (a *App) Start() {
	mux := http.NewServeMux()
	a.UserHandler.RegisterRoutes(mux)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
