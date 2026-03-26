package rbac

import (
	"context"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	"strings"
)

// SystemMessageService 系统消息业务层。
// 职责：
// 1. 面向当前用户的消息查询与已读管理；
// 2. 在 RBAC 变更后向受影响用户投递系统消息。
type SystemMessageService struct {
	dao *rbacdao.SystemMessageDao // 系统消息 DAO
}

// NewSystemMessageService 创建系统消息服务。
func NewSystemMessageService(d *rbacdao.SystemMessageDao) *SystemMessageService {
	// 返回当前处理结果。
	return &SystemMessageService{dao: d}
}

// ListUserMessages 查询用户消息列表与未读数。
func (s *SystemMessageService) ListUserMessages(
	// 处理当前语句逻辑。
	ctx context.Context,
	// 处理当前语句逻辑。
	userID uint,
	// 处理当前语句逻辑。
	page int,
	// 处理当前语句逻辑。
	pageSize int,
	// 处理当前语句逻辑。
	onlyUnread bool,
	// 进入新的代码块进行处理。
) ([]models.SystemMessage, int64, int64, error) {
	// 定义并初始化当前变量。
	items, total, err := s.dao.ListByUser(ctx, userID, page, pageSize, onlyUnread)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, 0, 0, err
	}

	// 定义并初始化当前变量。
	unread, err := s.dao.CountUnread(ctx, userID)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, 0, 0, err
	}

	// 返回当前处理结果。
	return items, total, unread, nil
}

// MarkRead 标记单条消息已读。
func (s *SystemMessageService) MarkRead(ctx context.Context, userID uint, messageID uint) error {
	// 返回当前处理结果。
	return s.dao.MarkReadByID(ctx, userID, messageID)
}

// MarkAllRead 标记全部消息已读。
func (s *SystemMessageService) MarkAllRead(ctx context.Context, userID uint) error {
	// 返回当前处理结果。
	return s.dao.MarkAllReadByUser(ctx, userID)
}

// PushToUser 给单个用户投递系统消息。
func (s *SystemMessageService) PushToUser(
	// 处理当前语句逻辑。
	ctx context.Context,
	// 处理当前语句逻辑。
	userID uint,
	// 处理当前语句逻辑。
	title string,
	// 处理当前语句逻辑。
	content string,
	// 处理当前语句逻辑。
	level string,
	// 处理当前语句逻辑。
	bizType string,
	// 处理当前语句逻辑。
	bizID uint,
	// 处理当前语句逻辑。
	operatorID uint,
	// 进入新的代码块进行处理。
) error {
	// 返回当前处理结果。
	return s.PushToUsers(ctx, []uint{userID}, title, content, level, bizType, bizID, operatorID)
}

// PushToUsers 给多个用户批量投递系统消息。
func (s *SystemMessageService) PushToUsers(
	// 处理当前语句逻辑。
	ctx context.Context,
	// 处理当前语句逻辑。
	userIDs []uint,
	// 处理当前语句逻辑。
	title string,
	// 处理当前语句逻辑。
	content string,
	// 处理当前语句逻辑。
	level string,
	// 处理当前语句逻辑。
	bizType string,
	// 处理当前语句逻辑。
	bizID uint,
	// 处理当前语句逻辑。
	operatorID uint,
	// 进入新的代码块进行处理。
) error {
	ids := uniquePositiveIDs(userIDs) // 去重并过滤无效 ID
	// 判断条件并进入对应分支逻辑。
	if len(ids) == 0 {
		// 返回当前处理结果。
		return nil
	}

	title = strings.TrimSpace(title)     // 清理标题
	content = strings.TrimSpace(content) // 清理内容
	// 判断条件并进入对应分支逻辑。
	if title == "" || content == "" {
		// 返回当前处理结果。
		return nil
	}

	msgLevel := normalizeMessageLevel(level) // 规范化消息等级
	// 定义并初始化当前变量。
	messages := make([]models.SystemMessage, 0, len(ids))
	// 循环处理当前数据集合。
	for _, uid := range ids {
		// 更新当前变量或字段值。
		messages = append(messages, models.SystemMessage{
			// 处理当前语句逻辑。
			UserID: uid,
			// 处理当前语句逻辑。
			Title: title,
			// 处理当前语句逻辑。
			Content: content,
			// 处理当前语句逻辑。
			Level: msgLevel,
			// 处理当前语句逻辑。
			IsRead: false,
			// 处理当前语句逻辑。
			OperatorID: operatorID,
			// 调用strings.TrimSpace完成当前处理。
			BizType: strings.TrimSpace(bizType),
			// 处理当前语句逻辑。
			BizID: bizID,
		})
	}

	// 返回当前处理结果。
	return s.dao.CreateBatch(ctx, messages)
}

