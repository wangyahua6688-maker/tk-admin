package rbac

import (
	"context"
	"go-admin-full/internal/models"

	"gorm.io/gorm"
)

// SystemMessageDao 系统消息数据访问层。
type SystemMessageDao struct {
	db *gorm.DB // 数据库连接
}

// NewSystemMessageDao 创建系统消息 DAO。
func NewSystemMessageDao(db *gorm.DB) *SystemMessageDao {
	return &SystemMessageDao{db: db}
}

// CreateBatch 批量写入系统消息。
func (d *SystemMessageDao) CreateBatch(ctx context.Context, messages []models.SystemMessage) error {
	if len(messages) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).Create(&messages).Error
}

// ListByUser 按用户分页查询系统消息。
func (d *SystemMessageDao) ListByUser(
	ctx context.Context,
	userID uint,
	page int,
	pageSize int,
	onlyUnread bool,
) ([]models.SystemMessage, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := d.db.WithContext(ctx).Model(&models.SystemMessage{}).Where("user_id = ?", userID) // 构造查询
	if onlyUnread {
		query = query.Where("is_read = ?", false) // 过滤未读
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.SystemMessage
	err := query.
		// 未读优先展示，其次按创建时间倒序
		Order("is_read ASC").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// CountUnread 统计用户未读消息数。
func (d *SystemMessageDao) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var total int64
	err := d.db.WithContext(ctx).
		Model(&models.SystemMessage{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&total).Error
	return total, err
}

// MarkReadByID 标记单条消息为已读（校验用户归属）。
func (d *SystemMessageDao) MarkReadByID(ctx context.Context, userID uint, messageID uint) error {
	var msg models.SystemMessage
	// 先校验消息归属，避免越权更新
	if err := d.db.WithContext(ctx).
		Select("id", "user_id").
		Where("id = ? AND user_id = ?", messageID, userID).
		First(&msg).Error; err != nil {
		return err
	}

	return d.db.WithContext(ctx).
		Model(&models.SystemMessage{}).
		Where("id = ? AND user_id = ?", messageID, userID).
		Update("is_read", true).Error
}

// MarkAllReadByUser 标记用户全部消息为已读。
func (d *SystemMessageDao) MarkAllReadByUser(ctx context.Context, userID uint) error {
	return d.db.WithContext(ctx).
		Model(&models.SystemMessage{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// FindUserIDsByRoleIDs 根据角色ID集合查询关联用户ID集合。
func (d *SystemMessageDao) FindUserIDsByRoleIDs(ctx context.Context, roleIDs []uint) ([]uint, error) {
	if len(roleIDs) == 0 {
		return []uint{}, nil
	}

	var userIDs []uint
	err := d.db.WithContext(ctx).
		// 通过用户-角色关系表反查用户列表
		Table("sys_user_roles").
		Distinct("user_id").
		Where("role_id IN ?", roleIDs).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// FindRoleIDsByPermissionIDs 根据权限ID集合查询关联角色ID集合。
func (d *SystemMessageDao) FindRoleIDsByPermissionIDs(ctx context.Context, permissionIDs []uint) ([]uint, error) {
	if len(permissionIDs) == 0 {
		return []uint{}, nil
	}

	var roleIDs []uint
	err := d.db.WithContext(ctx).
		// 通过角色-权限关系表反查角色列表
		Table("sys_role_permissions").
		Distinct("role_id").
		Where("permission_id IN ?", permissionIDs).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}
