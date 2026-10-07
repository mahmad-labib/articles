package db

import (
	"gorm.io/gorm"
)

type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Filter   string `json:"filter"`
}

func Paginate(paginate Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := paginate.Page
		PageSize := paginate.PageSize

		if page <= 0 {
			page = 1
		}

		switch {
		case PageSize > 100:
			PageSize = 100
		case PageSize <= 0:
			PageSize = 10
		}

		offset := (page - 1) * PageSize

		return db.Offset(offset).Limit(PageSize)
	}
}

func Filter(filter ArticleFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.CategoryID != 0 {
			db = db.Where("articles.category_id = ?", &filter.CategoryID)
		}
		if filter.AuthorID != 0 {
			db = db.Where("articles.user_id = ?", &filter.AuthorID)
		}
		if filter.TagID != 0 {
			db = db.
				Joins("JOIN article_tags ON article_tags.article_id = articles.id").
				Where("article_tags.tag_id = ?", &filter.TagID)
		}
		return db
	}
}
