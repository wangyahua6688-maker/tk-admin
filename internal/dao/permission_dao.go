package dao

import (
	"context"
	"go-admin-full/internal/models"

	"gorm.io/gorm"
)

type PermissionDAO interface {
	Create(ctx context.Context, p *models.Permission) error
	Update(ctx context.Context, p *models.Permission) error
	List(ctx context.Context) ([]models.Permission, error)
	Get(ctx context.Context, id uint) (*models.Permission, error)
	Delete(ctx context.Context, id uint) error
	GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Permission, error)
}

type permissionDAOImpl struct {
	db *gorm.DB
}

func NewPermissionDAO(db *gorm.DB) PermissionDAO {
	return &permissionDAOImpl{db: db}
}

func (d *permissionDAOImpl) Create(ctx context.Context, p *models.Permission) error {
	return d.db.WithContext(ctx).Create(p).Error
}

func (d *permissionDAOImpl) Update(ctx context.Context, p *models.Permission) error {
	return d.db.WithContext(ctx).Save(p).Error
}

func (d *permissionDAOImpl) List(ctx context.Context) ([]models.Permission, error) {
	var list []models.Permission
	err := d.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (d *permissionDAOImpl) Get(ctx context.Context, id uint) (*models.Permission, error) {
	var p models.Permission
	if err := d.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (d *permissionDAOImpl) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Permission{}, id).Error
}

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
