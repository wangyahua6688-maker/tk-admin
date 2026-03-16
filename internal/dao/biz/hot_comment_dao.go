package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListHotComments 查询热点评论列表。
func (d *UserOpsDAO) ListHotComments(ctx context.Context, limit int) ([]models.WComment, error) {
	query := d.db.WithContext(ctx).Model(&models.WComment{}).Where("status = 1 AND post_id > 0").Order("likes DESC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	items := make([]models.WComment, 0)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateHotComment 新增热点评论。
func (d *UserOpsDAO) CreateHotComment(ctx context.Context, item *models.WComment) error {
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateHotComment 更新热点评论。
func (d *UserOpsDAO) UpdateHotComment(ctx context.Context, id uint, updates map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&models.WComment{}).Where("id = ? AND post_id > 0", id).Updates(updates).Error
}

// DeleteHotComment 删除热点评论。
func (d *UserOpsDAO) DeleteHotComment(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Where("post_id > 0").Delete(&models.WComment{}, id).Error
}
