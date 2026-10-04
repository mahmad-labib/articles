package db

import (
	"context"
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db/models"
)

func (r *Repository) CreateComment(ctx context.Context, comment domain.Comment, userId int64, articleId int64) error {
	Comment := models.Comment{
		UserID:    userId,
		ArticleID: articleId,
		Body:      comment.Body,
	}

	result := r.db.WithContext(ctx).Create(&Comment)
	if result.Error != nil {
		return fmt.Errorf("failed to create new Comment with err: %w", result.Error)
	}

	return nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentId int64) error {
	result := r.db.WithContext(ctx).Delete(&models.Comment{}, commentId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete Comment with err: %w", result.Error)
	}

	return nil
}
