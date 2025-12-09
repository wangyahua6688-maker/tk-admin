package dao

import (
	"context"
	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type RolePermissionDao struct {
	db *gorm.DB
}

func NewRolePermissionDao(db *gorm.DB) *RolePermissionDao {
	return &RolePermissionDao{db: db}
}

func (d *RolePermissionDao) FindRole(ctx context.Context, id uint) (*models.Role, error) {
	var r models.Role
	if err := d.db.WithContext(ctx).First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (d *RolePermissionDao) FindPermissions(ctx context.Context, ids []uint) ([]models.Permission, error) {
	var perms []models.Permission
	err := d.db.WithContext(ctx).Find(&perms, ids).Error
	return perms, err
}

func (d *RolePermissionDao) ReplacePermissions(ctx context.Context, role *models.Role, perms []models.Permission) error {
	return d.db.WithContext(ctx).Model(role).Association("Permissions").Replace(&perms)
}

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
