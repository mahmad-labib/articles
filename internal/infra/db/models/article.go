package models

type Article struct {
	ID         int64
	Title      string
	Slug       string
	Body       string
	Status     string
	AuthorID   int64
	CategoryID int64

	Author   User
	Category Category
	Comments []Comment
	Tags     []Tag `gorm:"many2many:article_tags;"`
}
