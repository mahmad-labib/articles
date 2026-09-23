package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/domain"
	"main.go/internal/infra/db/models"
)

func (C *Core) CreateUser(ctx context.Context, domUser *domain.User) (int64, error) {
	User := models.User{
		Name:         domUser.Name,
		Email:        domUser.Email,
		PasswordHash: domUser.PasswordHash,
	}

	result := C.db.WithContext(ctx).Create(&User)
	if result.Error != nil {
		err := fmt.Errorf("failed to create new User with err: %w", result.Error)
		return 0, err
	}

	return User.ID, nil
}

func (C *Core) GetUserByID(ctx context.Context, ID int64) (*domain.User, error) {
	User := models.User{}
	Item := domain.User{}

	err := C.db.First(&User, ID).Error
	if err != nil {
		return &Item, err
	}

	Item.Email = User.Email
	Item.ID = User.ID
	Item.Name = User.Name

	return &Item, nil
}

func (C *Core) UpdateUser(ctx context.Context, domUser *domain.User) error {
	User := models.User{
		Name:  domUser.Name,
		Email: domUser.Email,
	}

	res := C.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", domUser.ID).Updates(User)

	if res.Error != nil {
		return fmt.Errorf("error when update User with id=%d: %w", domUser.ID, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("User: id=%d: does not exist", domUser.ID)
	}

	return nil
}

func (C *Core) ListUser(ctx context.Context, Page, PageSize int) ([]domain.User, error) {
	var Users []models.User

	err := C.db.Scopes(Paginate(Page, PageSize)).Find(&Users).Error
	if err != nil {
		return nil, err
	}

	list := make([]domain.User, 0, len(Users))
	for _, t := range Users {
		list = append(list, domain.User{
			ID:    t.ID,
			Name:  t.Name,
			Email: t.Email,
		})
	}

	return list, nil
}

func Paginate(page, PageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case PageSize > 100:
			PageSize = 100
		case PageSize <= 0:
			PageSize = 10
		}

		offset := (page - 1) * PageSize
		return db.Offset(offset).Limit(PageSize)
	}
}

func (C *Core) DeleteUser(ctx context.Context, id int) error {
	res := C.db.Delete(&models.User{}, id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
