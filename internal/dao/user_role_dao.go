package dao

import (
	"context"
	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type UserRoleDao struct {
	db *gorm.DB
}

func NewUserRoleDao(db *gorm.DB) *UserRoleDao {
	return &UserRoleDao{db: db}
}

func (d *UserRoleDao) FindUser(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserRoleDao) FindRoles(ctx context.Context, roleIDs []uint) ([]models.Role, error) {
	var roles []models.Role
	if err := d.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (d *UserRoleDao) ReplaceRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	return d.db.WithContext(ctx).Model(user).Association("Roles").Replace(&roles)
}

func (d *UserRoleDao) AppendRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	return d.db.WithContext(ctx).Model(user).Association("Roles").Append(&roles)
}

func (d *UserRoleDao) RemoveRoles(ctx context.Context, user *models.User, roles []models.Role) error {
	return d.db.WithContext(ctx).Model(user).Association("Roles").Delete(&roles)
}

func (d *UserRoleDao) GetUserRolesWithPermissions(ctx context.Context, id uint) ([]models.Role, error) {
	var user models.User
	err := d.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&user, id).Error

	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}
