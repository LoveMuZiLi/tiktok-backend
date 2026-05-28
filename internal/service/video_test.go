package service_test

import (
	"testing"

	"github.com/LoveMuZiLi/tiktok-backend/internal/config"
	"github.com/LoveMuZiLi/tiktok-backend/internal/infra/persistence"
	"github.com/LoveMuZiLi/tiktok-backend/internal/service"
)

func TestVideoService_List(t *testing.T) {
	t.Skip("integration test: run in CI with MySQL")

	cfg := config.Load()
	db, err := persistence.Connect(cfg)
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
