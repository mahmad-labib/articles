package service

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"main.go/internal/domain"
)

func (s *UserService) CreateCategory(ctx context.Context, category *domain.Category) (int64, error) {
	err := s.validator.ValidateCategory(category)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return id, gorm.ErrDuplicatedKey
		}

		return id, err
	}

	return id, nil
}

func (s *UserService) ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error) {
	list, err := s.repo.ListCategories(ctx, page, pageSize)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (s *UserService) GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *UserService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	err := s.validator.ValidateCategory(category)
	if err != nil {
		return err
	}

	err = s.repo.UpdateCategory(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteCategory(ctx context.Context, id int64) error {
	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
