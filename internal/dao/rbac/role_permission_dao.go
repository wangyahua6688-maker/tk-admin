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
	return &RolePermissionDao{db: db}
}

// FindRole 查询角色是否存在。
func (d *RolePermissionDao) FindRole(ctx context.Context, id uint) (*models.Role, error) {
	var r models.Role
	if err := d.db.WithContext(ctx).First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// FindPermissions 查询权限集合。
func (d *RolePermissionDao) FindPermissions(ctx context.Context, ids []uint) ([]models.Permission, error) {
	var perms []models.Permission
	err := d.db.WithContext(ctx).Find(&perms, ids).Error
	return perms, err
}

// ReplacePermissions 全量替换角色绑定权限。
func (d *RolePermissionDao) ReplacePermissions(ctx context.Context, role *models.Role, perms []models.Permission) error {
	return d.db.WithContext(ctx).Model(role).Association("Permissions").Replace(&perms)
}

// GetPermissions 查询角色已绑定权限。
func (d *RolePermissionDao) GetPermissions(ctx context.Context, id uint) ([]models.Permission, error) {
	var role models.Role
	err := d.db.WithContext(ctx).
		Preload("Permissions").
		First(&role, id).Error

	if err != nil {
		return nil, err
	}

	return role.Permissions, nil
}
