package service

import (
	"context"

	"github.com/mahmad-labib/articles/internal/domain"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) (*domain.Article, error)
	ListArticles(ctx context.Context, Page, PageSize int) ([]domain.Article, error)
	DeleteArticle(ctx context.Context, id int64) error
	GetArticleByID(ctx context.Context, id int64) (*domain.Article, error)
	UpdateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error
}

type ArticleService struct {
	repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) *ArticleService {
	return &ArticleService{
		repo: repo,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) (*domain.Article, error) {
	articleModel, err := s.repo.CreateArticle(ctx, article, userId, catId, tags)
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
