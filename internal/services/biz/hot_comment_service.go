package biz

import (
	"context"
	"errors"
	"strings"

	"go-admin-full/internal/models"
)

// ListHotComments 查询热点评论列表。
func (s *UserOpsService) ListHotComments(ctx context.Context, limit int) ([]models.WComment, error) {
	return s.dao.ListHotComments(ctx, limit)
}

// CreateHotComment 新增热点评论。
func (s *UserOpsService) CreateHotComment(ctx context.Context, postID uint, userID uint, parentID *uint, content string, likes *int64, status *int8) (*models.WComment, error) {
	// 必填参数校验。
	if postID == 0 || userID == 0 {
		return nil, errors.New("post_id/user_id required")
	}
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("content required")
	}

	// 组装评论模型。
	item := models.WComment{
		PostID:   postID,
		UserID:   userID,
		ParentID: 0,
		Content:  strings.TrimSpace(content),
		Likes:    0,
		Status:   1,
	}
	if parentID != nil {
		item.ParentID = *parentID
	}
	if likes != nil {
		item.Likes = *likes
	}
	if status != nil {
		item.Status = *status
	}

	// 写库。
	if err := s.dao.CreateHotComment(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateHotComment 更新热点评论。
func (s *UserOpsService) UpdateHotComment(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.dao.UpdateHotComment(ctx, id, updates)
}

// DeleteHotComment 删除热点评论。
func (s *UserOpsService) DeleteHotComment(ctx context.Context, id uint) error {
	return s.dao.DeleteHotComment(ctx, id)
}
