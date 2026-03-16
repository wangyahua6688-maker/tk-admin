package biz

import (
	"context"
	"strings"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// ListLotteryCategories 查询图库分类列表。
func (d *BizConfigDAO) ListLotteryCategories(ctx context.Context, keyword string) ([]models.WLotteryCategory, error) {
	query := d.db.WithContext(ctx).Order("sort ASC, id ASC")
	if strings.TrimSpace(keyword) != "" {
		like := "%" + strings.TrimSpace(keyword) + "%"
		query = query.Where("category_key LIKE ? OR name LIKE ? OR search_keywords LIKE ?", like, like, like)
	}

	var items []models.WLotteryCategory
	err := query.Find(&items).Error
	return items, err
}

// CreateLotteryCategory 新增图库分类。
func (d *BizConfigDAO) CreateLotteryCategory(ctx context.Context, item *models.WLotteryCategory) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateLotteryCategory 更新图库分类。
func (d *BizConfigDAO) UpdateLotteryCategory(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WLotteryCategory{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteLotteryCategory 删除图库分类。
func (d *BizConfigDAO) DeleteLotteryCategory(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WLotteryCategory{}, id).Error
}

// GetLotteryCategoryByID 按主键获取分类。
func (d *BizConfigDAO) GetLotteryCategoryByID(ctx context.Context, id uint) (*models.WLotteryCategory, error) {
	var cat models.WLotteryCategory
	if err := d.db.WithContext(ctx).First(&cat, id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

// GetLotteryCategoryByTag 按 key/name 查询分类（兼容旧请求）。
func (d *BizConfigDAO) GetLotteryCategoryByTag(ctx context.Context, tag string) (*models.WLotteryCategory, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var cat models.WLotteryCategory
	if err := d.db.WithContext(ctx).
		Where("category_key = ? OR name = ?", tag, tag).
		Order("id ASC").
		First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}
