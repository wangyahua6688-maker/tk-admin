package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
)

// GetUserByID 查询单个用户（可过滤状态）。
func (d *UserOpsDAO) GetUserByID(ctx context.Context, userID uint) (*models.WUser, error) {
	var user models.WUser
	if err := d.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetActiveUserType 获取启用用户的类型。
func (d *UserOpsDAO) GetActiveUserType(ctx context.Context, userID uint) (string, error) {
	var user models.WUser
	if err := d.db.WithContext(ctx).Select("id,user_type,status").Where("id = ? AND status = 1", userID).First(&user).Error; err != nil {
		return "", err
	}
	return strings.TrimSpace(user.UserType), nil
}

// GetUsersByIDs 批量获取用户基础信息。
func (d *UserOpsDAO) GetUsersByIDs(ctx context.Context, userIDs []uint) ([]models.WUser, error) {
	if len(userIDs) == 0 {
		return []models.WUser{}, nil
	}
	users := make([]models.WUser, 0)
	if err := d.db.WithContext(ctx).Select("id,username,nickname,user_type").Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
