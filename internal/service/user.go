package service

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"main.go/internal/domain"
	"main.go/internal/service/validation"
)

func (a *ArticleService) GetUserByEmail(ctx context.Context, Email string) (*domain.ApiUser, error) {
	user, err := a.repo.GetUserByEmail(ctx, Email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (a *ArticleService) hashPass(user *domain.User) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return user, err
	}
	user.Password = string(hash)

	return user, nil
}

func (a *ArticleService) CreateUser(ctx context.Context, domUser *domain.User) (int64, error) {
	err := a.validator.ValidateUser(domUser)
	if err != nil {
		return 0, err
	}

	domUser, err = a.hashPass(domUser)
	if err != nil {
		return 0, err
	}

	userWithEmail, err := a.GetUserByEmail(ctx, domUser.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if userWithEmail != nil && userWithEmail.Email != "" {
		return 0, validation.EmailUnique
	}

	id, err := a.repo.CreateUser(ctx, domUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *ArticleService) GetUserByID(ctx context.Context, ID int64) *domain.ApiUser {
	user := a.repo.GetUserByID(ctx, ID)

	return user
}

func (a *ArticleService) ListUser(ctx context.Context, Page, PageSize int) ([]domain.ApiUser, error) {
	usersList, err := a.repo.ListUser(ctx, Page, PageSize)
	if err != nil {
		return usersList, err
	}

	return usersList, nil
}

func (a *ArticleService) DeleteUser(ctx context.Context, Id int64) error {
	err := a.repo.DeleteUser(ctx, Id)
	if err != nil {
		return err
	}

	return nil
}
