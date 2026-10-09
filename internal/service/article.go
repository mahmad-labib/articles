package service

import (
	"context"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db"
	"github.com/mahmad-labib/articles/internal/service/validation"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) (*domain.Article, error)
	ListPublishedArticles(ctx context.Context, filter db.ArticleFilter) ([]domain.Article, error)
	ListProfileArticles(ctx context.Context, id int64, paginate db.Pagination) ([]domain.Article, error)
	DeleteArticle(ctx context.Context, id int64) error
	GetArticleByID(ctx context.Context, id int64) (*domain.Article, error)
	UpdateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) error
	ListCommentsByArticleID(ctx context.Context, paginate db.Pagination, articleId int64) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
	AddView(ctx context.Context, articleID int64, views int) error
}

type ArticleService struct {
	repo      ArticleRepository
	validator ArticleValidator
}

type ArticleValidator interface {
	Validate(article *domain.Article) error
}

func NewArticleService(repo ArticleRepository, validator ArticleValidator) *ArticleService {
	return &ArticleService{
		repo:      repo,
		validator: validator,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) (*domain.Article, error) {
	err := s.validator.Validate(&article)
	if err != nil {
		return nil, err
	}

	articleModel, err := s.repo.CreateArticle(ctx, article, tags)
	if err != nil {
		return nil, err
	}

	return articleModel, nil
}

func (s *ArticleService) GetArticleByID(ctx context.Context, id int64) (*domain.Article, error) {
	article, err := s.repo.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ArticleService) ListPublishedArticles(ctx context.Context, filter db.ArticleFilter) ([]domain.Article, error) {
	articles, err := s.repo.ListPublishedArticles(ctx, filter)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *ArticleService) ListProfileArticles(ctx context.Context, id int64, paginate db.Pagination) ([]domain.Article, error) {
	articles, err := s.repo.ListProfileArticles(ctx, id, paginate)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *ArticleService) UpdateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) (*domain.Article, error) {
	err := s.validator.Validate(&article)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateArticle(ctx, article, tags)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, id int64) error {
	article, err := s.repo.GetArticleByID(ctx, id)
	if err != nil {
		return err
	}

	if article.AuthorID != ctx.Value("userID") {
		return validation.Validation
	}

	err = s.repo.DeleteArticle(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArticleService) AddView(ctx context.Context, articleID int64, views int) error {
	err := s.repo.AddView(ctx, articleID, views)
	if err != nil {
		return err
	}

	return nil
}
