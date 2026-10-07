package db

import (
	"context"
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db/models"
)

func (r *Repository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	Comment := models.Comment{
		UserID:    comment.AuthorID,
		ArticleID: comment.ArticleID,
		Body:      comment.Body,
	}

	result := r.db.WithContext(ctx).Create(&Comment)
	if result.Error != nil {
		return fmt.Errorf("failed to create new Comment with err: %w", result.Error)
	}

	return nil
}

func (r *Repository) ListCommentsByArticleID(ctx context.Context, paginate Pagination, articleId int64) ([]domain.Comment, error) {
	var comments []models.Comment
	result := r.db.WithContext(ctx).Where("article_id = ?", articleId).Scopes(Paginate(paginate)).Find(&comments)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list Comments with err: %w", result.Error)
	}

	domainComments := make([]domain.Comment, len(comments))
	for i, comment := range comments {
		domainComments[i] = domain.Comment{
			ID:        comment.ID,
			Body:      comment.Body,
			AuthorID:  comment.UserID,
			ArticleID: comment.ArticleID,
			CreatedAt: comment.CreatedAt,
		}
	}

	return domainComments, nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentId int64) error {
	result := r.db.WithContext(ctx).Delete(&models.Comment{}, commentId)
	if result.Error != nil {
		return fmt.Errorf("failed to delete Comment with err: %w", result.Error)
	}

	return nil
}

func (r *Repository) GetCommentByID(ctx context.Context, commentId int64) (*domain.Comment, error) {
	var comment models.Comment
	result := r.db.WithContext(ctx).First(&comment, commentId)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get Comment with err: %w", result.Error)
	}

	domainComment := &domain.Comment{
		ID:        comment.ID,
		ArticleID: comment.ArticleID,
		AuthorID:  comment.UserID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
	}

	return domainComment, nil
}
