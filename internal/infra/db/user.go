package db

import (
	"context"
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
	"github.com/mahmad-labib/articles/internal/infra/db/models"
	"gorm.io/gorm"
)

func (r *Repository) CreateUser(ctx context.Context, domUser *domain.User) (int64, error) {
	User := models.User{
		Name:         domUser.Name,
		Email:        domUser.Email,
		PasswordHash: domUser.Password,
	}

	result := r.db.WithContext(ctx).Create(&User)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new User with err: %w", result.Error)

		return 0, err
	}

	return User.ID, nil
}

func (r *Repository) GetUserByID(ctx context.Context, ID int64) *domain.ApiUser {
	User := models.User{}
	Item := domain.ApiUser{}

	err := r.db.First(&User, ID).Error
	if err != nil {
		return nil
	}

	Item.Email = User.Email
	Item.ID = User.ID
	Item.Name = User.Name

	return &Item
}

func (r *Repository) GetUserByEmail(ctx context.Context, Email string) (*domain.ApiUser, error) {
	User := models.User{}
	Item := domain.ApiUser{}

	err := r.db.First(&User, &models.User{Email: Email}).Error
	if err != nil {
		return nil, err
	}

	Item.Email = User.Email
	Item.ID = User.ID
	Item.Name = User.Name

	return &Item, nil
}

func (r *Repository) GetUserByEmail_Login(ctx context.Context, Email string) (string, error) {
	User := models.User{}

	err := r.db.First(&User, &models.User{Email: Email}).Error
	if err != nil {
		return "", err
	}

	return User.PasswordHash, nil
}

func (r *Repository) UpdateUser(ctx context.Context, domUser *domain.ApiUser) error {
	User := models.User{
		Name:  domUser.Name,
		Email: domUser.Email,
	}

	res := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", domUser.ID).Updates(User)

	if res.Error != nil {
		return fmt.Errorf("error when update User with id=%d: %w", domUser.ID, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("User: id=%d: does not exist", domUser.ID)
	}

	return nil
}

func (r *Repository) ListUser(ctx context.Context, Page, PageSize int) ([]domain.ApiUser, error) {
	var Users []models.User

	err := r.db.Scopes(Paginate(Page, PageSize)).Find(&Users).Error
	if err != nil {
		return nil, err
	}

	list := make([]domain.ApiUser, 0, len(Users))
	for _, t := range Users {
		list = append(list, domain.ApiUser{
			ID:    t.ID,
			Name:  t.Name,
			Email: t.Email,
		})
	}

	return list, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	res := r.db.Delete(&models.User{}, id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) GetProfile(ctx context.Context, id int64) (*domain.Profile, error) {
	var user models.User

	if id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	err := r.db.
		Preload("Articles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "user_id")
		}).
		Find(&user, &models.User{ID: id}).Error
	if err != nil {
		return nil, err
	}

	var profile = domain.Profile{
		User: domain.User{
			Name:  user.Name,
			Email: user.Email,
		},
		Articles: []domain.Article{},
	}

	for _, val := range user.Articles {
		article := domain.Article{
			Title: val.Title,
			Slug:  val.Slug,
		}
		profile.Articles = append(profile.Articles, article)
	}

	return &profile, nil
}
