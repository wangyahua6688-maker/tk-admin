package rbac

import (
	"context"
	"errors"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// MenuPermissionDao 处理菜单与权限的关联关系。
type MenuPermissionDao struct {
	db *gorm.DB // 数据库连接
}

// NewMenuPermissionDao 创建菜单权限 DAO。
func NewMenuPermissionDao(db *gorm.DB) *MenuPermissionDao {
	// 返回当前处理结果。
	return &MenuPermissionDao{db: db}
}

// FindMenu 查询菜单是否存在。
func (d *MenuPermissionDao) FindMenu(ctx context.Context, id uint) (*models.Menu, error) {
	// 仅查一条菜单记录
	var m models.Menu
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &m, nil
}

// FindPermissions 查询权限集合。
func (d *MenuPermissionDao) FindPermissions(ctx context.Context, ids []uint) ([]models.Permission, error) {
	// 批量查询权限
	var perms []models.Permission
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Find(&perms, ids).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return perms, nil
}

// ReplacePermissions 全量替换菜单绑定权限。
func (d *MenuPermissionDao) ReplacePermissions(ctx context.Context, menu *models.Menu, perms []models.Permission) error {
	// 使用 GORM 关联替换
	return d.db.WithContext(ctx).Model(menu).Association("Permissions").Replace(&perms)
}

// GetPermissions 查询菜单已绑定权限。
func (d *MenuPermissionDao) GetPermissions(ctx context.Context, menuID uint) ([]models.Permission, error) {
	// 预加载权限关联
	var menu models.Menu
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Preload("Permissions").First(&menu, menuID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Permission{}, nil
		}
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return menu.Permissions, nil
}
