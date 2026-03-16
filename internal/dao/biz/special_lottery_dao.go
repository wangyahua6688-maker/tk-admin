package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListSpecialLotteries 查询彩种配置列表。
func (d *BizConfigDAO) ListSpecialLotteries(ctx context.Context, limit int) ([]models.WSpecialLottery, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	var items []models.WSpecialLottery
	err := d.db.WithContext(ctx).
		Order("sort ASC, id DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// CreateSpecialLottery 新增彩种配置。
func (d *BizConfigDAO) CreateSpecialLottery(ctx context.Context, item *models.WSpecialLottery) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateSpecialLottery 更新彩种配置。
func (d *BizConfigDAO) UpdateSpecialLottery(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WSpecialLottery{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteSpecialLottery 删除彩种配置。
func (d *BizConfigDAO) DeleteSpecialLottery(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WSpecialLottery{}, id).Error
}
