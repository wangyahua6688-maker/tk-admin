package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListOfficialPosts 查询官方发帖列表。
func (d *BizConfigDAO) ListOfficialPosts(ctx context.Context, limit int) ([]models.WPostArticle, error) {
	query := d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("is_official = 1").Order("id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	items := make([]models.WPostArticle, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateOfficialPost 新增官方发帖。
func (d *BizConfigDAO) CreateOfficialPost(ctx context.Context, item *models.WPostArticle) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateOfficialPost 更新官方发帖。
func (d *BizConfigDAO) UpdateOfficialPost(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WPostArticle{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteOfficialPost 删除官方发帖。
func (d *BizConfigDAO) DeleteOfficialPost(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.WPostArticle{}, id).Error
}
