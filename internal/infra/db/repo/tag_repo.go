package repo

import (
	"context"
	"fmt"

	"main.go/internal/domain"
	"main.go/internal/infra/db/models"
)

func (c *Core) CreateTag(ctx context.Context, tag *domain.Tag) (int64, error) {
	Tag := models.Tag{
		Name: tag.Name,
	}

	result := c.db.WithContext(ctx).Create(&Tag)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new Tag with: %w", result.Error)
		return 0, err
	}

	return Tag.ID, nil
}

func (c *Core) GetTagByID(ctx context.Context, id int64) (*domain.Tag, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id: %d", id)
	}

	Tag := models.Tag{}
	Item := domain.Tag{}

	err := c.db.WithContext(ctx).First(&Tag, id).Error
	if err != nil {
		return &Item, err
	}

	Item.ID = Tag.ID
	Item.Name = Tag.Name

	return &Item, nil
}

func (c *Core) ListTags(ctx context.Context, page, pageSize int) ([]domain.Tag, error) {
	tags := []models.Tag{}
	items := []domain.Tag{}

	err := c.db.WithContext(ctx).Scopes(Paginate(page, pageSize)).Find(&tags).Error
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

func (c *Core) UpdateTag(ctx context.Context, tag *domain.Tag) error {
	if tag.ID <= 0 {
		return fmt.Errorf("invalid id: %d", tag.ID)
	}

	Tag := models.Tag{
		Name: tag.Name,
	}

	result := c.db.WithContext(ctx).Model(&Tag).Where("id = ?", tag.ID).Updates(Tag)
	if result.Error != nil {
		err := fmt.Errorf("failed to update Tag with: %w", result.Error)
		return err
	}

	return nil
}

func (c *Core) DeleteTag(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid id: %d", id)
	}

	Tag := models.Tag{}
	err := c.db.WithContext(ctx).First(&Tag, id).Error
	if err != nil {
		return fmt.Errorf("failed to find Tag with id: %d", id)
	}

	result := c.db.WithContext(ctx).Delete(&Tag, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete Tag with id: %d", id)
	}

	return nil
}
