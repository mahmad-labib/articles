package service

import (
	"context"

	"main.go/internal/domain"
)

func (s *UserService) CreateArticle(ctx context.Context, article domain.Article, userId int64, catId int64, tags []domain.Tag) (*domain.Article, error) {
	Article, err := s.repo.CreateArticle(ctx, article, userId, catId, tags)
	if err != nil {
		return nil, err
	}

	return Article, nil
}

func (s *UserService) GetArticleByID(ctx context.Context, id int64) (*domain.Article, error) {
	Article, err := s.repo.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return Article, nil
}
