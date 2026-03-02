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

	// 通过 user_roles 关联拿到用户角色，并过滤软删除角色。
	err := d.db.WithContext(ctx).
		Table("roles r").
		Joins("JOIN user_roles ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.deleted_at IS NULL", userID).
		Find(&list).Error

	return list, err
}
