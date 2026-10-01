package service

import (
	"context"

	"main.go/internal/domain"
)

func (s *UserService) CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) error {
	err := s.repo.CreateArticle(ctx, article, userId, catId, tags)
	if err != nil {
		return err
	}
	return nil
}
