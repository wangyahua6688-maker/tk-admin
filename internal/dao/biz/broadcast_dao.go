package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListBroadcasts 查询系统广播列表。
func (d *BizConfigDAO) ListBroadcasts(ctx context.Context, limit int) ([]models.WBroadcast, error) {
	// 判断条件并进入对应分支逻辑。
	if limit <= 0 || limit > 1000 {
		// 更新当前变量或字段值。
		limit = 200
	}
	// 声明当前变量。
	var items []models.WBroadcast
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).
		// 调用Order完成当前处理。
		Order("sort ASC, id DESC").
		// 调用Limit完成当前处理。
		Limit(limit).
		// 调用Find完成当前处理。
		Find(&items).Error
	// 返回当前处理结果。
	return items, err
}

// CreateBroadcast 新增系统广播。
func (d *BizConfigDAO) CreateBroadcast(ctx context.Context, item *models.WBroadcast) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateBroadcast 更新系统广播。
func (d *BizConfigDAO) UpdateBroadcast(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WBroadcast{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteBroadcast 删除系统广播。
func (d *BizConfigDAO) DeleteBroadcast(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WBroadcast{}, id).Error
}
