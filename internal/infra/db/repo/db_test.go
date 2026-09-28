package repo

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"main.go/internal/config"
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
