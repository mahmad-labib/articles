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

func (r *Repository) CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) (*domain.Article, error) {
	Article := models.Article{
		UserID:     userId,
		CategoryID: catId,
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

	if Tags != nil {
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

	err := r.db.Preload("User").Preload("Comments").Preload("Category").Preload("Tags").First(&Article, id).Error
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

	category := domain.Category{
		ID:   Article.CategoryID,
		Name: Article.Category.Name,
		Slug: Article.Category.Slug,
	}

	author := domain.Author{
		ID:    Article.User.ID,
		Name:  Article.User.Name,
		Email: Article.User.Email,
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
		Category:   category,
		Author:     author,
	}

	return &article, nil
}

func (r *Repository) UpdateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error {
	Article := models.Article{
		ID:         article.ID,
		Title:      article.Title,
		Slug:       article.Slug,
		Body:       article.Body,
		Status:     string(article.Status),
		CategoryID: catId,
	}

	res := r.db.WithContext(ctx).Model(&models.Article{}).Where("id = ?", article.ID).Updates(Article)
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

func (r *Repository) ListArticles(ctx context.Context, Page, PageSize int) ([]domain.Article, error) {
	var Articles []APIArticle

	err := r.db.
		Scopes(Paginate(Page, PageSize)).
		Preload("Tags").
		Preload("Category").
		Preload("Author").
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
			ID:    a.ID,
			Title: a.Title,
			Tags:  tags,
			Category: domain.Category{
				ID:   a.Category.ID,
				Name: a.Category.Name,
			},
			Author: domain.Author{
				ID:   a.Author.ID,
				Name: a.Author.Name,
			},
		}
		articles = append(articles, article)
	}

	return articles, nil
}
