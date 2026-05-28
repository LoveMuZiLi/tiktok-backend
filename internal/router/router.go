package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"tiktok-app/backend/internal/config"
	"tiktok-app/backend/internal/handler"
	"tiktok-app/backend/internal/service"
)

func New(cfg config.Config, db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	videos := service.NewVideoService(db)
	health := handler.NewHealthHandler(db)
	videoH := handler.NewVideoHandler(videos)

	r.Get("/api/v1/health", health.Health)
	r.Route("/api/v1/videos", func(r chi.Router) {
		r.Get("/", videoH.List)
		r.Get("/{id}", videoH.Get)
		r.Post("/{id}/like", videoH.Like)
	})

	return r
}
