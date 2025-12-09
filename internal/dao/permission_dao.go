package dao

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type PermissionDao struct {
	db *gorm.DB
}

func NewPermissionDao(db *gorm.DB) *PermissionDao { return &PermissionDao{db: db} }

func (d *PermissionDao) Create(ctx context.Context, p *models.Permission) error {
	return d.db.WithContext(ctx).Create(p).Error
}

func (d *PermissionDao) FindAll(ctx context.Context) ([]models.Permission, error) {
	var out []models.Permission
	err := d.db.WithContext(ctx).Find(&out).Error
	return out, err
}

func (d *PermissionDao) FindByID(ctx context.Context, id uint) (*models.Permission, error) {
	var p models.Permission
	if err := d.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (d *PermissionDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Permission{}, id).Error
}
