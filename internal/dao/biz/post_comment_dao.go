package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListPostComments 查询指定帖子下的评论。
func (d *UserOpsDAO) ListPostComments(ctx context.Context, postID uint, limit int) ([]models.WComment, error) {
	query := d.db.WithContext(ctx).Model(&models.WComment{}).Where("post_id = ?", postID).Order("parent_id ASC, id ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	rows := make([]models.WComment, 0)
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// CreatePostComment 新增评论。
func (d *UserOpsDAO) CreatePostComment(ctx context.Context, item *models.WComment) error {
	return d.db.WithContext(ctx).Create(item).Error
}
