package repo

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/gorm"
)

type Core struct {
	db *gorm.DB
}

func NewArticlesDB(db *gorm.DB) *Core {
	return &Core{db: db}
}
