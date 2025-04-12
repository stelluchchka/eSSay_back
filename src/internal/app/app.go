package app

import (
	"essay/src/internal/config"
	"essay/src/internal/database"
	"essay/src/internal/middleware"
	"essay/src/internal/services"
	"essay/src/internal/transport/handlers"
	"log"
	"net/http"
	"time"
)

type App struct {
	DB *database.DB

	UserService *services.UserService

	UserHandler *handlers.UserHandler

	stopChan chan struct{} // для graceful shutdown
}

func NewApp() *App {
	db := database.GetPostgreSQLConnection()

	userService := services.NewUserService(db.Instance)

	userHandler := handlers.NewUserHandler(userService)

	app := &App{
		DB:          db,
		UserService: userService,
		UserHandler: userHandler,
		stopChan:    make(chan struct{}),
	}

	// Запускаем периодический сброс проверок
	go app.startCheckResetter()

	return app
}

func (a *App) startCheckResetter() {
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := a.UserService.ResetAllChecks(); err != nil {
				log.Printf("Error resetting check counts: %v", err)
			} else {
				log.Print("Successfully reset all user check counts")
			}
		case <-a.stopChan:
			return
		}
	}
}

func (a *App) Close() {
	close(a.stopChan) // останавливаем горутину сброса проверок
	a.DB.Close()
}

func (a *App) ServeMux() http.Handler {
	config.InitSessionStore()
	mux := http.NewServeMux()

	a.UserHandler.RegisterRoutes(mux)

	return middleware.NewCORSMiddleware(config.СorsConfig)(mux)
}
