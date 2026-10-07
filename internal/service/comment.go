package service

import (
	"context"
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db"
	"github.com/mahmad-labib/articles/internal/service/validation"
	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *domain.Comment) error
	ListCommentsByArticleID(ctx context.Context, paginate db.Pagination, articleId int64) ([]domain.Comment, error)
	DeleteComment(ctx context.Context, commentId int64) error
	GetCommentByID(ctx context.Context, commentId int64) (*domain.Comment, error)
	GetArticleByID(ctx context.Context, articleId int64) (*domain.Article, error)
}

type CommentValidator interface {
	Validate(comment *domain.Comment) error
}

type CommentService struct {
	repo      CommentRepository
	validator CommentValidator
}

func NewCommentService(repo CommentRepository, validator CommentValidator) *CommentService {
	return &CommentService{
		repo:      repo,
		validator: validator,
	}
}

func (s *CommentService) ListCommentsByArticleID(ctx context.Context, paginate db.Pagination, articleId int64) ([]domain.Comment, error) {
	comments, err := s.repo.ListCommentsByArticleID(ctx, paginate, articleId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) CreateComment(ctx context.Context, comment *domain.Comment) error {
	article, err := s.repo.GetArticleByID(ctx, comment.ArticleID)
	if err != nil {
		return fmt.Errorf("article not found: %w", gorm.ErrRecordNotFound)
	}

	if article.ID != comment.ArticleID {
		return fmt.Errorf("article not found: %w", gorm.ErrRecordNotFound)
	}

	err = s.validator.Validate(comment)
	if err != nil {
		return err
	}

	err = s.repo.CreateComment(ctx, comment)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentService) DeleteComment(ctx context.Context, commentId int64) error {
	comment, err := s.repo.GetCommentByID(ctx, commentId)
	if err != nil {
		return err
	}

	if comment == nil {
		return gorm.ErrRecordNotFound
	}

	if comment.AuthorID != ctx.Value("userID") {
		return validation.Validation
	}

	err = s.repo.DeleteComment(ctx, commentId)
	if err != nil {
		return err
	}

	return nil
}
