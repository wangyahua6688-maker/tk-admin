package dao

import (
	"context"
	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (d *UserDao) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) Create(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *UserDao) ListAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := d.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
