package biz

import (
	"context"
	"strings"

	"go-admin/internal/models"
)

// ListHomePopups 查询首屏弹窗列表。
func (d *BizConfigDAO) ListHomePopups(ctx context.Context, position string, limit int) ([]models.WHomePopup, error) {
	// 更新当前变量或字段值。
	position = strings.TrimSpace(position)
	// 判断条件并进入对应分支逻辑。
	if position == "" {
		// 更新当前变量或字段值。
		position = "home"
	}
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Model(&models.WHomePopup{}).Where("position = ?", position).Order("sort ASC, id DESC")
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 定义并初始化当前变量。
	items := make([]models.WHomePopup, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// CreateHomePopup 新增首屏弹窗。
func (d *BizConfigDAO) CreateHomePopup(ctx context.Context, item *models.WHomePopup) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateHomePopup 更新首屏弹窗。
func (d *BizConfigDAO) UpdateHomePopup(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WHomePopup{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteHomePopup 删除首屏弹窗。
func (d *BizConfigDAO) DeleteHomePopup(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WHomePopup{}, id).Error
}
