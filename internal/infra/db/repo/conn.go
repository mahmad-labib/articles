package repo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	return db, nil
}
