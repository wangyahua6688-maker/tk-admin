package rbac

import (
	"context"
	"go-admin-full/internal/models"

	"gorm.io/gorm"
)

// MenuDAO 定义菜单数据访问接口。
type MenuDAO interface {
	Create(ctx context.Context, m *models.Menu) error
	Update(ctx context.Context, m *models.Menu) error
	List(ctx context.Context) ([]models.Menu, error)
	ListByUserID(ctx context.Context, userID uint) ([]models.Menu, error)
	Get(ctx context.Context, id uint) (*models.Menu, error)
	Delete(ctx context.Context, id uint) error
}

// menuDAOImpl 使用 GORM 实现菜单 DAO。
type menuDAOImpl struct {
	db *gorm.DB // 数据库连接
}

// NewMenuDAO 创建菜单 DAO 实例。
func NewMenuDAO(db *gorm.DB) MenuDAO {
	return &menuDAOImpl{db: db}
}

// Create 新增菜单。
func (d *menuDAOImpl) Create(ctx context.Context, m *models.Menu) error {
	return d.db.WithContext(ctx).Create(m).Error
}

// Update 更新菜单。
func (d *menuDAOImpl) Update(ctx context.Context, m *models.Menu) error {
	return d.db.WithContext(ctx).Save(m).Error
}

// List 查询全部菜单（按层级与排序号排序）。
func (d *menuDAOImpl) List(ctx context.Context) ([]models.Menu, error) {
	var list []models.Menu
	err := d.db.WithContext(ctx).
		// 先按父级排序，再按 order_num 排序，最后按 id 稳定排序
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
		// 仅检查是否存在 admin 角色
		Select("COUNT(1) > 0").
		Joins("JOIN sys_roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.code = ? AND r.deleted_at IS NULL", userID, "admin").
		Scan(&isAdmin).Error; err != nil {
		return nil, err
	}

	// admin 角色直接返回全部菜单
	if isAdmin {
		return d.List(ctx)
	}

	var list []models.Menu
	err := d.db.WithContext(ctx).
		// 通过用户角色与角色权限反查菜单权限，再去重得到用户可访问菜单。
		Table("sys_menus m").
		// 去重避免同一权限映射多次
		Select("DISTINCT m.*").
		Joins("JOIN sys_menu_permissions mp ON mp.menu_id = m.id").
		Joins("JOIN sys_role_permissions rp ON rp.permission_id = mp.permission_id").
		Joins("JOIN sys_user_roles ur ON ur.role_id = rp.role_id").
		Joins("JOIN sys_roles r ON r.id = ur.role_id").
		Joins("JOIN sys_permissions p ON p.id = rp.permission_id").
		// 过滤软删除数据
		Where("ur.user_id = ? AND m.deleted_at IS NULL AND r.deleted_at IS NULL AND p.deleted_at IS NULL", userID).
		Order("m.parent_id ASC, m.order_num ASC, m.id ASC").
		Scan(&list).Error

	return list, err
}

// Get 根据 ID 查询菜单。
func (d *menuDAOImpl) Get(ctx context.Context, id uint) (*models.Menu, error) {
	var m models.Menu
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// Delete 删除菜单。
func (d *menuDAOImpl) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Menu{}, id).Error
}
