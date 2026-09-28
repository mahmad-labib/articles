package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	Server Server
	DB     DB
	JWT    string
}

func New(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		return &Config{}, err
	}

	return &Config{
		Server: ServerValues(),
		DB:     CreateDSN(),
		JWT:    JWT(),
	}, nil
}
