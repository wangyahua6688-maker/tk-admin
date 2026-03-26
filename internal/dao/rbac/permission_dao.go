package rbac

import (
	"context"
	"go-admin/internal/models"

	"gorm.io/gorm"
)

// PermissionDAO 定义权限数据访问接口。
type PermissionDAO interface {
	// 调用Create完成当前处理。
	Create(ctx context.Context, p *models.Permission) error
	// 调用Update完成当前处理。
	Update(ctx context.Context, p *models.Permission) error
	// 调用List完成当前处理。
	List(ctx context.Context) ([]models.Permission, error)
	// 调用Get完成当前处理。
	Get(ctx context.Context, id uint) (*models.Permission, error)
	// 调用Delete完成当前处理。
	Delete(ctx context.Context, id uint) error
	// 调用GetByRoleIDs完成当前处理。
	GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Permission, error)
}

// permissionDAOImpl 使用 GORM 实现权限 DAO。
type permissionDAOImpl struct {
	db *gorm.DB // 数据库连接
}

// NewPermissionDAO 创建权限 DAO 实例。
func NewPermissionDAO(db *gorm.DB) PermissionDAO {
	// 返回当前处理结果。
	return &permissionDAOImpl{db: db}
}

// Create 新增权限记录。
func (d *permissionDAOImpl) Create(ctx context.Context, p *models.Permission) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(p).Error
}

// Update 更新权限记录。
func (d *permissionDAOImpl) Update(ctx context.Context, p *models.Permission) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Save(p).Error
}

// List 查询全部权限。
func (d *permissionDAOImpl) List(ctx context.Context) ([]models.Permission, error) {
	// 声明当前变量。
	var list []models.Permission
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).Find(&list).Error
	// 返回当前处理结果。
	return list, err
}

// Get 根据 ID 获取权限。
func (d *permissionDAOImpl) Get(ctx context.Context, id uint) (*models.Permission, error) {
	// 声明当前变量。
	var p models.Permission
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&p, id).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &p, nil
}

// Delete 删除权限记录。
func (d *permissionDAOImpl) Delete(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.Permission{}, id).Error
}

// GetByRoleIDs 根据角色 ID 列表查询权限集合。
func (d *permissionDAOImpl) GetByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Permission, error) {
	// 声明当前变量。
	var list []models.Permission

	// 通过角色-权限关联表加载权限，去重并过滤软删除数据。
	err := d.db.WithContext(ctx).
		// 调用Table完成当前处理。
		Table("sys_permissions p").
		// 调用Select完成当前处理。
		Select("DISTINCT p.*").
		// 更新当前变量或字段值。
		Joins("JOIN sys_role_permissions rp ON rp.permission_id = p.id").
		// 更新当前变量或字段值。
		Joins("JOIN sys_roles r ON r.id = rp.role_id").
		// 调用Where完成当前处理。
		Where("rp.role_id IN ? AND p.deleted_at IS NULL AND r.deleted_at IS NULL", roleIDs).
		// 调用Find完成当前处理。
		Find(&list).Error

	// 返回当前处理结果。
	return list, err
}
