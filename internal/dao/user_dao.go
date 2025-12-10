package dao

import (
	"context"
	"fmt"

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

// GetByID 根据用户ID查询用户。
func (d *UserDao) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDao) ListAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := d.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateByID 更新用户可编辑字段。
func (d *UserDao) UpdateByID(ctx context.Context, id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}
	return d.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID 删除用户。
func (d *UserDao) DeleteByID(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
