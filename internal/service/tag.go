package service

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"main.go/internal/domain"
)

func (s *UserService) CreateTag(ctx context.Context, tag *domain.Tag) (int64, error) {
	err := s.validator.ValidateTag(tag)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateTag(ctx, tag)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, gorm.ErrDuplicatedKey
		}

		return 0, err
	}

	return id, nil
}

func (s *UserService) GetTagByID(ctx context.Context, id int64) (*domain.Tag, error) {
	item, err := s.repo.GetTagByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *UserService) ListTags(ctx context.Context, page, pageSize int) ([]domain.Tag, error) {
	items, err := s.repo.ListTags(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *UserService) UpdateTag(ctx context.Context, tag *domain.Tag) error {
	err := s.validator.ValidateTag(tag)
	if err != nil {
		return err
	}

	err = s.repo.UpdateTag(ctx, tag)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return gorm.ErrDuplicatedKey
		}

		return err
	}

	return nil
}

func (s *UserService) DeleteTag(ctx context.Context, id int64) error {
	err := s.repo.DeleteTag(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
