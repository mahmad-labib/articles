package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/domain"
	"main.go/internal/infra/db/models"
)

type APIArticle struct {
	ID       int64
	Title    string
	Category models.Category
	Author   models.User
	Tags     []models.Tag `gorm:"many2many:article_tags;"`
}

func (C *Core) CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error {
	Article := models.Article{
		AuthorID:   userId,
		CategoryID: catId,
		Title:      article.Title,
		Slug:       article.Slug,
		Body:       article.Body,
		Status:     string(article.Status),
	}

	Tags := []models.Tag{}
	for _, t := range tags {
		Tags = append(Tags, models.Tag{ID: t.ID})
	}

	result := C.db.WithContext(ctx).Create(&Article)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new Article with err: %w", result.Error)
		return err
	}

	if err := C.db.Model(&Article).Association("Tags").Append(Tags); err != nil {
		return fmt.Errorf("failed to add Tags to Article with err: %w", result.Error)
	}

	return nil
}

func (C *Core) GetArticleByID(ctx context.Context, id int64) (*domain.Article, error) {
	Article := models.Article{}

	err := C.db.Preload("Comments").Preload("Tags").First(&Article, id).Error
	if err != nil {
		return nil, err
	}

	comments := make([]domain.Comment, 0, len(Article.Comments))
	for _, comment := range Article.Comments {
		comments = append(comments, domain.Comment{
			ID:        comment.ID,
			ArticleID: comment.ArticleID,
			AuthorID:  comment.AuthorID,
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
		ID:       Article.ID,
		Title:    Article.Title,
		Slug:     Article.Slug,
		Body:     Article.Body,
		Comments: comments,
		Tags:     tags,
	}

	return &article, nil
}

func (C *Core) UpdateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error {
	Article := models.Article{
		ID:         article.ID,
		Title:      article.Title,
		Slug:       article.Slug,
		Body:       article.Body,
		Status:     string(article.Status),
		CategoryID: catId,
	}

	res := C.db.WithContext(ctx).Model(&models.Article{}).Where("id = ?", article.ID).Updates(Article)
	if res.Error != nil {
		return fmt.Errorf("failed to update Article with err: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("article: id=%d: does not exist", article.ID)
	}

	return nil
}

func (C *Core) DeleteArticle(ctx context.Context, id int64) error {
	res := C.db.WithContext(ctx).Delete(&models.Article{}, id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (C *Core) ListArticles(ctx context.Context, Page, PageSize int) ([]domain.Article, error) {
	var Articles []APIArticle
	err := C.db.Scopes(Paginate(Page, PageSize)).Preload("Tags").Find(&Articles).Error
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
		}
		articles = append(articles, article)
	}

	return articles, nil
}
