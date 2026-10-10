package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID          int64
	Title       string
	Slug        string
	Body        string
	Status      string
	UserID      int64
	CategoryID  int64
	DeletedAt   gorm.DeletedAt
	PublishedAt time.Time

	User     User
	Category Category
	Comments []Comment
	Tags     []Tag `gorm:"many2many:article_tags;"`
}
