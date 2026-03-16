package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
)

// ListBanners 查询 Banner 列表。
func (d *BizConfigDAO) ListBanners(ctx context.Context, bannerType string, limit int) ([]models.WBanner, error) {
	if limit <= 0 || limit > 1000 {
		limit = 300
	}
	query := d.db.WithContext(ctx).Order("type ASC, sort ASC, id DESC").Limit(limit)
	if strings.TrimSpace(bannerType) != "" {
		query = query.Where("type = ?", strings.TrimSpace(bannerType))
	}

	var items []models.WBanner
	err := query.Find(&items).Error
	return items, err
}

// CreateBanner 新增 Banner。
func (d *BizConfigDAO) CreateBanner(ctx context.Context, item *models.WBanner) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateBanner 更新 Banner。
func (d *BizConfigDAO) UpdateBanner(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WBanner{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteBanner 删除 Banner。
func (d *BizConfigDAO) DeleteBanner(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WBanner{}, id).Error
}
