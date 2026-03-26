package rbac

import (
	"context"
	"go-admin/internal/models"
	"gorm.io/gorm"
)

// UserRoleDao 负责用户与角色的关联关系。
type UserRoleDao struct {
	db *gorm.DB // 数据库连接
}

// NewUserRoleDao 创建用户角色 DAO。
func NewUserRoleDao(db *gorm.DB) *UserRoleDao {
	// 返回当前处理结果。
	return &UserRoleDao{db: db}
}

// FindUser 查询用户是否存在。
func (d *UserRoleDao) FindUser(ctx context.Context, id uint) (*models.User, error) {
	// 声明当前变量。
	var user models.User
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &user, nil
}

// FindRoles 查询角色集合。
func (d *UserRoleDao) FindRoles(ctx context.Context, roleIDs []uint) ([]models.Role, error) {
	// 声明当前变量。
	var roles []models.Role
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return roles, nil
}

// ReplaceRoles 全量替换用户角色。
func (d *UserRoleDao) ReplaceRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(user).Association("Roles").Replace(&roles)
}

// AppendRoles 追加用户角色。
func (d *UserRoleDao) AppendRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(user).Association("Roles").Append(&roles)
}

// RemoveRoles 移除用户角色。
func (d *UserRoleDao) RemoveRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(user).Association("Roles").Delete(&roles)
}

// GetUserIDsByRoleID 查询拥有指定角色的所有用户 ID。
// 用于角色权限变更后批量清除 Redis 权限缓存，保证下次请求时重新加载最新权限。
func (d *UserRoleDao) GetUserIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	var userIDs []uint
	// 直接查关联表，只取 user_id，避免加载完整用户对象造成不必要的内存开销
	err := d.db.WithContext(ctx).
		Table("sys_user_roles").
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

// GetUserRolesWithPermissions 查询用户角色并预加载权限。
func (d *UserRoleDao) GetUserRolesWithPermissions(ctx context.Context, id uint) ([]models.Role, error) {
	// 声明当前变量。
	var user models.User
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).
		// 调用Preload完成当前处理。
		Preload("Roles.Permissions").
		// 调用First完成当前处理。
		First(&user, id).Error

	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return user.Roles, nil
}
