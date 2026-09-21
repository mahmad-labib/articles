package domain

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	AuthorID  int64     `json:"author_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
