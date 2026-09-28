package service

import (
	"context"

	"main.go/internal/domain"
)

type ArticlesRepository interface {
	CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error
	ListArticles(ctx context.Context, Page, PageSize int) ([]domain.Article, error)
	DeleteArticle(ctx context.Context, id int64) error
	GetArticleByID(ctx context.Context, id int64) (*domain.Article, error)
	UpdateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error
	CreateCategory(ctx context.Context, category domain.Category) error
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, category domain.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	GetCategoryWithArticles(ctx context.Context, id int64) (*domain.Category, error)
	CreateComment(ctx context.Context, comment domain.Comment, userId int64, articleId int64) error
	DeleteComment(ctx context.Context, commentId int64) error
	CreateUser(ctx context.Context, domUser *domain.User) (int64, error)
	GetUserByID(ctx context.Context, ID int64) *domain.ApiUser
	UpdateUser(ctx context.Context, domUser *domain.ApiUser) error
	ListUser(ctx context.Context, Page, PageSize int) ([]domain.ApiUser, error)
	DeleteUser(ctx context.Context, ID int64) error
	GetUserByEmail(ctx context.Context, Email string) *domain.ApiUser
	GetUserByEmail_Login(ctx context.Context, Email string) string
	GetProfile(ctx context.Context, id int64) (*domain.Profile, error)
}

type Validator interface {
	ValidateUser(item *domain.User) error
}

type ArticleService struct {
	repo      ArticlesRepository
	validator Validator
}

func NewArticleService(repo ArticlesRepository, validator Validator) *ArticleService {
	return &ArticleService{
		repo:      repo,
		validator: validator,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error {
	err := s.repo.CreateArticle(ctx, article, userId, catId, tags)
	if err != nil {
		return err
	}
	return nil
}
