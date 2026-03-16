package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
)

// ListClientUsers 查询客户端用户列表。
func (d *UserOpsDAO) ListClientUsers(ctx context.Context, userType string, limit int) ([]models.WUser, error) {
	query := d.db.WithContext(ctx).Model(&models.WUser{}).Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	// 按用户类型筛选。
	userType = strings.TrimSpace(userType)
	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}
	items := make([]models.WUser, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateClientUser 新增客户端用户。
func (d *UserOpsDAO) CreateClientUser(ctx context.Context, item *models.WUser) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateClientUser 更新客户端用户。
func (d *UserOpsDAO) UpdateClientUser(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WUser{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteClientUser 删除客户端用户。
func (d *UserOpsDAO) DeleteClientUser(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WUser{}, id).Error
}
