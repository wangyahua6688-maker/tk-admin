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
		Order("order_num ASC").
		Find(&list).Error
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

// 根据管理员权限获取菜单
func (d *menuDAOImpl) GetMenusByAdminID(ctx context.Context, adminID int64) ([]*models.Menu, error) {
	var menus []*models.Menu

	err := d.db.WithContext(ctx).
		Joins("JOIN role_menus rm ON rm.menu_id = menus.id").
		Joins("JOIN admin_roles ar ON rm.role_id = ar.role_id").
		Where("ar.admin_id = ?", adminID).
		Group("menus.id").
		Find(&menus).Error

	return menus, err
}
