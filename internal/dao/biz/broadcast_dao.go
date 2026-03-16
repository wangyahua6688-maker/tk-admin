package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListBroadcasts 查询系统广播列表。
func (d *BizConfigDAO) ListBroadcasts(ctx context.Context, limit int) ([]models.WBroadcast, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	var items []models.WBroadcast
	err := d.db.WithContext(ctx).
		Order("sort ASC, id DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// CreateBroadcast 新增系统广播。
func (d *BizConfigDAO) CreateBroadcast(ctx context.Context, item *models.WBroadcast) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateBroadcast 更新系统广播。
func (d *BizConfigDAO) UpdateBroadcast(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WBroadcast{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteBroadcast 删除系统广播。
func (d *BizConfigDAO) DeleteBroadcast(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WBroadcast{}, id).Error
}
