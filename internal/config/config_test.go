package config

import (
	"path/filepath"
	"testing"
)

func TestConfig_Server(t *testing.T) {
	envPath := filepath.Join("..", "..", ".env")
	cfg, err := New(envPath)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if cfg.Server.Port == "" {
		t.Fatalf("expected Port got nil")
	}
	if cfg.Server.IdleTimeout == 0 {
		t.Fatalf("expected IdleTimeout got nil")
	}
	if cfg.Server.ReadTimeout == 0 {
		t.Fatalf("expected ReadTimeout got nil")
	}
	if cfg.Server.WriteTimeout == 0 {
		t.Fatalf("expected WriteTimeout got nil")
	}

}

func TestConfig_DB(t *testing.T) {
	envPath := filepath.Join("..", "..", ".env")
	cfg, err := New(envPath)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if cfg.DB.DSN == "" {
		t.Fatalf("expected DSN got nil")
	}
}
