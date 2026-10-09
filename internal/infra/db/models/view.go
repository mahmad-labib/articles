package models

type View struct {
	ArticleID int64
	Views     int64

	Article Article
}
