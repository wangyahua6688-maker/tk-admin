package biz

import (
	"context"

	"go-admin/internal/models"
)

// ListBroadcasts 查询系统广播列表。
func (s *BizConfigService) ListBroadcasts(ctx context.Context, limit int) ([]models.WBroadcast, error) {
	// 返回当前处理结果。
	return s.dao.ListBroadcasts(ctx, limit)
}

// CreateBroadcast 新增系统广播。
func (s *BizConfigService) CreateBroadcast(ctx context.Context, item *models.WBroadcast) error {
	// 返回当前处理结果。
	return s.dao.CreateBroadcast(ctx, item)
}

// UpdateBroadcast 更新系统广播。
func (s *BizConfigService) UpdateBroadcast(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateBroadcast(ctx, id, updates)
}

// DeleteBroadcast 删除系统广播。
func (s *BizConfigService) DeleteBroadcast(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteBroadcast(ctx, id)
}
