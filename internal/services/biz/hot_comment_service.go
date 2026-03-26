package biz

import (
	"context"
	"errors"
	"strings"

	"go-admin/internal/models"
)

// ListHotComments 查询热点评论列表。
func (s *UserOpsService) ListHotComments(ctx context.Context, limit int) ([]models.WComment, error) {
	// 返回当前处理结果。
	return s.dao.ListHotComments(ctx, limit)
}

// CreateHotComment 新增热点评论。
func (s *UserOpsService) CreateHotComment(ctx context.Context, postID uint, userID uint, parentID *uint, content string, likes *int64, status *int8) (*models.WComment, error) {
	// 必填参数校验。
	if postID == 0 || userID == 0 {
		// 返回当前处理结果。
		return nil, errors.New("post_id/user_id required")
	}
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(content) == "" {
		// 返回当前处理结果。
		return nil, errors.New("content required")
	}

	// 组装评论模型。
	item := models.WComment{
		// 处理当前语句逻辑。
		PostID: postID,
		// 处理当前语句逻辑。
		UserID: userID,
		// 处理当前语句逻辑。
		ParentID: 0,
		// 调用strings.TrimSpace完成当前处理。
		Content: strings.TrimSpace(content),
		// 处理当前语句逻辑。
		Likes: 0,
		// 处理当前语句逻辑。
		Status: 1,
	}
	// 判断条件并进入对应分支逻辑。
	if parentID != nil {
		// 更新当前变量或字段值。
		item.ParentID = *parentID
	}
	// 判断条件并进入对应分支逻辑。
	if likes != nil {
		// 更新当前变量或字段值。
		item.Likes = *likes
	}
	// 判断条件并进入对应分支逻辑。
	if status != nil {
		// 更新当前变量或字段值。
		item.Status = *status
	}

	// 写库。
	if err := s.dao.CreateHotComment(ctx, &item); err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &item, nil
}

// UpdateHotComment 更新热点评论。
func (s *UserOpsService) UpdateHotComment(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdateHotComment(ctx, id, updates)
}

// DeleteHotComment 删除热点评论。
func (s *UserOpsService) DeleteHotComment(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeleteHotComment(ctx, id)
}
