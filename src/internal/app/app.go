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

	UserService    *services.UserService
	EssayService   *services.EssayService
	ContentService *services.ContentService

	UserHandler    *handlers.UserHandler
	EssayHandler   *handlers.EssayHandler
	ContentHandler *handlers.ContentHandler
}

func NewApp() *App {
	db := database.GetPostgreSQLConnection()

	userService := services.NewUserService(db.Instance)
	essayService := services.NewEssayService(db.Instance)
	contentService := services.NewContentService(db.Instance)

	userHandler := handlers.NewUserHandler(userService)
	essayHandler := handlers.NewEssayHandler(essayService)
	contentHandler := handlers.NewContentHandler(contentService, essayService)

	return &App{
		DB: db,

		UserService:    userService,
		EssayService:   essayService,
		ContentService: contentService,

		UserHandler:    userHandler,
		EssayHandler:   essayHandler,
		ContentHandler: contentHandler,
	}
}

func (a *App) Close() {
	a.DB.Close()
}

func (a *App) ServeMux() http.Handler {
	config.InitSessionStore()
	mux := http.NewServeMux()

	a.UserHandler.RegisterRoutes(mux)
	a.EssayHandler.RegisterRoutes(mux)
	a.ContentHandler.RegisterRoutes(mux)

	return middleware.NewCORSMiddleware(config.СorsConfig)(mux)
}
