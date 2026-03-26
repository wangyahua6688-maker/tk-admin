package biz

import (
	"context"
	"strings"

	"go-admin/internal/models"
)

// ListClientUsers 查询客户端用户列表。
func (d *UserOpsDAO) ListClientUsers(ctx context.Context, userType string, limit int) ([]models.WUser, error) {
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Model(&models.WUser{}).Order("id DESC")
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 按用户类型筛选。
	userType = strings.TrimSpace(userType)
	// 判断条件并进入对应分支逻辑。
	if userType != "" {
		// 更新当前变量或字段值。
		query = query.Where("user_type = ?", userType)
	}
	// 定义并初始化当前变量。
	items := make([]models.WUser, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// CreateClientUser 新增客户端用户。
func (d *UserOpsDAO) CreateClientUser(ctx context.Context, item *models.WUser) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateClientUser 更新客户端用户。
func (d *UserOpsDAO) UpdateClientUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WUser{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteClientUser 删除客户端用户。
func (d *UserOpsDAO) DeleteClientUser(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WUser{}, id).Error
}
