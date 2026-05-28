package config

import "testing"

func TestLoad_Defaults(t *testing.T) {
	cfg := Load()
	if cfg.Port != "8080" {
		t.Fatalf("port: got %q", cfg.Port)
	}
	if cfg.DBName != "tiktok" {
		t.Fatalf("db name: got %q", cfg.DBName)
	}
	if cfg.DSN() == "" {
		t.Fatal("empty dsn")
	}
}
