package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListHotComments 查询热点评论列表。
func (d *UserOpsDAO) ListHotComments(ctx context.Context, limit int) ([]models.WComment, error) {
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Model(&models.WComment{}).Where("status = 1 AND post_id > 0").Order("likes DESC, id DESC")
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 定义并初始化当前变量。
	items := make([]models.WComment, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&items).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return items, nil
}

// CreateHotComment 新增热点评论。
func (d *UserOpsDAO) CreateHotComment(ctx context.Context, item *models.WComment) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}

// UpdateHotComment 更新热点评论。
func (d *UserOpsDAO) UpdateHotComment(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(&models.WComment{}).Where("id = ? AND post_id > 0", id).Updates(updates).Error
}

// DeleteHotComment 删除热点评论。
func (d *UserOpsDAO) DeleteHotComment(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Where("post_id > 0").Delete(&models.WComment{}, id).Error
}
