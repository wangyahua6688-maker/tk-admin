package biz

import (
	"context"

	"go-admin/internal/models"
)

// ListPostArticles 查询帖子列表（按官方/网友区分）。
func (s *UserOpsService) ListPostArticles(ctx context.Context, isOfficial bool, limit int) ([]models.WPostArticle, error) {
	// 定义并初始化当前变量。
	flag := int8(0)
	// 判断条件并进入对应分支逻辑。
	if isOfficial {
		// 更新当前变量或字段值。
		flag = 1
	}
	// 返回当前处理结果。
	return s.dao.ListPostArticles(ctx, flag, limit)
}

// CreatePostArticle 新增帖子。
func (s *UserOpsService) CreatePostArticle(ctx context.Context, item *models.WPostArticle) error {
	// 返回当前处理结果。
	return s.dao.CreatePostArticle(ctx, item)
}

// UpdatePostArticle 更新帖子。
func (s *UserOpsService) UpdatePostArticle(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 返回当前处理结果。
	return s.dao.UpdatePostArticle(ctx, id, updates)
}

// DeletePostArticle 删除帖子。
func (s *UserOpsService) DeletePostArticle(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return s.dao.DeletePostArticle(ctx, id)
}

// IsUserTypes 判断用户是否为指定类型。
func (s *UserOpsService) IsUserTypes(ctx context.Context, userID uint, expectedTypes ...string) bool {
	// 判断条件并进入对应分支逻辑。
	if userID == 0 {
		// 返回当前处理结果。
		return false
	}
	// 定义并初始化当前变量。
	current, err := s.dao.GetActiveUserType(ctx, userID)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return false
	}
	// 循环处理当前数据集合。
	for _, t := range expectedTypes {
		// 判断条件并进入对应分支逻辑。
		if current == t {
			// 返回当前处理结果。
			return true
		}
	}
	// 返回当前处理结果。
	return false
}
