package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/service/validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, ID int64) *domain.ApiUser
	UpdateUser(ctx context.Context, domUser *domain.ApiUser) error
	ListUser(ctx context.Context, Page, PageSize int) ([]domain.ApiUser, error)
	DeleteUser(ctx context.Context, ID int64) error
	GetUserByEmail(ctx context.Context, Email string) (*domain.ApiUser, error)
	GetUserByEmail_Login(ctx context.Context, Email string) (string, error)
	GetProfile(ctx context.Context, id int64) (*domain.Profile, error)
	CreateUser(ctx context.Context, domUser *domain.User) (int64, error)
}

type UserValidator interface {
	Validate(item *domain.User) error
	ValidateEmail(email string) error
}

type UserService struct {
	repo      UserRepository
	validator UserValidator
}

func NewUserService(repo UserRepository, validator UserValidator) *UserService {
	return &UserService{
		repo:      repo,
		validator: validator,
	}
}

var (
	PasswordInvalid = errors.New("password is not correct")
)

func (s *UserService) GetUserByEmail(ctx context.Context, Email string) (*domain.ApiUser, error) {
	user, err := s.repo.GetUserByEmail(ctx, Email)
	if err != nil {
		return nil, fmt.Errorf("%v %v", validation.NotFoundErr, gorm.ErrRecordNotFound)
	}

	return user, nil
}

func (s *UserService) hashPass(user *domain.User) (*domain.User, error) {
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

func (s *UserService) CreateUser(ctx context.Context, domUser *domain.User) (int64, error) {
	err := s.validator.Validate(domUser)
	if err != nil {
		return 0, err
	}

	domUser, err = s.hashPass(domUser)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateUser(ctx, domUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, fmt.Errorf("%v this email already exist", validation.Validation)
		}

		return 0, err
	}

	return id, nil
}

func (s *UserService) GetUserByID(ctx context.Context, ID int64) *domain.ApiUser {
	user := s.repo.GetUserByID(ctx, ID)

	return user
}

func (s *UserService) ListUser(ctx context.Context, Page, PageSize int) ([]domain.ApiUser, error) {
	usersList, err := s.repo.ListUser(ctx, Page, PageSize)
	if err != nil {
		return usersList, err
	}

	return usersList, nil
}

func (s *UserService) DeleteUser(ctx context.Context, Id int64) error {
	err := s.repo.DeleteUser(ctx, Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string, jwt *jwtauth.JWTAuth) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%v %v", validation.NotFoundErr, gorm.ErrRecordNotFound)
	}

	userPass, err := s.repo.GetUserByEmail_Login(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%v %v", validation.NotFoundErr, gorm.ErrRecordNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userPass), []byte(password))
	if err != nil {
		return "", fmt.Errorf("%v %v", validation.Validation, PasswordInvalid)
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

func (s *UserService) GetProfile(ctx context.Context, id int64) (*domain.Profile, error) {
	if id == 0 {
		return nil, fmt.Errorf("%v %v", validation.Validation, gorm.ErrInvalidValue)
	}

	user, err := s.repo.GetProfile(ctx, id)
	if err != nil {
		return user, err
	}

	return user, nil
}
