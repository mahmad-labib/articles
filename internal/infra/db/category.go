package db

import (
	"context"
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db/models"
)

func (r *Repository) CreateCategory(ctx context.Context, category *domain.Category) (int64, error) {
	Category := models.Category{
		Name: category.Name,
		Slug: category.Slug,
	}

	result := r.db.WithContext(ctx).Create(&Category)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new Category with: %w", result.Error)

		return 0, err
	}

	return Category.ID, nil
}

func (r *Repository) GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %d", id)
	}

	Category := models.Category{}
	Item := domain.Category{}

	err := r.db.First(&Category, id).Error
	if err != nil {
		return &Item, err
	}

	Item = domain.Category{
		ID:   Category.ID,
		Name: Category.Name,
	}

	return &Item, nil
}

func (r *Repository) ListCategories(ctx context.Context, paginate Pagination) ([]domain.Category, error) {
	categories := []models.Category{}
	items := []domain.Category{}

	err := r.db.WithContext(ctx).Scopes(Paginate(paginate)).Find(&categories).Error
	if err != nil {
		return items, err
	}

	for _, category := range categories {
		items = append(items, domain.Category{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		})
	}

	return items, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	Category := models.Category{
		Name: category.Name,
		Slug: category.Slug,
	}

	res := r.db.WithContext(ctx).Model(&models.Category{}).Where("id = ?", category.ID).Updates(Category)
	if res.Error != nil {
		return fmt.Errorf("error when update Category with id=%d: %w", category.ID, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("Category: id=%d: does not exist", category.ID)
	}

	return nil
}

func (r *Repository) DeleteCategory(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&models.Category{}, id)
	if res.Error != nil {
		return fmt.Errorf("error when deleting Category with id=%d: %w", id, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("Category: id=%d: does not exist", id)
	}

	return nil
}

func (r *Repository) GetCategoryWithArticles(ctx context.Context, id int64) (*domain.Category, error) {
	Category := models.Category{}
	Item := domain.Category{}

	err := r.db.Preload("Articles").First(&Category, id).Error
	if err != nil {
		return &Item, err
	}

	Item.Name = Category.Name
	Item.ID = Category.ID

	Item.Articles = make([]domain.Article, 0, len(Category.Articles))
	for _, article := range Category.Articles {
		Item.Articles = append(Item.Articles, domain.Article{
			ID:    article.ID,
			Title: article.Title,
		})
	}

	return &Item, nil
}

func (r *Repository) GetCategoryListBySlug(ctx context.Context, slug string) (*[]domain.Category, error) {
	Category := []models.Category{}
	Item := []domain.Category{}

	err := r.db.Where("slug LIKE ?", "%"+slug+"%").Find(&Category).Error
	if err != nil {
		return &Item, err
	}

	for _, cat := range Category {
		Item = append(Item, domain.Category{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	return &Item, nil
}
