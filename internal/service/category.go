package service

import (
	"context"
	"errors"

	"github.com/mahmad-labib/articles/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category) (int64, error)
	ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error)
	GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

type CategoryValidator interface {
	Validate(category *domain.Category) error
}

type CategoryService struct {
	repo      CategoryRepository
	validator CategoryValidator
}

func NewCategoryService(repo CategoryRepository, validator CategoryValidator) *CategoryService {
	return &CategoryService{
		repo:      repo,
		validator: validator,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (int64, error) {
	err := s.validator.Validate(category)
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

func (s *CategoryService) ListCategories(ctx context.Context, page, pageSize int) ([]domain.Category, error) {
	list, err := s.repo.ListCategories(ctx, page, pageSize)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int64) (*domain.Category, error) {
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	err := s.validator.Validate(category)
	if err != nil {
		return err
	}

	err = s.repo.UpdateCategory(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int64) error {
	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
