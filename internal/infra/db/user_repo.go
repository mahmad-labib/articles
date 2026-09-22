package db

import (
	"context"
	"fmt"

	"main.go/internal/domain"
	"main.go/internal/infra/db/models"
)

func (A *ArticleDB) Create(ctx context.Context, domainUser domain.User) (int64, error) {
	User := models.User{
		Name:         domainUser.Name,
		Email:        domainUser.Email,
		PasswordHash: domainUser.PasswordHash,
	}

	result := A.db.WithContext(ctx).Create(&User)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new User with err: %w", result.Error)
		return 0, err
	}

	return User.ID, nil
}
