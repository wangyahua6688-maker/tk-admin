package rbac

import (
	"context"
	"go-admin-full/internal/models"

	"gorm.io/gorm"
)

// RoleDAO 定义角色数据访问接口。
type RoleDAO interface {
	Create(ctx context.Context, r *models.Role) error
	Update(ctx context.Context, r *models.Role) error
	List(ctx context.Context) ([]models.Role, error)
	Get(ctx context.Context, id uint) (*models.Role, error)
	Delete(ctx context.Context, id uint) error
	GetByUserID(ctx context.Context, userID uint) ([]models.Role, error)
}

// roleDAOImpl 使用 GORM 实现角色 DAO。
type roleDAOImpl struct {
	db *gorm.DB // 数据库连接
}

// NewRoleDAO 创建角色 DAO 实例。
func NewRoleDAO(db *gorm.DB) RoleDAO {
	return &roleDAOImpl{db: db}
}

// Create 新增角色。
func (d *roleDAOImpl) Create(ctx context.Context, r *models.Role) error {
	return d.db.WithContext(ctx).Create(r).Error
}

// Update 更新角色。
func (d *roleDAOImpl) Update(ctx context.Context, r *models.Role) error {
	return d.db.WithContext(ctx).Save(r).Error
}

// List 查询全部角色。
func (d *roleDAOImpl) List(ctx context.Context) ([]models.Role, error) {
	var list []models.Role
	err := d.db.WithContext(ctx).Find(&list).Error
	return list, err
}

// Get 根据 ID 查询角色。
func (d *roleDAOImpl) Get(ctx context.Context, id uint) (*models.Role, error) {
	var r models.Role
	if err := d.db.WithContext(ctx).First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// Delete 删除角色。
func (d *roleDAOImpl) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}

// GetByUserID 根据用户 ID 查询其关联角色列表。
func (d *roleDAOImpl) GetByUserID(ctx context.Context, userID uint) ([]models.Role, error) {
	var list []models.Role

	// 通过 sys_user_roles 关联拿到用户角色，并过滤软删除角色。
	err := d.db.WithContext(ctx).
		Table("sys_roles r").
		Joins("JOIN sys_user_roles ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.deleted_at IS NULL", userID).
		Find(&list).Error

	return list, err
}
