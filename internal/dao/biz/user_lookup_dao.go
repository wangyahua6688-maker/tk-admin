package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
)

// GetUserType 获取指定用户的类型（仅查询启用用户）。
func (d *BizConfigDAO) GetUserType(ctx context.Context, userID uint) (string, error) {
	// 声明当前变量。
	var user models.WUser
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Select("id,user_type,status").Where("id = ? AND status = 1", userID).First(&user).Error; err != nil {
		// 返回当前处理结果。
		return "", err
	}
	// 返回当前处理结果。
	return strings.TrimSpace(user.UserType), nil
}

// IsPostExists 判断帖子是否存在。
func (d *BizConfigDAO) IsPostExists(ctx context.Context, postID uint) (bool, error) {
	// 声明当前变量。
	var count int64
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("id = ?", postID).Count(&count).Error; err != nil {
		// 返回当前处理结果。
		return false, err
	}
	// 返回当前处理结果。
	return count > 0, nil
}
