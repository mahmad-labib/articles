package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db/models"
	"gorm.io/gorm"
)

type APIArticle struct {
	ID       int64
	Title    string
	Category models.Category
	Author   models.User
	Tags     []models.Tag `gorm:"many2many:article_tags;"`
}

type ArticleFilter struct {
	CategoryID int64 `json:"category_id"`
	TagID      int64 `json:"tag_id"`
	AuthorID   int64 `json:"author_id"`
	Page       int   `json:"Page"`
	PageSize   int   `json:"pageSize"`
}

func (r *Repository) CreateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) (*domain.Article, error) {
	Article := models.Article{
		UserID:     article.AuthorID,
		CategoryID: article.CategoryID,
		Title:      article.Title,
		Slug:       strings.ToLower(article.Title),
		Body:       article.Body,
		Status:     string(article.Status),
	}

	result := r.db.WithContext(ctx).Create(&Article)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
			return nil, fmt.Errorf("failed to create new Article with err: category %w", gorm.ErrForeignKeyViolated)
		}

		return nil, fmt.Errorf("failed to create new Article with err: %w", result.Error)
	}

	Tags := []models.ArticleTags{}

	for _, t := range tags {
		err := r.db.First(&models.Tag{}, t.ID).Error
		if err == nil {
			Tags = append(Tags, models.ArticleTags{TagID: t.ID, ArticleID: Article.ID})
		}
	}

	if len(Tags) > 0 {
		result = r.db.WithContext(ctx).Create(&Tags)
		if result.Error != nil {
			err := fmt.Errorf("failed to create new ArticleTags with err: %w", result.Error)

			return nil, err
		}
	}

	newArticle, err := r.GetArticleByID(ctx, Article.ID)
	if err != nil {
		return nil, err
	}

	return newArticle, nil
}

func (r *Repository) GetArticleByID(ctx context.Context, id int64) (*domain.Article, error) {
	Article := models.Article{}

	err := r.db.Preload("Comments").Preload("Tags").First(&Article, id).Error
	if err != nil {
		return nil, err
	}

	comments := make([]domain.Comment, 0, len(Article.Comments))
	for _, comment := range Article.Comments {
		comments = append(comments, domain.Comment{
			ID:        comment.ID,
			ArticleID: comment.ArticleID,
			AuthorID:  comment.UserID,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
		})
	}

	tags := make([]domain.Tag, 0, len(Article.Tags))
	for _, tag := range Article.Tags {
		tags = append(tags, domain.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	article := domain.Article{
		ID:         Article.ID,
		AuthorID:   Article.UserID,
		CategoryID: Article.CategoryID,
		Status:     domain.ArticleStatus(Article.Status),
		Title:      Article.Title,
		Slug:       Article.Slug,
		Body:       Article.Body,
		Comments:   comments,
		Tags:       tags,
	}

	return &article, nil
}

func (r *Repository) UpdateArticle(ctx context.Context, article domain.Article, tags []domain.Tag) error {
	Article := models.Article{
		Title:      article.Title,
		Slug:       article.Slug,
		Body:       article.Body,
		Status:     string(article.Status),
		CategoryID: article.CategoryID,
	}

	res := r.db.WithContext(ctx).Model(&models.Article{}).Where("id = ?", article.ID).Where("user_id = ?", ctx.Value("userID")).Updates(Article)
	if res.Error != nil {
		return fmt.Errorf("failed to update Article with err: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("article: id=%d: does not exist", article.ID)
	}

	return nil
}

func (r *Repository) DeleteArticle(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&models.Article{}, id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) ListPublishedArticles(ctx context.Context, filter ArticleFilter) ([]domain.Article, error) {
	var Articles []models.Article

	err := r.db.
		Scopes(Paginate(Pagination{
			Page:     filter.Page,
			PageSize: filter.PageSize,
		})).
		Scopes(Filter(filter)).
		Preload("User").
		Where("status = ?", "published").
		Find(&Articles).Error
	if err != nil {
		return nil, err
	}

	articles := make([]domain.Article, 0, len(Articles))
	for _, a := range Articles {
		tags := make([]domain.Tag, 0, len(a.Tags))
		for _, tag := range a.Tags {
			tags = append(tags, domain.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			})
		}

		article := domain.Article{
			ID:         a.ID,
			Title:      a.Title,
			Tags:       tags,
			Slug:       a.Slug,
			CategoryID: a.CategoryID,
			AuthorID:   a.UserID,
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *Repository) ListProfileArticles(ctx context.Context, id int64, paginate Pagination) ([]domain.Article, error) {
	var Articles []models.Article

	err := r.db.
		Scopes(Paginate(paginate)).
		Preload("Tags").
		Preload("User").
		Where("user_id = ?", id).
		Find(&Articles).Error
	if err != nil {
		return nil, err
	}

	articles := make([]domain.Article, 0, len(Articles))
	for _, a := range Articles {
		tags := make([]domain.Tag, 0, len(a.Tags))
		for _, tag := range a.Tags {
			tags = append(tags, domain.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			})
		}

		article := domain.Article{
			ID:         a.ID,
			Title:      a.Title,
			Slug:       a.Slug,
			Tags:       tags,
			CategoryID: a.CategoryID,
			AuthorID:   a.UserID,
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (r *Repository) GetArticleListByTitle(ctx context.Context, title string) (*[]domain.Article, error) {
	var Articles []models.Article

	err := r.db.Preload("Comments").Preload("Tags").Where("slug LIKE ?", "%"+title+"%").Find(&Articles).Error
	if err != nil {
		return nil, err
	}

	articles := make([]domain.Article, 0, len(Articles))
	for _, a := range Articles {
		tags := make([]domain.Tag, 0, len(a.Tags))
		for _, tag := range a.Tags {
			tags = append(tags, domain.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			})
		}

		comments := make([]domain.Comment, 0, len(a.Comments))
		for _, comment := range a.Comments {
			comments = append(comments, domain.Comment{
				ID:        comment.ID,
				ArticleID: comment.ArticleID,
				AuthorID:  comment.UserID,
				Body:      comment.Body,
				CreatedAt: comment.CreatedAt,
			})
		}

		article := domain.Article{
			ID:         a.ID,
			Title:      a.Title,
			Slug:       a.Slug,
			Tags:       tags,
			CategoryID: a.CategoryID,
			AuthorID:   a.UserID,
			Comments:   comments,
		}
		articles = append(articles, article)
	}

	return &articles, nil
}

func (r *Repository) AddView(ctx context.Context, articleID int64, views int) error {
	// 1. Check if the article exists
	var article models.Article
	if err := r.db.WithContext(ctx).First(&article, articleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("article not found: %w", err)
		}
		return fmt.Errorf("failed to check article existence: %w", err)
	}

	// 2. Try to perform an atomic increment assuming the record exists
	res := r.db.WithContext(ctx).
		Model(&models.View{}).
		Where("article_id = ?", articleID).
		Update("views", gorm.Expr("views + ?", views))

	if res.Error != nil {
		return fmt.Errorf("failed to update view count: %w", res.Error)
	}

	// 3. If no row was affected, create the initial record
	if res.RowsAffected == 0 {
		view := models.View{
			ArticleID: articleID,
			Views:     1,
		}
		if err := r.db.WithContext(ctx).Create(&view).Error; err != nil {
			return fmt.Errorf("failed to create view record: %w", err)
		}
	}

	return nil
}
