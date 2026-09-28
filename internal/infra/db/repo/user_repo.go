package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/domain"
	"main.go/internal/infra/db/models"
)

func (c *Core) CreateUser(ctx context.Context, domUser *domain.User) (int64, error) {
	User := models.User{
		Name:         domUser.Name,
		Email:        domUser.Email,
		PasswordHash: domUser.Password,
	}

	result := c.db.WithContext(ctx).Create(&User)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new User with err: %w", result.Error)
		return 0, err
	}

	return User.ID, nil
}

func (c *Core) GetUserByID(ctx context.Context, ID int64) *domain.ApiUser {
	User := models.User{}
	Item := domain.ApiUser{}

	err := c.db.First(&User, ID).Error
	if err != nil {
		return nil
	}

	Item.Email = User.Email
	Item.ID = User.ID
	Item.Name = User.Name

	return &Item
}

func (c *Core) GetUserByEmail(ctx context.Context, Email string) (*domain.ApiUser, error) {
	User := models.User{}
	Item := domain.ApiUser{}

	err := c.db.First(&User, &models.User{Email: Email}).Error
	if err != nil {
		return nil, err
	}

	Item.Email = User.Email
	Item.ID = User.ID
	Item.Name = User.Name

	return &Item, nil
}

func (c *Core) UpdateUser(ctx context.Context, domUser *domain.ApiUser) error {
	User := models.User{
		Name:  domUser.Name,
		Email: domUser.Email,
	}

	res := c.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", domUser.ID).Updates(User)

	if res.Error != nil {
		return fmt.Errorf("error when update User with id=%d: %w", domUser.ID, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("User: id=%d: does not exist", domUser.ID)
	}

	return nil
}

func (c *Core) ListUser(ctx context.Context, Page, PageSize int) ([]domain.ApiUser, error) {
	var Users []models.User

	err := c.db.Scopes(Paginate(Page, PageSize)).Find(&Users).Error
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

func (c *Core) DeleteUser(ctx context.Context, id int64) error {
	res := c.db.Delete(&models.User{}, id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
