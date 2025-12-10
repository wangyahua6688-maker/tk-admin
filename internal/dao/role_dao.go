package dao

import (
	"context"
	"go-admin-full/internal/models"

	"gorm.io/gorm"
)

type RoleDAO interface {
	Create(ctx context.Context, r *models.Role) error
	Update(ctx context.Context, r *models.Role) error
	List(ctx context.Context) ([]models.Role, error)
	Get(ctx context.Context, id uint) (*models.Role, error)
	Delete(ctx context.Context, id uint) error
	GetByUserID(ctx context.Context, userID uint) ([]models.Role, error)
}

type roleDAOImpl struct {
	db *gorm.DB
}

func NewRoleDAO(db *gorm.DB) RoleDAO {
	return &roleDAOImpl{db: db}
}

func (d *roleDAOImpl) Create(ctx context.Context, r *models.Role) error {
	return d.db.WithContext(ctx).Create(r).Error
}

func (d *roleDAOImpl) Update(ctx context.Context, r *models.Role) error {
	return d.db.WithContext(ctx).Save(r).Error
}

func (d *roleDAOImpl) List(ctx context.Context) ([]models.Role, error) {
	var list []models.Role
	err := d.db.WithContext(ctx).Find(&list).Error
	return list, err
}

func (d *roleDAOImpl) Get(ctx context.Context, id uint) (*models.Role, error) {
	var r models.Role
	if err := d.db.WithContext(ctx).First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (d *roleDAOImpl) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}

func (d *roleDAOImpl) GetByUserID(ctx context.Context, userID uint) ([]models.Role, error) {
	var list []models.Role

	err := d.db.WithContext(ctx).
		Table("roles r").
		Joins("JOIN user_roles ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.deleted_at IS NULL", userID).
		Find(&list).Error

	return list, err
}

// AssignPermissions 给角色分配权限（全量替换）。
func (d *roleDAOImpl) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		for _, pid := range permissionIDs {
			rp := models.RolePermission{RoleID: roleID, PermissionID: pid}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// 角色对应的权限
func (d *roleDAOImpl) ListPermissions(ctx context.Context, roleID uint) ([]*models.Permission, error) {
	var items []*models.Permission
	err := d.db.WithContext(ctx).
		Table("permissions p").
		Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
		Where("rp.role_id = ?", roleID).
		Scan(&items).Error
	return items, err
}
