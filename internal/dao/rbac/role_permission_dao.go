package rbac

import (
	"context"
	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// RolePermissionDao 负责角色与权限的关联关系。
type RolePermissionDao struct {
	db *gorm.DB // 数据库连接
}

// NewRolePermissionDao 创建角色权限 DAO。
func NewRolePermissionDao(db *gorm.DB) *RolePermissionDao {
	// 返回当前处理结果。
	return &RolePermissionDao{db: db}
}

// FindRole 查询角色是否存在。
func (d *RolePermissionDao) FindRole(ctx context.Context, id uint) (*models.Role, error) {
	// 声明当前变量。
	var r models.Role
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&r, id).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &r, nil
}

// FindPermissions 查询权限集合。
func (d *RolePermissionDao) FindPermissions(ctx context.Context, ids []uint) ([]models.Permission, error) {
	// 声明当前变量。
	var perms []models.Permission
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).Find(&perms, ids).Error
	// 返回当前处理结果。
	return perms, err
}

// ReplacePermissions 全量替换角色绑定权限。
func (d *RolePermissionDao) ReplacePermissions(ctx context.Context, role *models.Role, perms []models.Permission) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Model(role).Association("Permissions").Replace(&perms)
}

// GetPermissions 查询角色已绑定权限。
func (d *RolePermissionDao) GetPermissions(ctx context.Context, id uint) ([]models.Permission, error) {
	// 声明当前变量。
	var role models.Role
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).
		// 调用Preload完成当前处理。
		Preload("Permissions").
		// 调用First完成当前处理。
		First(&role, id).Error

	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, err
	}

	// 返回当前处理结果。
	return role.Permissions, nil
}
