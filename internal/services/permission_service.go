package services

import (
    "context"

    "go-admin-full/internal/models"
    "gorm.io/gorm"
)

type PermissionService struct {
    db *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService { return &PermissionService{db: db} }

func (s *PermissionService) Create(ctx context.Context, p *models.Permission) error {
    return s.db.WithContext(ctx).Create(p).Error
}
func (s *PermissionService) List(ctx context.Context) ([]models.Permission, error) {
    var items []models.Permission
    err := s.db.WithContext(ctx).Find(&items).Error
    return items, err
}
func (s *PermissionService) Get(ctx context.Context, id uint) (*models.Permission, error) {
    var p models.Permission
    if err := s.db.WithContext(ctx).First(&p, id).Error; err != nil {
        return nil, err
    }
    return &p, nil
}
func (s *PermissionService) Delete(ctx context.Context, id uint) error {
    return s.db.WithContext(ctx).Delete(&models.Permission{}, id).Error
}
