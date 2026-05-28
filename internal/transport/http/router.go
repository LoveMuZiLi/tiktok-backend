package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"

	"github.com/LoveMuZiLi/tiktok-backend/internal/config"
	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
	"github.com/LoveMuZiLi/tiktok-backend/internal/transport/http/handler"
)

func NewRouter(cfg config.Config, db *gorm.DB) http.Handler {
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

	users := service.NewUserService(db)
	videos := service.NewVideoService(db)
	follows := service.NewFollowService(db)
	inbox := service.NewInboxService(db)
	notifications := service.NewNotificationService(db)

	health := handler.NewHealthHandler(db)
	userH := handler.NewUserHandler(users)
	videoH := handler.NewVideoHandler(videos)
	followH := handler.NewFollowHandler(follows)
	inboxH := handler.NewInboxHandler(inbox)
	notifyH := handler.NewNotificationHandler(notifications)

	r.Get("/api/v1/health", health.Health)

	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", userH.List)
		r.Post("/", userH.Create)
		r.Get("/{id}", userH.Get)
		r.Get("/{id}/profile", userH.Profile)
		r.Put("/{id}", userH.Update)
		r.Delete("/{id}", userH.Delete)
	})

	r.Route("/api/v1/videos", func(r chi.Router) {
		r.Get("/", videoH.List)
		r.Post("/", videoH.Create)
		r.Get("/{id}", videoH.Get)
		r.Put("/{id}", videoH.Update)
		r.Delete("/{id}", videoH.Delete)
		r.Post("/{id}/like", videoH.Like)
	})

	r.Route("/api/v1/follows", func(r chi.Router) {
		r.Post("/", followH.Create)
		r.Get("/{followingId}/status", followH.Status)
		r.Delete("/{followingId}", followH.Delete)
	})

	r.Route("/api/v1/conversations", func(r chi.Router) {
		r.Get("/", inboxH.ListChats)
		r.Post("/open", inboxH.OpenChat)
		r.Get("/{id}/messages", inboxH.ListMessages)
		r.Post("/{id}/messages", inboxH.SendMessage)
		r.Delete("/{id}/messages/{messageId}", inboxH.DeleteMessage)
	})

	r.Route("/api/v1/notifications", func(r chi.Router) {
		r.Get("/", notifyH.List)
		r.Post("/", notifyH.Create)
		r.Patch("/{id}/read", notifyH.MarkRead)
		r.Delete("/{id}", notifyH.Delete)
	})

	return r
}
