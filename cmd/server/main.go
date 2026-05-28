package main

import (
	"log"
	"net/http"

	"tiktok-app/backend/internal/config"
	"tiktok-app/backend/internal/database"
	"tiktok-app/backend/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	r := router.New(cfg, db)

	addr := ":" + cfg.Port
	log.Printf("tiktok-api listening on %s (mysql %s@%s:%s/%s)", addr, cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
