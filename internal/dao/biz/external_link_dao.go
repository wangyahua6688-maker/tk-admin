package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListExternalLinks 查询外链列表。
func (d *BizConfigDAO) ListExternalLinks(ctx context.Context, limit int) ([]models.WExternalLink, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	var items []models.WExternalLink
	err := d.db.WithContext(ctx).
		Order("position ASC, sort ASC, id DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// CreateExternalLink 新增外链。
func (d *BizConfigDAO) CreateExternalLink(ctx context.Context, item *models.WExternalLink) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateExternalLink 更新外链。
func (d *BizConfigDAO) UpdateExternalLink(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WExternalLink{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteExternalLink 删除外链。
func (d *BizConfigDAO) DeleteExternalLink(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WExternalLink{}, id).Error
}