// ListUserIDsByRoleIDs 获取角色关联的用户ID集合。
func (s *SystemMessageService) ListUserIDsByRoleIDs(ctx context.Context, roleIDs []uint) ([]uint, error) {
	// 返回当前处理结果。
	return s.dao.FindUserIDsByRoleIDs(ctx, uniquePositiveIDs(roleIDs))
}

// ListUserIDsByPermissionIDs 获取权限关联角色下的用户ID集合。
func (s *SystemMessageService) ListUserIDsByPermissionIDs(ctx context.Context, permissionIDs []uint) ([]uint, error) {
	roleIDs, err := s.dao.FindRoleIDsByPermissionIDs(ctx, uniquePositiveIDs(permissionIDs)) // 先查角色
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return s.ListUserIDsByRoleIDs(ctx, roleIDs)
}

// PushToUsersByRoleIDs 向角色关联用户推送消息。
func (s *SystemMessageService) PushToUsersByRoleIDs(
	// 处理当前语句逻辑。
	ctx context.Context,
	// 处理当前语句逻辑。
	roleIDs []uint,
	// 处理当前语句逻辑。
	title string,
	// 处理当前语句逻辑。
	content string,
	// 处理当前语句逻辑。
	level string,
	// 处理当前语句逻辑。
	bizType string,
	// 处理当前语句逻辑。
	bizID uint,
	// 处理当前语句逻辑。
	operatorID uint,
	// 进入新的代码块进行处理。
) error {
	// 定义并初始化当前变量。
	userIDs, err := s.ListUserIDsByRoleIDs(ctx, roleIDs)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return err
	}
	// 返回当前处理结果。
	return s.PushToUsers(ctx, userIDs, title, content, level, bizType, bizID, operatorID)
}

// PushToUsersByPermissionIDs 向权限关联角色下的用户推送消息。
func (s *SystemMessageService) PushToUsersByPermissionIDs(
	// 处理当前语句逻辑。
	ctx context.Context,
	// 处理当前语句逻辑。
	permissionIDs []uint,
	// 处理当前语句逻辑。
	title string,
	// 处理当前语句逻辑。
	content string,
	// 处理当前语句逻辑。
	level string,
	// 处理当前语句逻辑。
	bizType string,
	// 处理当前语句逻辑。
	bizID uint,
	// 处理当前语句逻辑。
	operatorID uint,
	// 进入新的代码块进行处理。
) error {
	// 定义并初始化当前变量。
	userIDs, err := s.ListUserIDsByPermissionIDs(ctx, permissionIDs)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return err
	}
	// 返回当前处理结果。
	return s.PushToUsers(ctx, userIDs, title, content, level, bizType, bizID, operatorID)
}

// uniquePositiveIDs 处理uniquePositiveIDs相关逻辑。
func uniquePositiveIDs(in []uint) []uint {
	// 使用 map 去重
	set := make(map[uint]struct{}, len(in))
	// 定义并初始化当前变量。
	out := make([]uint, 0, len(in))
	// 循环处理当前数据集合。
	for _, v := range in {
		// 判断条件并进入对应分支逻辑。
		if v == 0 {
			// 处理当前语句逻辑。
			continue
		}
		// 判断条件并进入对应分支逻辑。
		if _, ok := set[v]; ok {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		set[v] = struct{}{}
		// 更新当前变量或字段值。
		out = append(out, v)
	}
	// 返回当前处理结果。
	return out
}

// normalizeMessageLevel 处理normalizeMessageLevel相关逻辑。
func normalizeMessageLevel(level string) string {
	// 根据表达式进入多分支处理。
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "success":
		return "success"
	case "warning":
		return "warning"
	default:
		return "info"
	}
}
