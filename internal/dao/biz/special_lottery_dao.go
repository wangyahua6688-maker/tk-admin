package biz

import (
	"context"

	"go-admin/internal/models"
)

// ListSpecialLotteries 查询彩种配置列表。
func (d *BizConfigDAO) ListSpecialLotteries(ctx context.Context, limit int) ([]models.WSpecialLottery, error) {
	// 判断条件并进入对应分支逻辑。
	if limit <= 0 || limit > 1000 {
		// 更新当前变量或字段值。
		limit = 200
	}
	// 声明当前变量。
	var items []models.WSpecialLottery
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

// CreateSpecialLottery 新增彩种配置。
func (d *BizConfigDAO) CreateSpecialLottery(ctx context.Context, item *models.WSpecialLottery) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateSpecialLottery 更新彩种配置。
func (d *BizConfigDAO) UpdateSpecialLottery(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WSpecialLottery{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteSpecialLottery 删除彩种配置。
func (d *BizConfigDAO) DeleteSpecialLottery(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WSpecialLottery{}, id).Error
}
