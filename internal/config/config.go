package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

type Config struct {
	Server Server
}

func New(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		slog.Warn("No .env file found, will use default values")
	}

	return &Config{
		Server: ServerValues(),
	}, nil
}
