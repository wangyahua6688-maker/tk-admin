package rbac

import (
	"context"
	"go-admin/internal/models"

	"gorm.io/gorm"
)

// SystemMessageDao 系统消息数据访问层。
type SystemMessageDao struct {
	db *gorm.DB // 数据库连接
}

// NewSystemMessageDao 创建系统消息 DAO。
func NewSystemMessageDao(db *gorm.DB) *SystemMessageDao {
	// 返回当前处理结果。
	return &SystemMessageDao{db: db}
}

// CreateBatch 批量写入系统消息。
func (d *SystemMessageDao) CreateBatch(ctx context.Context, messages []models.SystemMessage) error {
	// 判断条件并进入对应分支逻辑。
	if len(messages) == 0 {
		// 返回当前处理结果。
		return nil
	}
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(&messages).Error
}

// ListByUser 按用户分页查询系统消息。
func (d *SystemMessageDao) ListByUser(
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
) ([]models.SystemMessage, int64, error) {
	// 判断条件并进入对应分支逻辑。
	if page <= 0 {
		// 更新当前变量或字段值。
		page = 1
	}
	// 判断条件并进入对应分支逻辑。
	if pageSize <= 0 {
		// 更新当前变量或字段值。
		pageSize = 20
	}

	query := d.db.WithContext(ctx).Model(&models.SystemMessage{}).Where("user_id = ?", userID) // 构造查询
	// 判断条件并进入对应分支逻辑。
	if onlyUnread {
		query = query.Where("is_read = ?", false) // 过滤未读
	}

	// 声明当前变量。
	var total int64
	// 判断条件并进入对应分支逻辑。
	if err := query.Count(&total).Error; err != nil {
		// 返回当前处理结果。
		return nil, 0, err
	}

	// 声明当前变量。
	var items []models.SystemMessage
	// 定义并初始化当前变量。
	err := query.
		// 未读优先展示，其次按创建时间倒序
		Order("is_read ASC").
		// 调用Order完成当前处理。
		Order("created_at DESC").
		// 调用Offset完成当前处理。
		Offset((page - 1) * pageSize).
		// 调用Limit完成当前处理。
		Limit(pageSize).
		// 调用Find完成当前处理。
		Find(&items).Error
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, 0, err
	}

	// 返回当前处理结果。
	return items, total, nil
}

// CountUnread 统计用户未读消息数。
func (d *SystemMessageDao) CountUnread(ctx context.Context, userID uint) (int64, error) {
	// 声明当前变量。
	var total int64
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).
		// 调用Model完成当前处理。
		Model(&models.SystemMessage{}).
		// 更新当前变量或字段值。
		Where("user_id = ? AND is_read = ?", userID, false).
		// 调用Count完成当前处理。
		Count(&total).Error
	// 返回当前处理结果。
	return total, err
}

// MarkReadByID 标记单条消息为已读（校验用户归属）。
func (d *SystemMessageDao) MarkReadByID(ctx context.Context, userID uint, messageID uint) error {
	// 声明当前变量。
	var msg models.SystemMessage
	// 先校验消息归属，避免越权更新
	if err := d.db.WithContext(ctx).
		// 调用Select完成当前处理。
		Select("id", "user_id").
		// 更新当前变量或字段值。
		Where("id = ? AND user_id = ?", messageID, userID).
		// 调用First完成当前处理。
		First(&msg).Error; err != nil {
		// 返回当前处理结果。
		return err
	}

	// 返回当前处理结果。
	return d.db.WithContext(ctx).
		// 调用Model完成当前处理。
		Model(&models.SystemMessage{}).
		// 更新当前变量或字段值。
		Where("id = ? AND user_id = ?", messageID, userID).
		// 调用Update完成当前处理。
		Update("is_read", true).Error
}

// MarkAllReadByUser 标记用户全部消息为已读。
func (d *SystemMessageDao) MarkAllReadByUser(ctx context.Context, userID uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).
		// 调用Model完成当前处理。
		Model(&models.SystemMessage{}).
		// 更新当前变量或字段值。
		Where("user_id = ? AND is_read = ?", userID, false).
		// 调用Update完成当前处理。
		Update("is_read", true).Error
}

// FindUserIDsByRoleIDs 根据角色ID集合查询关联用户ID集合。
func (d *SystemMessageDao) FindUserIDsByRoleIDs(ctx context.Context, roleIDs []uint) ([]uint, error) {
	// 判断条件并进入对应分支逻辑。
	if len(roleIDs) == 0 {
		// 返回当前处理结果。
		return []uint{}, nil
	}

	// 声明当前变量。
	var userIDs []uint
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).
		// 通过用户-角色关系表反查用户列表
		Table("sys_user_roles").
		// 调用Distinct完成当前处理。
		Distinct("user_id").
		// 调用Where完成当前处理。
		Where("role_id IN ?", roleIDs).
		// 调用Pluck完成当前处理。
		Pluck("user_id", &userIDs).Error
	// 返回当前处理结果。
	return userIDs, err
}

// FindRoleIDsByPermissionIDs 根据权限ID集合查询关联角色ID集合。
func (d *SystemMessageDao) FindRoleIDsByPermissionIDs(ctx context.Context, permissionIDs []uint) ([]uint, error) {
	// 判断条件并进入对应分支逻辑。
	if len(permissionIDs) == 0 {
		// 返回当前处理结果。
		return []uint{}, nil
	}

	// 声明当前变量。
	var roleIDs []uint
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).
		// 通过角色-权限关系表反查角色列表
		Table("sys_role_permissions").
		// 调用Distinct完成当前处理。
		Distinct("role_id").
		// 调用Where完成当前处理。
		Where("permission_id IN ?", permissionIDs).
		// 调用Pluck完成当前处理。
		Pluck("role_id", &roleIDs).Error
	// 返回当前处理结果。
	return roleIDs, err
}
