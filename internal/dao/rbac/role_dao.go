package rbac

import (
	"context"
	"go-admin/internal/models"

	"gorm.io/gorm"
)

// RoleDAO 定义角色数据访问接口。
type RoleDAO interface {
	// 调用Create完成当前处理。
	Create(ctx context.Context, r *models.Role) error
	// 调用Update完成当前处理。
	Update(ctx context.Context, r *models.Role) error
	// 调用List完成当前处理。
	List(ctx context.Context) ([]models.Role, error)
	// 调用Get完成当前处理。
	Get(ctx context.Context, id uint) (*models.Role, error)
	// 调用Delete完成当前处理。
	Delete(ctx context.Context, id uint) error
	// 调用GetByUserID完成当前处理。
	GetByUserID(ctx context.Context, userID uint) ([]models.Role, error)
}

// roleDAOImpl 使用 GORM 实现角色 DAO。
type roleDAOImpl struct {
	db *gorm.DB // 数据库连接
}

// NewRoleDAO 创建角色 DAO 实例。
func NewRoleDAO(db *gorm.DB) RoleDAO {
	// 返回当前处理结果。
	return &roleDAOImpl{db: db}
}

// Create 新增角色。
func (d *roleDAOImpl) Create(ctx context.Context, r *models.Role) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(r).Error
}

// Update 更新角色。
func (d *roleDAOImpl) Update(ctx context.Context, r *models.Role) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Save(r).Error
}

// List 查询全部角色。
func (d *roleDAOImpl) List(ctx context.Context) ([]models.Role, error) {
	// 声明当前变量。
	var list []models.Role
	// 定义并初始化当前变量。
	err := d.db.WithContext(ctx).Find(&list).Error
	// 返回当前处理结果。
	return list, err
}

// Get 根据 ID 查询角色。
func (d *roleDAOImpl) Get(ctx context.Context, id uint) (*models.Role, error) {
	// 声明当前变量。
	var r models.Role
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).First(&r, id).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &r, nil
}

// Delete 删除角色。
func (d *roleDAOImpl) Delete(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}

// GetByUserID 根据用户 ID 查询其关联角色列表。
func (d *roleDAOImpl) GetByUserID(ctx context.Context, userID uint) ([]models.Role, error) {
	// 声明当前变量。
	var list []models.Role

	// 通过 sys_user_roles 关联拿到用户角色，并过滤软删除角色。
	err := d.db.WithContext(ctx).
		// 调用Table完成当前处理。
		Table("sys_roles r").
		// 更新当前变量或字段值。
		Joins("JOIN sys_user_roles ur ON ur.role_id = r.id").
		// 更新当前变量或字段值。
		Where("ur.user_id = ? AND r.deleted_at IS NULL", userID).
		// 调用Find完成当前处理。
		Find(&list).Error

	// 返回当前处理结果。
	return list, err
}
