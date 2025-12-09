package dao

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type MenuDao struct {
	db *gorm.DB
}

func NewMenuDao(db *gorm.DB) *MenuDao { return &MenuDao{db: db} }

// Create a menu entry
func (d *MenuDao) Create(ctx context.Context, m *models.Menu) error {
	return d.db.WithContext(ctx).Create(m).Error
}

// FindAll returns menus (preload permissions)
func (d *MenuDao) FindAll(ctx context.Context) ([]models.Menu, error) {
	var out []models.Menu
	err := d.db.WithContext(ctx).Preload("Permissions").Order("`order` ASC").Find(&out).Error
	return out, err
}

func (d *MenuDao) FindByID(ctx context.Context, id uint) (*models.Menu, error) {
	var m models.Menu
	if err := d.db.WithContext(ctx).Preload("Permissions").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *MenuDao) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Menu{}, id).Error
}
