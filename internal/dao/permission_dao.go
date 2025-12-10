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

	err := d.db.WithContext(ctx).
		Table("permissions p").
		Select("DISTINCT p.*").
		Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN roles r ON r.id = rp.role_id").
		Where("rp.role_id IN ? AND p.deleted_at IS NULL AND r.deleted_at IS NULL", roleIDs).
		Find(&list).Error

	return list, err
}

func (d *permissionDAOImpl) GetByAdminID(ctx context.Context, adminID int64) ([]*models.Permission, error) {
	var list []*models.Permission

	err := d.db.WithContext(ctx).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN admin_roles ar ON rp.role_id = ar.role_id").
		Where("ar.admin_id = ?", adminID).
		Group("permissions.id").
		Find(&list).Error

	return list, err
}

// 用于路由自动注册
func (d *permissionDAOImpl) CreateIfNotExists(ctx context.Context, method, path string) error {
	var count int64
	d.db.WithContext(ctx).
		Model(&models.Permission{}).
		Where("method = ? AND path = ?", method, path).
		Count(&count)

	if count > 0 {
		return nil
	}

	return d.db.WithContext(ctx).Create(&models.Permission{
		Method: method,
		Path:   path,
		Name:   method + " " + path,
	}).Error
}
