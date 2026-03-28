package biz

import (
	"context"

	bizdao "go-admin/internal/dao/biz"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// BannerService Banner 业务服务。
type BannerService struct {
	dao *bizdao.BannerDAO
}

// NewBannerService 创建 Banner 服务。
func NewBannerService(db *gorm.DB) *BannerService {
	return &BannerService{dao: bizdao.NewBannerDAO(db)}
}

// ListBanners 查询 Banner 列表。
func (s *BannerService) ListBanners(ctx context.Context, bannerType string, limit int) ([]models.WBanner, error) {
	// 返回当前处理结果。
	return s.dao.ListBanners(ctx, bannerType, limit)
}

// CreateBanner 新增 Banner。
func (s *BannerService) CreateBanner(ctx context.Context, item *models.WBanner) error {
	// 返回当前处理结果。
	return s.dao.CreateBanner(ctx, item)
}

// UpdateBanner 更新 Banner。
func (s *BannerService) UpdateBanner(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateBanner(ctx, id, updates)
}

// DeleteBanner 删除 Banner。
func (s *BannerService) DeleteBanner(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteBanner(ctx, id)
}
