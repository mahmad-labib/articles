package models

type Category struct {
	ID   int64
	Name string
	Slug string

	Articles []Article
}
