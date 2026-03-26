package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
)

// ListHomePopups 查询首屏弹窗列表。
func (s *BizConfigService) ListHomePopups(ctx context.Context, position string, limit int) ([]models.WHomePopup, error) {
	// 更新当前变量或字段值。
	position = strings.TrimSpace(position)
	// 判断条件并进入对应分支逻辑。
	if position == "" {
		// 更新当前变量或字段值。
		position = "home"
	}
	// 返回当前处理结果。
	return s.dao.ListHomePopups(ctx, position, limit)
}

// CreateHomePopup 新增首屏弹窗。
func (s *BizConfigService) CreateHomePopup(ctx context.Context, item *models.WHomePopup) error {
	// 直通 DAO 创建记录。
	return s.dao.CreateHomePopup(ctx, item)
}

// UpdateHomePopup 更新首屏弹窗。
func (s *BizConfigService) UpdateHomePopup(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 直通 DAO 更新记录。
	return s.dao.UpdateHomePopup(ctx, id, updates)
}

// DeleteHomePopup 删除首屏弹窗。
func (s *BizConfigService) DeleteHomePopup(ctx context.Context, id uint) error {
	// 直通 DAO 删除记录。
	return s.dao.DeleteHomePopup(ctx, id)
}
