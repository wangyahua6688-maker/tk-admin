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
	return &SystemMessageService{dao: d}
}

// ListUserMessages 查询用户消息列表与未读数。
func (s *SystemMessageService) ListUserMessages(
	ctx context.Context,
	userID uint,
	page int,
	pageSize int,
	onlyUnread bool,
) ([]models.SystemMessage, int64, int64, error) {
	items, total, err := s.dao.ListByUser(ctx, userID, page, pageSize, onlyUnread)
	if err != nil {
		return nil, 0, 0, err
	}

	unread, err := s.dao.CountUnread(ctx, userID)
	if err != nil {
		return nil, 0, 0, err
	}

	return items, total, unread, nil
}

// MarkRead 标记单条消息已读。
func (s *SystemMessageService) MarkRead(ctx context.Context, userID uint, messageID uint) error {
	return s.dao.MarkReadByID(ctx, userID, messageID)
}

// MarkAllRead 标记全部消息已读。
func (s *SystemMessageService) MarkAllRead(ctx context.Context, userID uint) error {
	return s.dao.MarkAllReadByUser(ctx, userID)
}

// PushToUser 给单个用户投递系统消息。
func (s *SystemMessageService) PushToUser(
	ctx context.Context,
	userID uint,
	title string,
	content string,
	level string,
	bizType string,
	bizID uint,
	operatorID uint,
) error {
	return s.PushToUsers(ctx, []uint{userID}, title, content, level, bizType, bizID, operatorID)
}

// PushToUsers 给多个用户批量投递系统消息。
func (s *SystemMessageService) PushToUsers(
	ctx context.Context,
	userIDs []uint,
	title string,
	content string,
	level string,
	bizType string,
	bizID uint,
	operatorID uint,
) error {
	ids := uniquePositiveIDs(userIDs) // 去重并过滤无效 ID
	if len(ids) == 0 {
		return nil
	}

	title = strings.TrimSpace(title)     // 清理标题
	content = strings.TrimSpace(content) // 清理内容
	if title == "" || content == "" {
		return nil
	}

	msgLevel := normalizeMessageLevel(level) // 规范化消息等级
	messages := make([]models.SystemMessage, 0, len(ids))
	for _, uid := range ids {
		messages = append(messages, models.SystemMessage{
			UserID:     uid,
			Title:      title,
			Content:    content,
			Level:      msgLevel,
			IsRead:     false,
			OperatorID: operatorID,
			BizType:    strings.TrimSpace(bizType),
			BizID:      bizID,
		})
	}

	return s.dao.CreateBatch(ctx, messages)
}

// ListUserIDsByRoleIDs 获取角色关联的用户ID集合。
func (s *SystemMessageService) ListUserIDsByRoleIDs(ctx context.Context, roleIDs []uint) ([]uint, error) {
	return s.dao.FindUserIDsByRoleIDs(ctx, uniquePositiveIDs(roleIDs))
}

// ListUserIDsByPermissionIDs 获取权限关联角色下的用户ID集合。
func (s *SystemMessageService) ListUserIDsByPermissionIDs(ctx context.Context, permissionIDs []uint) ([]uint, error) {
	roleIDs, err := s.dao.FindRoleIDsByPermissionIDs(ctx, uniquePositiveIDs(permissionIDs)) // 先查角色
	if err != nil {
		return nil, err
	}
	return s.ListUserIDsByRoleIDs(ctx, roleIDs)
}

// PushToUsersByRoleIDs 向角色关联用户推送消息。
func (s *SystemMessageService) PushToUsersByRoleIDs(
	ctx context.Context,
	roleIDs []uint,
	title string,
	content string,
	level string,
	bizType string,
	bizID uint,
	operatorID uint,
) error {
	userIDs, err := s.ListUserIDsByRoleIDs(ctx, roleIDs)
	if err != nil {
		return err
	}
	return s.PushToUsers(ctx, userIDs, title, content, level, bizType, bizID, operatorID)
}

// PushToUsersByPermissionIDs 向权限关联角色下的用户推送消息。
func (s *SystemMessageService) PushToUsersByPermissionIDs(
	ctx context.Context,
	permissionIDs []uint,
	title string,
	content string,
	level string,
	bizType string,
	bizID uint,
	operatorID uint,
) error {
	userIDs, err := s.ListUserIDsByPermissionIDs(ctx, permissionIDs)
	if err != nil {
		return err
	}
	return s.PushToUsers(ctx, userIDs, title, content, level, bizType, bizID, operatorID)
}

func uniquePositiveIDs(in []uint) []uint {
	// 使用 map 去重
	set := make(map[uint]struct{}, len(in))
	out := make([]uint, 0, len(in))
	for _, v := range in {
		if v == 0 {
			continue
		}
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func normalizeMessageLevel(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "success":
		return "success"
	case "warning":
		return "warning"
	default:
		return "info"
	}
}
