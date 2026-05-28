package service_test

import (
	"testing"

	"tiktok-app/backend/internal/config"
	"tiktok-app/backend/internal/database"
	"tiktok-app/backend/internal/service"
)

func TestVideoService_List(t *testing.T) {
	t.Skip("integration test: run in CI with MySQL")

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}

	svc := service.NewVideoService(db)
	feed, err := svc.List(0, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(feed.Items) == 0 {
		t.Fatal("expected seeded videos")
	}
}
