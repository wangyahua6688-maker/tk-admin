package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListSMSChannels 查询短信通道配置列表。
func (d *BizConfigDAO) ListSMSChannels(ctx context.Context, status *int, limit int) ([]models.WSMSChannel, error) {
	query := d.db.WithContext(ctx).Model(&models.WSMSChannel{}).Order("status DESC, id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	// 仅当传入 0/1 时筛选状态。
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	items := make([]models.WSMSChannel, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateSMSChannel 新增短信通道配置。
func (d *BizConfigDAO) CreateSMSChannel(ctx context.Context, item *models.WSMSChannel) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateSMSChannel 更新短信通道配置。
func (d *BizConfigDAO) UpdateSMSChannel(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WSMSChannel{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteSMSChannel 删除短信通道配置。
func (d *BizConfigDAO) DeleteSMSChannel(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WSMSChannel{}, id).Error
}
