package rbac

import (
	"context"
	"go-admin-full/internal/models"

	"gorm.io/gorm"
)

// PermissionDAO 定义权限数据访问接口。
type PermissionDAO interface {
	Create(ctx context.Context, p *models.Permission) error
	Update(ctx context.Context, p *models.Permission) error
	List(ctx context.Context) ([]models.Permission, error)
	Get(ctx context.Context, id uint) (*models.Permission, error)
	Delete(ctx context.Context, id uint) error
	GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Permission, error)
}

// permissionDAOImpl 使用 GORM 实现权限 DAO。
type permissionDAOImpl struct {
	db *gorm.DB // 数据库连接
}

// NewPermissionDAO 创建权限 DAO 实例。
func NewPermissionDAO(db *gorm.DB) PermissionDAO {
	return &permissionDAOImpl{db: db}
}

// Create 新增权限记录。
func (d *permissionDAOImpl) Create(ctx context.Context, p *models.Permission) error {
	return d.db.WithContext(ctx).Create(p).Error
}

// Update 更新权限记录。
func (d *permissionDAOImpl) Update(ctx context.Context, p *models.Permission) error {
	return d.db.WithContext(ctx).Save(p).Error
}

// List 查询全部权限。
func (d *permissionDAOImpl) List(ctx context.Context) ([]models.Permission, error) {
	var list []models.Permission
	err := d.db.WithContext(ctx).Find(&list).Error
	return list, err
}

// Get 根据 ID 获取权限。
func (d *permissionDAOImpl) Get(ctx context.Context, id uint) (*models.Permission, error) {
	var p models.Permission
	if err := d.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// Delete 删除权限记录。
func (d *permissionDAOImpl) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Permission{}, id).Error
}

// GetByRoleIDs 根据角色 ID 列表查询权限集合。
func (d *permissionDAOImpl) GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Permission, error) {
	var list []models.Permission

	// 通过角色-权限关联表加载权限，去重并过滤软删除数据。
	err := d.db.WithContext(ctx).
		Table("sys_permissions p").
		Select("DISTINCT p.*").
		Joins("JOIN sys_role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN sys_roles r ON r.id = rp.role_id").
		Where("rp.role_id IN ? AND p.deleted_at IS NULL AND r.deleted_at IS NULL", roleIDs).
		Find(&list).Error

	return list, err
}
