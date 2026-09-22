package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/gorm"
)

type ArticleDB struct {
	db *gorm.DB
}

func NewArticlesDB(db *gorm.DB) *ArticleDB {
	return &ArticleDB{db: db}
}
