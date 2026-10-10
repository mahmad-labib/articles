package handlers

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db"
	"github.com/mahmad-labib/articles/internal/service"
)

type UserService interface {
	CreateUser(ctx context.Context, domUser *domain.User) (int64, error)
	Login(ctx context.Context, email, password string, jwt *jwtauth.JWTAuth) (string, error)
	GetProfile(ctx context.Context, id int64) (*domain.Profile, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, category *domain.Category) (int64, error)
	ListCategories(ctx context.Context, paginate db.Pagination) ([]domain.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	GetCategoryListBySlug(ctx context.Context, slug string) (*[]domain.Category, error)
}

type TagService interface {
	CreateTag(ctx context.Context, tag *domain.Tag) (int64, error)
	ListTags(ctx context.Context, paginate db.Pagination) ([]domain.Tag, error)
	GetTagByID(ctx context.Context, id int64) (*domain.Tag, error)
	UpdateTag(ctx context.Context, tag *domain.Tag) error
	DeleteTag(ctx context.Context, id int64) error
	GetTagsByName(ctx context.Context, name string) (*[]domain.Tag, error)
}

type ArticleService interface {
	CreateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) (*domain.Article, error)
	GetArticleByID(ctx context.Context, id int64) (*domain.Article, error)
	ListPublishedArticles(ctx context.Context, filter db.ArticleFilter) ([]domain.Article, error)
	ListProfileArticles(ctx context.Context, id int64, paginate db.Pagination) ([]domain.Article, error)
	UpdateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) (*domain.Article, error)
	DeleteArticle(ctx context.Context, id int64) error
	AddView(ctx context.Context, articleID int64, views int) error
	GetArticleListByTitle(ctx context.Context, title string) (*[]domain.Article, error)
}

type CommentService interface {
	ListCommentsByArticleID(ctx context.Context, paginate db.Pagination, articleId int64) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	DeleteComment(ctx context.Context, commentId int64) error
}

type Handler struct {
	views *service.BatchCounter[int64]
	jwt   *jwtauth.JWTAuth
	usvc  UserService
	tsvc  TagService
	csvc  CategoryService
	asvc  ArticleService
	cosvc CommentService
}

func NewHandler(
	views *service.BatchCounter[int64],
	jwt *jwtauth.JWTAuth,
	usvc UserService,
	tsvc TagService,
	csvc CategoryService,
	asvc ArticleService,
	cosvc CommentService,
) *Handler {
	return &Handler{
		views: views,
		jwt:   jwt,
		usvc:  usvc,
		tsvc:  tsvc,
		csvc:  csvc,
		asvc:  asvc,
		cosvc: cosvc,
	}
}
