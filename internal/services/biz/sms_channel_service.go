package biz

import (
	"context"
	"strconv"
	"strings"

	"go-admin-full/internal/models"
)

// ListSMSChannels 查询短信通道配置列表。
func (s *BizConfigService) ListSMSChannels(ctx context.Context, status string, limit int) ([]models.WSMSChannel, error) {
	// 将字符串状态转换为可选筛选值。
	var statusFilter *int
	trimmed := strings.TrimSpace(status)
	if trimmed == "0" || trimmed == "1" {
		v, _ := strconv.Atoi(trimmed)
		statusFilter = &v
	}
	// 调用 DAO 执行查询。
	return s.dao.ListSMSChannels(ctx, statusFilter, limit)
}

// CreateSMSChannel 新增短信通道配置。
func (s *BizConfigService) CreateSMSChannel(ctx context.Context, item *models.WSMSChannel) error {
	// 直通 DAO 创建记录。
	return s.dao.CreateSMSChannel(ctx, item)
}

// UpdateSMSChannel 更新短信通道配置。
func (s *BizConfigService) UpdateSMSChannel(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 直通 DAO 更新记录。
	return s.dao.UpdateSMSChannel(ctx, id, updates)
}

// DeleteSMSChannel 删除短信通道配置。
func (s *BizConfigService) DeleteSMSChannel(ctx context.Context, id uint) error {
	// 直通 DAO 删除记录。
	return s.dao.DeleteSMSChannel(ctx, id)
}
