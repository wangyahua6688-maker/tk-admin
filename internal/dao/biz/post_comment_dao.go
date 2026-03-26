package biz

import (
	"context"

	"go-admin/internal/models"
)

// ListPostComments 查询指定帖子下的评论。
func (d *UserOpsDAO) ListPostComments(ctx context.Context, postID uint, limit int) ([]models.WComment, error) {
	// 定义并初始化当前变量。
	query := d.db.WithContext(ctx).Model(&models.WComment{}).Where("post_id = ?", postID).Order("parent_id ASC, id ASC")
	// 判断条件并进入对应分支逻辑。
	if limit > 0 {
		// 更新当前变量或字段值。
		query = query.Limit(limit)
	}
	// 定义并初始化当前变量。
	rows := make([]models.WComment, 0)
	// 判断条件并进入对应分支逻辑。
	if err := query.Find(&rows).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return rows, nil
}

// CreatePostComment 新增评论。
func (d *UserOpsDAO) CreatePostComment(ctx context.Context, item *models.WComment) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(item).Error
}
