package biz

import (
	"context"
	"strings"
)

// IsUserTypes 判断用户是否属于指定类型（任一命中即返回 true）。
func (s *BizConfigService) IsUserTypes(ctx context.Context, userID uint, expectedTypes ...string) bool {
	// 未传用户直接判定失败。
	if userID == 0 {
		return false
	}
	// 读取用户类型。
	current, err := s.dao.GetUserType(ctx, userID)
	if err != nil {
		return false
	}
	current = strings.TrimSpace(current)
	for _, t := range expectedTypes {
		if current == strings.TrimSpace(t) {
			return true
		}
	}
	return false
}

// IsPostExists 校验帖子是否存在。
func (s *BizConfigService) IsPostExists(ctx context.Context, postID uint) (bool, error) {
	if postID == 0 {
		return false, nil
	}
	return s.dao.IsPostExists(ctx, postID)
}
