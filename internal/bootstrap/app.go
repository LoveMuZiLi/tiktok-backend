package bootstrap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LoveMuZiLi/tiktok-backend/internal/config"
	"github.com/LoveMuZiLi/tiktok-backend/internal/infra/persistence"
	httpx "github.com/LoveMuZiLi/tiktok-backend/internal/transport/http"
)

func Run() error {
	cfg := config.Load()

	db, err := persistence.Connect(cfg)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	handler := httpx.NewRouter(cfg, db)
	addr := ":" + cfg.Port
	log.Printf("tiktok-api listening on %s (mysql %s@%s:%s/%s)", addr, cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return http.ListenAndServe(addr, handler)
}
