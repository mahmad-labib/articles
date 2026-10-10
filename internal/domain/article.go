package domain

import "time"

type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusScheduled ArticleStatus = "scheduled"
	ArticleStatusPublished ArticleStatus = "published"
)

type Article struct {
	ID          int64
	Title       string
	Slug        string
	Body        string
	Status      ArticleStatus
	AuthorID    int64
	CategoryID  int64
	PublishedAt time.Time
	CreatedAt   time.Time

	Comments []Comment
	Tags     []Tag
}

type ProfileArticle struct {
	ID    int64
	Title string
	Slug  string
}
