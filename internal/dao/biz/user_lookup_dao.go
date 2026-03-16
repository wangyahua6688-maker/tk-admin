package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
)

// GetUserType 获取指定用户的类型（仅查询启用用户）。
func (d *BizConfigDAO) GetUserType(ctx context.Context, userID uint) (string, error) {
	var user models.WUser
	if err := d.db.WithContext(ctx).Select("id,user_type,status").Where("id = ? AND status = 1", userID).First(&user).Error; err != nil {
		return "", err
	}
	return strings.TrimSpace(user.UserType), nil
}

// IsPostExists 判断帖子是否存在。
func (d *BizConfigDAO) IsPostExists(ctx context.Context, postID uint) (bool, error) {
	var count int64
	if err := d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("id = ?", postID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
