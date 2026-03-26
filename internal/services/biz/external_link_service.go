package biz

import (
	"context"

	"go-admin/internal/models"
)

// ListExternalLinks 查询外链列表。
func (s *BizConfigService) ListExternalLinks(ctx context.Context, limit int) ([]models.WExternalLink, error) {
	// 返回当前处理结果。
	return s.dao.ListExternalLinks(ctx, limit)
}

// CreateExternalLink 新增外链。
func (s *BizConfigService) CreateExternalLink(ctx context.Context, item *models.WExternalLink) error {
	// 返回当前处理结果。
	return s.dao.CreateExternalLink(ctx, item)
}

// UpdateExternalLink 更新外链。
func (s *BizConfigService) UpdateExternalLink(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateExternalLink(ctx, id, updates)
}

// DeleteExternalLink 删除外链。
func (s *BizConfigService) DeleteExternalLink(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteExternalLink(ctx, id)
}
