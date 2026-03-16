package biz

import (
	"context"
	"errors"
	"strings"

	"go-admin-full/internal/models"
)

// ListPostComments 查询帖子评论并补充用户信息。
func (s *UserOpsService) ListPostComments(ctx context.Context, postID uint) ([]map[string]interface{}, error) {
	// 拉取评论列表。
	rows, err := s.dao.ListPostComments(ctx, postID, 500)
	if err != nil {
		return nil, err
	}

	// 收集用户ID。
	userIDs := make([]uint, 0)
	seen := make(map[uint]struct{})
	for _, row := range rows {
		if _, ok := seen[row.UserID]; ok {
			continue
		}
		seen[row.UserID] = struct{}{}
		userIDs = append(userIDs, row.UserID)
	}

	// 批量查询用户信息。
	userMap := make(map[uint]models.WUser)
	if len(userIDs) > 0 {
		users, uerr := s.dao.GetUsersByIDs(ctx, userIDs)
		if uerr != nil {
			return nil, uerr
		}
		for _, u := range users {
			userMap[u.ID] = u
		}
	}

	// 组装返回结构。
	items := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		u := userMap[row.UserID]
		items = append(items, map[string]interface{}{
			"id":         row.ID,
			"post_id":    row.PostID,
			"user_id":    row.UserID,
			"parent_id":  row.ParentID,
			"content":    row.Content,
			"likes":      row.Likes,
			"status":     row.Status,
			"created_at": row.CreatedAt,
			"username":   u.Username,
			"nickname":   u.Nickname,
			"user_type":  u.UserType,
		})
	}

	return items, nil
}

// CreatePostComment 新增帖子评论（机器人/官方）。
func (s *UserOpsService) CreatePostComment(ctx context.Context, postID uint, userID uint, parentID uint, content string, status *int8) (*models.WComment, error) {
	// 参数校验。
	if postID == 0 {
		return nil, errors.New("post_id required")
	}
	if userID == 0 {
		return nil, errors.New("user_id required")
	}
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("content required")
	}

	// 校验用户类型。
	if !s.IsUserTypes(ctx, userID, "robot", "official") {
		return nil, errors.New("user must be robot or official account")
	}

	// 校验帖子存在。
	if _, err := s.dao.GetPostArticleByID(ctx, postID); err != nil {
		return nil, errors.New("post not found")
	}

	// 组装评论模型。
	item := models.WComment{
		PostID:   postID,
		UserID:   userID,
		ParentID: parentID,
		Content:  strings.TrimSpace(content),
		Likes:    0,
		Status:   1,
	}
	if status != nil {
		item.Status = *status
	}

	// 落库。
	if err := s.dao.CreatePostComment(ctx, &item); err != nil {
		return nil, err
	}
	return &item, nil
}
