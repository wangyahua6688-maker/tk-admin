package biz

import (
	"context"
	"strings"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// BannerDAO Banner 数据访问层。
type BannerDAO struct {
	db *gorm.DB
}

// NewBannerDAO 创建 Banner DAO。
func NewBannerDAO(db *gorm.DB) *BannerDAO {
	return &BannerDAO{db: db}
}

// ListBanners 查询 Banner 列表。
func (d *BannerDAO) ListBanners(ctx context.Context, bannerType string, limit int) ([]models.WBanner, error) {
	// 判断条件并进入对应分支逻辑。
	if limit <= 0 || limit > 1000 {
		// 更新当前变量或字段值。
		limit = 300
	}
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Order("type ASC, sort ASC, id DESC").Limit(limit)
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(bannerType) != "" {
		// 更新当前变量或字段值。
		query = query.Where("type = ?", strings.TrimSpace(bannerType))
	}

	// 声明当前变量。
	var items []models.WBanner
	// 定义并初始化当前变量。
	err := query.Find(&items).Error
	// 返回当前处理结果。
	return items, err
}

// CreateBanner 新增 Banner。
func (d *BannerDAO) CreateBanner(ctx context.Context, item *models.WBanner) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateBanner 更新 Banner。
func (d *BannerDAO) UpdateBanner(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WBanner{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteBanner 删除 Banner。
func (d *BannerDAO) DeleteBanner(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WBanner{}, id).Error
}
