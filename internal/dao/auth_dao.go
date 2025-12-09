package dao

import (
	"context"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
)

type AuthDao struct {
	db *gorm.DB
}

func NewAuthDao(db *gorm.DB) *AuthDao {
	return &AuthDao{db: db}
}

// GetUserByUsername 根据用户名查询用户
func (d *AuthDao) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	db := d.db
	if ctx != nil {
		// 优先使用上下文中的数据库连接
		if ctxDB := utils.DBFromContext(ctx); ctxDB != nil {
			db = ctxDB
		}
	}

	if err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户（注册用）
func (d *AuthDao) CreateUser(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

// UpdateUserToken 更新用户token相关信息（可选，根据业务需求）
func (d *AuthDao) UpdateUserToken(ctx context.Context, userID uint, refreshToken string) error {
	return d.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", refreshToken).Error
}

// GetUserByID 根据用户ID查询用户
func (d *AuthDao) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	if err := d.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
