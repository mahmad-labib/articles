package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/service/validation"
	"gorm.io/gorm"
)

type TagRepository interface {
	CreateTag(ctx context.Context, tag *domain.Tag) (int64, error)
	GetTagByID(ctx context.Context, id int64) (*domain.Tag, error)
	ListTags(ctx context.Context, page, pageSize int) ([]domain.Tag, error)
	UpdateTag(ctx context.Context, tag *domain.Tag) error
	DeleteTag(ctx context.Context, id int64) error
}

type TagValidator interface {
	Validate(tag *domain.Tag) error
}

type TagService struct {
	repo      TagRepository
	validator TagValidator
}

func NewTagService(repo TagRepository, validator TagValidator) *TagService {
	return &TagService{
		repo:      repo,
		validator: validator,
	}
}

func (s *TagService) CreateTag(ctx context.Context, tag *domain.Tag) (int64, error) {
	err := s.validator.Validate(tag)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateTag(ctx, tag)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, fmt.Errorf("%v %v", validation.Validation, gorm.ErrDuplicatedKey)
		}

		return 0, err
	}

	return id, nil
}

func (s *TagService) GetTagByID(ctx context.Context, id int64) (*domain.Tag, error) {
	item, err := s.repo.GetTagByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *TagService) ListTags(ctx context.Context, page, pageSize int) ([]domain.Tag, error) {
	items, err := s.repo.ListTags(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *TagService) UpdateTag(ctx context.Context, tag *domain.Tag) error {
	err := s.validator.Validate(tag)
	if err != nil {
		return err
	}

	err = s.repo.UpdateTag(ctx, tag)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("%v %v", validation.Validation, gorm.ErrDuplicatedKey)
		}

		return err
	}

	return nil
}

func (s *TagService) DeleteTag(ctx context.Context, id int64) error {
	err := s.repo.DeleteTag(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
