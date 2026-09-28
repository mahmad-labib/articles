package service

import (
	"context"
	"errors"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"main.go/internal/domain"
	"main.go/internal/service/validation"
)

var (
	passwordInvalid = errors.New("password is not correct")
)

func (a *ArticleService) GetUserByEmail(ctx context.Context, Email string) (*domain.ApiUser, error) {
	user, err := a.repo.GetUserByEmail(ctx, Email)
	if err != nil {
		return nil, gorm.ErrRecordNotFound
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

	userWithEmail, _ := a.repo.GetUserByEmail(ctx, domUser.Email)
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

func (a *ArticleService) Login(ctx context.Context, email, password string, jwt *jwtauth.JWTAuth) (string, error) {
	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", gorm.ErrRecordNotFound
	}

	userPass, err := a.repo.GetUserByEmail_Login(ctx, email)
	if err != nil {
		return "", gorm.ErrRecordNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPass), []byte(password))
	if err != nil {
		return "", passwordInvalid
	}

	_, token, err := jwt.Encode(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})
	if err != nil {
		return token, err
	}

	return token, nil
}

func (a *ArticleService) GetProfile(ctx context.Context, id int64) (*domain.Profile, error) {
	user, err := a.repo.GetProfile(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}
