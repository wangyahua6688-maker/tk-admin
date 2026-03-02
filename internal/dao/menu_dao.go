package dao

import (
	"context"
	"go-admin-full/internal/models"

	"gorm.io/gorm"
)

type MenuDAO interface {
	Create(ctx context.Context, m *models.Menu) error
	Update(ctx context.Context, m *models.Menu) error
	List(ctx context.Context) ([]models.Menu, error)
	ListByUserID(ctx context.Context, userID uint) ([]models.Menu, error)
	Get(ctx context.Context, id uint) (*models.Menu, error)
	Delete(ctx context.Context, id uint) error
}

type menuDAOImpl struct {
	db *gorm.DB
}

func NewMenuDAO(db *gorm.DB) MenuDAO {
	return &menuDAOImpl{db: db}
}

func (d *menuDAOImpl) Create(ctx context.Context, m *models.Menu) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *menuDAOImpl) Update(ctx context.Context, m *models.Menu) error {
	return d.db.WithContext(ctx).Save(m).Error
}

func (d *menuDAOImpl) List(ctx context.Context) ([]models.Menu, error) {
	var list []models.Menu
	err := d.db.WithContext(ctx).
		Order("parent_id ASC, order_num ASC, id ASC").
		Find(&list).Error
	return list, err
}

// ListByUserID 根据用户拥有的权限查询可访问菜单。
// 逻辑说明：
// 1. admin 角色默认返回全部菜单；
// 2. 普通角色通过 sys_user_roles -> sys_role_permissions -> sys_menu_permissions 关联获取。
func (d *menuDAOImpl) ListByUserID(ctx context.Context, userID uint) ([]models.Menu, error) {
	var isAdmin bool
	if err := d.db.WithContext(ctx).
		// 统一使用 sys_* 真实表名，避免因历史命名导致联调时 SQL 报表不存在。
		Table("sys_user_roles ur").
		Select("COUNT(1) > 0").
		Joins("JOIN sys_roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.code = ? AND r.deleted_at IS NULL", userID, "admin").
		Scan(&isAdmin).Error; err != nil {
		return nil, err
	}

	if isAdmin {
		return d.List(ctx)
	}

	var list []models.Menu
	err := d.db.WithContext(ctx).
		// 通过用户角色与角色权限反查菜单权限，再去重得到用户可访问菜单。
		Table("sys_menus m").
		Select("DISTINCT m.*").
		Joins("JOIN sys_menu_permissions mp ON mp.menu_id = m.id").
		Joins("JOIN sys_role_permissions rp ON rp.permission_id = mp.permission_id").
		Joins("JOIN sys_user_roles ur ON ur.role_id = rp.role_id").
		Joins("JOIN sys_roles r ON r.id = ur.role_id").
		Joins("JOIN sys_permissions p ON p.id = rp.permission_id").
		Where("ur.user_id = ? AND m.deleted_at IS NULL AND r.deleted_at IS NULL AND p.deleted_at IS NULL", userID).
		Order("m.parent_id ASC, m.order_num ASC, m.id ASC").
		Scan(&list).Error

	return list, err
}

func (d *menuDAOImpl) Get(ctx context.Context, id uint) (*models.Menu, error) {
	var m models.Menu
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *menuDAOImpl) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Menu{}, id).Error
}
