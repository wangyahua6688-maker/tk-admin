package biz

import (
	"context"
	"strings"
)

// IsUserTypes 判断用户是否属于指定类型（任一命中即返回 true）。
func (s *BizConfigService) IsUserTypes(ctx context.Context, userID uint, expectedTypes ...string) bool {
	// 未传用户直接判定失败。
	if userID == 0 {
		// 返回当前处理结果。
		return false
	}
	// 读取用户类型。
	current, err := s.dao.GetUserType(ctx, userID)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return false
	}
	// 更新当前变量或字段值。
	current = strings.TrimSpace(current)
	// 循环处理当前数据集合。
	for _, t := range expectedTypes {
		// 判断条件并进入对应分支逻辑。
		if current == strings.TrimSpace(t) {
			// 返回当前处理结果。
			return true
		}
	}
	// 返回当前处理结果。
	return false
}

// IsPostExists 校验帖子是否存在。
func (s *BizConfigService) IsPostExists(ctx context.Context, postID uint) (bool, error) {
	// 判断条件并进入对应分支逻辑。
	if postID == 0 {
		// 返回当前处理结果。
		return false, nil
	}
	// 返回当前处理结果。
	return s.dao.IsPostExists(ctx, postID)
}
