package repo

import (
	"path/filepath"
	"testing"

	"main.go/internal/config"
)

func TestConfig_Server(t *testing.T) {
	envPath := filepath.Join("..", "..", "..", ".env")
	cfg, err := config.New(envPath)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	db, err := NewClient(cfg.DB.DSN)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}

	if db == nil {
		t.Fatalf("expected db got nil")
	}
}

func TestConfig_Server_Invalid(t *testing.T) {
	cfg, err := config.New("")
	if err == nil {
		t.Fatalf("expected error got nil")
	}

	db, err := NewClient(cfg.DB.DSN)
	if err == nil {
		t.Fatalf("expected error got nil")
	}

	if db != nil {
		t.Fatalf("expected no db got %v", db)
	}
}
