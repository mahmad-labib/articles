package models

type Tag struct {
	ID   int64
	Name string

	Articles []Article `gorm:"many2many:article_tags;"`
}
