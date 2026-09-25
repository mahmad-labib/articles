package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func envPath(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file path")
	}

	dir := filepath.Dir(file)
	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatal("could not find .env in project tree")
	return ""
}

func TestConfig_Server(t *testing.T) {
	envPath := envPath(t)
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
