package services

import (
	"context"

	"go-admin-full/internal/models"
	"gorm.io/gorm"
)

type UserRoleService struct {
	db *gorm.DB
}

func NewUserRoleService(db *gorm.DB) *UserRoleService {
	return &UserRoleService{db: db}
}

//
// 用户绑定角色（覆盖）
//

func (s *UserRoleService) BindRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	var roles []models.Role
	if err := s.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	// Replace = 清空原角色并绑定新角色
	return s.db.WithContext(ctx).Model(&user).Association("Roles").Replace(&roles)
}

//
// 用户添加角色（追加）
//

func (s *UserRoleService) AddRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	var roles []models.Role
	if err := s.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&user).Association("Roles").Append(&roles)
}

//
// 用户移除部分角色
//

func (s *UserRoleService) RemoveRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}

	var roles []models.Role
	if err := s.db.WithContext(ctx).Find(&roles, roleIDs).Error; err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(&user).Association("Roles").Delete(&roles)
}

//
// 查询用户角色（自动携带角色的权限，用于菜单过滤、权限中间件等）
//

func (s *UserRoleService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	var user models.User

	// 自动预加载 user.roles.permissions
	err := s.db.WithContext(ctx).
		Preload("Roles.Permissions").
		First(&user, userID).Error

	if err != nil {
		return nil, err
	}

	return user.Roles, nil
}
