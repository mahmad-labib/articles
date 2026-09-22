package models

type ArticleTags struct {
	ArticleID int64
	TagID     int64

	Article Article
	Tag     Tag
}
