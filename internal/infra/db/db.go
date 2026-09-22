package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// if err := db.Ping(); err != nil {
	// 	db.Close()
	// 	return nil, err
	// }

	return db, nil
}
