package config

import (
	"fmt"
	"os"
)

type DB struct {
	DSN string
}

func CreateDSN() DB {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	DB_options := DB{}
	DB_options.DSN = fmt.Sprintf("postgres://%s:%s@%v:%s/%s", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

	return DB_options
}
