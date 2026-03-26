package biz

import (
	"context"
	"strings"

	"go-admin/internal/models"
)

// GetActiveUserType 获取启用用户的类型。
func (d *UserOpsDAO) GetActiveUserType(ctx context.Context, userID uint) (string, error) {
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

// GetUsersByIDs 批量获取用户基础信息。
func (d *UserOpsDAO) GetUsersByIDs(ctx context.Context, userIDs []uint) ([]models.WUser, error) {
	// 判断条件并进入对应分支逻辑。
	if len(userIDs) == 0 {
		// 返回当前处理结果。
		return []models.WUser{}, nil
	}
	// 定义并初始化当前变量。
	users := make([]models.WUser, 0)
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Select("id,username,nickname,user_type").Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return users, nil
}
