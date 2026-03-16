package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
)

// ListHomePopups 查询首屏弹窗列表。
func (d *BizConfigDAO) ListHomePopups(ctx context.Context, position string, limit int) ([]models.WHomePopup, error) {
	position = strings.TrimSpace(position)
	if position == "" {
		position = "home"
	}
	query := d.db.WithContext(ctx).Model(&models.WHomePopup{}).Where("position = ?", position).Order("sort ASC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	items := make([]models.WHomePopup, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateHomePopup 新增首屏弹窗。
func (d *BizConfigDAO) CreateHomePopup(ctx context.Context, item *models.WHomePopup) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateHomePopup 更新首屏弹窗。
func (d *BizConfigDAO) UpdateHomePopup(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WHomePopup{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteHomePopup 删除首屏弹窗。
func (d *BizConfigDAO) DeleteHomePopup(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WHomePopup{}, id).Error
}
