package biz

import (
	"context"
	"strings"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// ListLotteryCategories 查询图库分类列表。
func (d *LotteryDAO) ListLotteryCategories(ctx context.Context, keyword string) ([]models.WLotteryCategory, error) {
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Order("sort ASC, id ASC")
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(keyword) != "" {
		// 定义并初始化当前变量。
		like := "%" + strings.TrimSpace(keyword) + "%"
		// 更新当前变量或字段值。
		query = query.Where("category_key LIKE ? OR name LIKE ? OR search_keywords LIKE ?", like, like, like)
	}

	// 声明当前变量。
	var items []models.WLotteryCategory
	// 定义并初始化当前变量。
	err := query.Find(&items).Error
	// 返回当前处理结果。
	return items, err
}

// CreateLotteryCategory 新增图库分类。
func (d *LotteryDAO) CreateLotteryCategory(ctx context.Context, item *models.WLotteryCategory) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateLotteryCategory 更新图库分类。
func (d *LotteryDAO) UpdateLotteryCategory(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WLotteryCategory{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteLotteryCategory 删除图库分类。
func (d *LotteryDAO) DeleteLotteryCategory(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WLotteryCategory{}, id).Error
}

// GetLotteryCategoryByID 按主键获取分类。
func (d *LotteryDAO) GetLotteryCategoryByID(ctx context.Context, id uint) (*models.WLotteryCategory, error) {
	// 声明当前变量。
	var cat models.WLotteryCategory
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&cat, id).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &cat, nil
}

// GetLotteryCategoryByTag 按 key/name 查询分类（兼容旧请求）。
func (d *LotteryDAO) GetLotteryCategoryByTag(ctx context.Context, tag string) (*models.WLotteryCategory, error) {
	// 更新当前变量或字段值。
	tag = strings.TrimSpace(tag)
	// 判断条件并进入对应分支逻辑。
	if tag == "" {
		// 返回当前处理结果。
		return nil, gorm.ErrRecordNotFound
	}
	// 声明当前变量。
	var cat models.WLotteryCategory
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).
		// 更新当前变量或字段值。
		Where("category_key = ? OR name = ?", tag, tag).
		// 调用Order完成当前处理。
		Order("id ASC").
		// 调用First完成当前处理。
		First(&cat).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &cat, nil
}
