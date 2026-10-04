package db

import "gorm.io/gorm"

func Paginate(page, PageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
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
