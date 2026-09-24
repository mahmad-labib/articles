package db

import (
	"context"
	"fmt"

	"main.go/internal/domain"
	"main.go/internal/infra/db/models"
)

func (C *Core) CreateComment(ctx context.Context, comment domain.Comment, userId int64, articleId int64) error {
	Comment := models.Comment{
		AuthorID:  userId,
		ArticleID: articleId,
		Body:      comment.Body,
	}

	result := C.db.WithContext(ctx).Create(&Comment)
	if result.Error != nil {
		return fmt.Errorf("failed to create new Comment with err: %w", result.Error)
	}

	return nil
}

func (C *Core) DeleteComment(ctx context.Context, commentId int64) error {
	result := C.db.WithContext(ctx).Delete(&models.Comment{}, commentId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete Comment with err: %w", result.Error)
	}

	return nil
}
