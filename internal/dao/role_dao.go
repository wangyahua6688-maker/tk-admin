package dao

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type RoleDao struct {
	db *gorm.DB
}

func NewRoleDao(db *gorm.DB) *RoleDao { return &RoleDao{db: db} }

func (d *RoleDao) Create(ctx context.Context, r *models.Role) error {
	return d.db.WithContext(ctx).Create(r).Error
}

func (d *RoleDao) FindAll(ctx context.Context) ([]models.Role, error) {
	var out []models.Role
	err := d.db.WithContext(ctx).Preload("Permissions").Find(&out).Error
	return out, err
}

func (d *RoleDao) FindByID(ctx context.Context, id uint) (*models.Role, error) {
	var r models.Role
	err := d.db.WithContext(ctx).Preload("Permissions").First(&r, id).Error
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (d *RoleDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}
