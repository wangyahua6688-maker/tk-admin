package rbac

import (
	"context"
	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// UserRoleDao 负责用户与角色的关联关系。
type UserRoleDao struct {
	db *gorm.DB // 数据库连接
}

// NewUserRoleDao 创建用户角色 DAO。
func NewUserRoleDao(db *gorm.DB) *UserRoleDao {
	return &UserRoleDao{db: db}
}

// FindUser 查询用户是否存在。
func (d *UserRoleDao) FindUser(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindRoles 查询角色集合。
func (d *UserRoleDao) FindRoles(ctx context.Context, roleIDs []uint) ([]models.Role, error) {
	var roles []models.Role
	if err := d.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// ReplaceRoles 全量替换用户角色。
func (d *UserRoleDao) ReplaceRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	return d.db.WithContext(ctx).Model(user).Association("Roles").Replace(&roles)
}

// AppendRoles 追加用户角色。
func (d *UserRoleDao) AppendRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	return d.db.WithContext(ctx).Model(user).Association("Roles").Append(&roles)
}

// RemoveRoles 移除用户角色。
func (d *UserRoleDao) RemoveRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	return d.db.WithContext(ctx).Model(user).Association("Roles").Delete(&roles)
}

// GetUserRolesWithPermissions 查询用户角色并预加载权限。
func (d *UserRoleDao) GetUserRolesWithPermissions(ctx context.Context, id uint) ([]models.Role, error) {
	var user models.User
	err := d.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&user, id).Error

	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}
