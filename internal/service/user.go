package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"main.go/internal/domain"
)

var (
	empty          = errors.New("field is required")
	nameLength     = errors.New("name length min is 5 and max is 30")
	emailLength    = errors.New("email length min is 5 and max is 30")
	emailUnique    = errors.New("this email already exists")
	passwordSpaces = errors.New("password cant conatin spaces")
	emailSpaces    = errors.New("email cant conatin spaces")
)

func (a *ArticleService) GetUserByEmail(ctx context.Context, Email string) (*domain.ApiUser, error) {
	user, err := a.repo.GetUserByEmail(ctx, Email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (a *ArticleService) validateUser(ctx context.Context, user *domain.User) error {
	pass := utf8.RuneCountInString(user.Password)
	if pass == 0 {
		return fmt.Errorf("password %v", empty)
	}

	for _, r := range user.Password {
		if unicode.IsSpace(r) {
			return passwordSpaces
		}
	}

	name := utf8.RuneCountInString(user.Name)
	if name == 0 {
		return fmt.Errorf("name %v", empty)
	}
	if name < 2 || name > 30 {
		return nameLength
	}

	for _, r := range user.Email {
		if unicode.IsSpace(r) {
			return emailSpaces
		}
	}

	email := utf8.RuneCountInString(user.Email)
	if email == 0 {
		return fmt.Errorf("email %v", empty)
	}
	if email > 100 {
		return emailLength
	}

	userWithEmail, err := a.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if userWithEmail != nil && userWithEmail.Email != "" {
		return emailUnique
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	return nil
}

func (a *ArticleService) CreateUser(ctx context.Context, domUser *domain.User) (int64, error) {
	err := a.validateUser(ctx, domUser)
	if err != nil {
		return 0, err
	}

	id, err := a.repo.CreateUser(ctx, domUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *ArticleService) GetUserByID(ctx context.Context, ID int64) (*domain.ApiUser, error) {
	user, err := a.repo.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return user, nil
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
