package app

import (
	"essay/src/internal/config"
	"essay/src/internal/database"
	"essay/src/internal/middleware"
	"essay/src/internal/services"
	"essay/src/internal/transport/handlers"
	"net/http"
)

type App struct {
	DB *database.DB

	UserService *services.UserService

	UserHandler *handlers.UserHandler
}

func NewApp() *App {
	db := database.GetPostgreSQLConnection()

	userService := services.NewUserService(db.Instance)

	userHandler := handlers.NewUserHandler(userService)

	return &App{
		DB: db,

		UserService: userService,

		UserHandler: userHandler,
	}
}

func (a *App) Close() {
	a.DB.Close()
}

func (a *App) ServeMux() http.Handler {
	config.InitSessionStore()
	mux := http.NewServeMux()

	a.UserHandler.RegisterRoutes(mux)

	return middleware.NewCORSMiddleware(config.СorsConfig)(mux)
}
