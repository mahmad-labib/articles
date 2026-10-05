package handlers

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mahmad-labib/articles/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, domUser *domain.User) (int64, error)
	Login(ctx context.Context, email, password string, jwt *jwtauth.JWTAuth) (string, error)
	GetProfile(ctx context.Context, id int64) (*domain.Profile, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, category *domain.Category) (int64, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type TagService interface {
	CreateTag(ctx context.Context, tag *domain.Tag) (int64, error)
	ListTags(ctx context.Context, page, pageSize int) ([]domain.Tag, error)
	GetTagByID(ctx context.Context, id int64) (*domain.Tag, error)
	UpdateTag(ctx context.Context, tag *domain.Tag) error
	DeleteTag(ctx context.Context, id int64) error
}

type ArticleService interface {
	CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) (*domain.Article, error)
	GetArticleByID(ctx context.Context, id int64) (*domain.Article, error)
}

type Handler struct {
	jwt  *jwtauth.JWTAuth
	usvc UserService
	tsvc TagService
	csvc CategoryService
	asvc ArticleService
}

func NewHandler(
	jwt *jwtauth.JWTAuth,
	usvc UserService,
	tsvc TagService,
	csvc CategoryService,
	asvc ArticleService,
) *Handler {
	return &Handler{
		jwt:  jwt,
		usvc: usvc,
		tsvc: tsvc,
		csvc: csvc,
		asvc: asvc,
	}
}
