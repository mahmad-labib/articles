package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"main.go/internal/domain"
)

var (
	emailErr       = errors.New("email")
	empty          = errors.New("field is required")
	nameLength     = errors.New("name length min is 2 and max is 30")
	emailUnique    = errors.New("this email already exists")
	passwordSpaces = errors.New("password cant conatin spaces")
)

func isValidEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	return addr.Address == email
}

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
			return fmt.Errorf("%v email cant conatin spaces", emailErr)
		}
	}

	if res := isValidEmail(user.Email); !res {
		return fmt.Errorf("invalid %v", emailErr)
	}

	email := utf8.RuneCountInString(user.Email)
	if email == 0 {
		return fmt.Errorf("email %v", empty)
	}
	if email > 100 {
		return fmt.Errorf("%v length max is 100", emailErr)
	}

	userWithEmail, err := a.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if userWithEmail != nil && userWithEmail.Email != "" {
		return emailUnique
	}

	return nil
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
	err := a.validateUser(ctx, domUser)
	if err != nil {
		return 0, err
	}

	domUser, err = a.hashPass(domUser)
	if err != nil {
		return 0, err
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
