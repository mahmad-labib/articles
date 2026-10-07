package db

import (
	"context"
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db/models"
)

func (r *Repository) CreateTag(ctx context.Context, tag *domain.Tag) (int64, error) {
	Tag := models.Tag{
		Name: tag.Name,
	}

	result := r.db.WithContext(ctx).Create(&Tag)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new Tag with: %w", result.Error)

		return 0, err
	}

	return Tag.ID, nil
}

func (r *Repository) GetTagByID(ctx context.Context, id int64) (*domain.Tag, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %d", id)
	}

	Tag := models.Tag{}
	Item := domain.Tag{}

	err := r.db.WithContext(ctx).First(&Tag, id).Error
	if err != nil {
		return &Item, err
	}

	Item.ID = Tag.ID
	Item.Name = Tag.Name

	return &Item, nil
}

func (r *Repository) ListTags(ctx context.Context, paginate Pagination) ([]domain.Tag, error) {
	tags := []models.Tag{}
	items := []domain.Tag{}

	err := r.db.WithContext(ctx).Scopes(Paginate(paginate)).Find(&tags).Error
	if err != nil {
		return items, err
	}

	for _, tag := range tags {
		items = append(items, domain.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return items, nil
}

func (r *Repository) UpdateTag(ctx context.Context, tag *domain.Tag) error {
	if tag.ID <= 0 {
		return fmt.Errorf("invalid id: %d", tag.ID)
	}

	Tag := models.Tag{
		Name: tag.Name,
	}

	result := r.db.WithContext(ctx).Model(&Tag).Where("id = ?", tag.ID).Updates(Tag)
	if result.Error != nil {
		err := fmt.Errorf("failed to update Tag with: %w", result.Error)

		return err
	}

	return nil
}

func (r *Repository) DeleteTag(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid id: %d", id)
	}

	Tag := models.Tag{}

	err := r.db.WithContext(ctx).First(&Tag, id).Error
	if err != nil {
		return fmt.Errorf("failed to find Tag with id: %d", id)
	}

	result := r.db.WithContext(ctx).Delete(&Tag, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete Tag with id: %d", id)
	}

	return nil
}
