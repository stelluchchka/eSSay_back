package app

import (
	"essay/src/internal/database"
	"essay/src/internal/services"
	"essay/src/internal/transport/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	UserService  *services.UserService
	EssayService *services.EssayService

	UserHandler  *handlers.UserHandler
	EssayHandler *handlers.EssayHandler
}

func NewApp() *App {
	userService := services.NewUserService()
	essayService := services.NewEssayService()

	userHandler := handlers.NewUserHandler(userService)
	essayHandler := handlers.NewEssayHandler(essayService)

	return &App{
		UserService:  userService,
		EssayService: essayService,

		UserHandler:  userHandler,
		EssayHandler: essayHandler,
	}
}

func (a *App) Start() {
	mux := http.NewServeMux()

	a.UserHandler.RegisterRoutes(mux)
	a.EssayHandler.RegisterRoutes(mux)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :8080")
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			log.Fatal("Error starting server:", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	database.CloseDB()

	log.Println("Server stopped")
}
