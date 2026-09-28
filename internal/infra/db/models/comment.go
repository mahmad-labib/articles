package models

import "time"

type Comment struct {
	ID        int64
	ArticleID int64
	UserID    int64
	Body      string
	CreatedAt time.Time

	Article Article
	User    User
}
