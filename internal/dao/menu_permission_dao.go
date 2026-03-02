package dao

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

// MenuPermissionDao 处理菜单与权限的关联关系。
type MenuPermissionDao struct {
	db *gorm.DB
}

func NewMenuPermissionDao(db *gorm.DB) *MenuPermissionDao {
	return &MenuPermissionDao{db: db}
}

// FindMenu 查询菜单是否存在。
func (d *MenuPermissionDao) FindMenu(ctx context.Context, id uint) (*models.Menu, error) {
	var m models.Menu
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// FindPermissions 查询权限集合。
func (d *MenuPermissionDao) FindPermissions(ctx context.Context, ids []uint) ([]models.Permission, error) {
	var perms []models.Permission
	if err := d.db.WithContext(ctx).Find(&perms, ids).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// ReplacePermissions 全量替换菜单绑定权限。
func (d *MenuPermissionDao) ReplacePermissions(ctx context.Context, menu *models.Menu, perms []models.Permission) error {
	return d.db.WithContext(ctx).Model(menu).Association("Permissions").Replace(&perms)
}

// GetPermissions 查询菜单已绑定权限。
func (d *MenuPermissionDao) GetPermissions(ctx context.Context, menuID uint) ([]models.Permission, error) {
	var menu models.Menu
	if err := d.db.WithContext(ctx).Preload("Permissions").First(&menu, menuID).Error; err != nil {
		return nil, err
	}
	return menu.Permissions, nil
}
