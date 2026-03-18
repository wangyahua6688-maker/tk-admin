package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListOfficialPosts 查询官方发帖列表。
func (d *BizConfigDAO) ListOfficialPosts(ctx context.Context, limit int) ([]models.WPostArticle, error) {
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("is_official = 1").Order("id DESC")
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 定义并初始化当前变量。
	items := make([]models.WPostArticle, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// CreateOfficialPost 新增官方发帖。
func (d *BizConfigDAO) CreateOfficialPost(ctx context.Context, item *models.WPostArticle) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateOfficialPost 更新官方发帖。
func (d *BizConfigDAO) UpdateOfficialPost(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteOfficialPost 删除官方发帖。
func (d *BizConfigDAO) DeleteOfficialPost(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.WPostArticle{}, id).Error
}
