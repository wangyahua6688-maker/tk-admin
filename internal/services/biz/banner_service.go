package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListBanners 查询 Banner 列表。
func (s *BizConfigService) ListBanners(ctx context.Context, bannerType string, limit int) ([]models.WBanner, error) {
	return s.dao.ListBanners(ctx, bannerType, limit)
}

// CreateBanner 新增 Banner。
func (s *BizConfigService) CreateBanner(ctx context.Context, item *models.WBanner) error {
	return s.dao.CreateBanner(ctx, item)
}

// UpdateBanner 更新 Banner。
func (s *BizConfigService) UpdateBanner(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.dao.UpdateBanner(ctx, id, updates)
}

// DeleteBanner 删除 Banner。
func (s *BizConfigService) DeleteBanner(ctx context.Context, id uint) error {
	return s.dao.DeleteBanner(ctx, id)
}
