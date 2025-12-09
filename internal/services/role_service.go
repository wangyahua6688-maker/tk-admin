package services

import (
    "context"

    "go-admin-full/internal/models"
    "gorm.io/gorm"
)

type RoleService struct {
    db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
    return &RoleService{db: db}
}

func (s *RoleService) Create(ctx context.Context, role *models.Role) error {
    return s.db.WithContext(ctx).Create(role).Error
}

func (s *RoleService) List(ctx context.Context) ([]models.Role, error) {
    var roles []models.Role
    err := s.db.WithContext(ctx).Preload("Permissions").Find(&roles).Error
    return roles, err
}

func (s *RoleService) Get(ctx context.Context, id uint) (*models.Role, error) {
    var r models.Role
    if err := s.db.WithContext(ctx).Preload("Permissions").First(&r, id).Error; err != nil {
        return nil, err
    }
    return &r, nil
}

func (s *RoleService) Delete(ctx context.Context, id uint) error {
    return s.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}
