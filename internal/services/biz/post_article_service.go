package biz

import (
	"context"

	"go-admin-full/internal/models"
)

// ListPostArticles 查询帖子列表（按官方/网友区分）。
func (s *UserOpsService) ListPostArticles(ctx context.Context, isOfficial bool, limit int) ([]models.WPostArticle, error) {
	flag := int8(0)
	if isOfficial {
		flag = 1
	}
	return s.dao.ListPostArticles(ctx, flag, limit)
}

// CreatePostArticle 新增帖子。
func (s *UserOpsService) CreatePostArticle(ctx context.Context, item *models.WPostArticle) error {
	return s.dao.CreatePostArticle(ctx, item)
}

// UpdatePostArticle 更新帖子。
func (s *UserOpsService) UpdatePostArticle(ctx context.Context, id uint, updates map[string]interface{}) error {
	return s.dao.UpdatePostArticle(ctx, id, updates)
}

// DeletePostArticle 删除帖子。
func (s *UserOpsService) DeletePostArticle(ctx context.Context, id uint) error {
	return s.dao.DeletePostArticle(ctx, id)
}

// IsUserTypes 判断用户是否为指定类型。
func (s *UserOpsService) IsUserTypes(ctx context.Context, userID uint, expectedTypes ...string) bool {
	if userID == 0 {
		return false
	}
	current, err := s.dao.GetActiveUserType(ctx, userID)
	if err != nil {
		return false
	}
	for _, t := range expectedTypes {
		if current == t {
			return true
		}
	}
	return false
}
