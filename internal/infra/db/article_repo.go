package db

import (
	"context"
	"fmt"

	"main.go/internal/domain"
	"main.go/internal/infra/db/models"
)

func (C *Core) CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error {
	Article := models.Article{
		AuthorID:   userId,
		CategoryID: catId,
		Title:      article.Title,
		Slug:       article.Slug,
		Body:       article.Body,
		Status:     string(article.Status),
	}

	Tags := []models.Tag{}
	for _, t := range tags {
		Tags = append(Tags, models.Tag{ID: t.ID})
	}

	result := C.db.WithContext(ctx).Create(&Article)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new Article with err: %w", result.Error)
		return err
	}

	if err := C.db.Model(&Article).Association("Tags").Append(Tags); err != nil {
		return fmt.Errorf("failed to add Tags to Article with err: %w", result.Error)
	}

	return nil
}
