package rbac

import (
	"context"
	"errors"
	"fmt"

	"go-admin/internal/models"
	"gorm.io/gorm"
)

// UserDao 用户数据访问层。
type UserDao struct {
	db *gorm.DB // 数据库连接
}

// NewUserDao 创建用户 DAO。
func NewUserDao(db *gorm.DB) *UserDao {
	// 返回当前处理结果。
	return &UserDao{db: db}
}

// GetByUsername 根据用户名查询用户。
func (d *UserDao) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	// 声明当前变量。
	var user models.User
	// 按用户名查询
	if err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &user, nil
}

// Create 创建用户。
func (d *UserDao) Create(ctx context.Context, user *models.User) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据用户ID查询用户。
func (d *UserDao) GetByID(ctx context.Context, id uint) (*models.User, error) {
	// 声明当前变量。
	var user models.User
	// 按 ID 查询
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return &user, nil
}

// ListAll 查询全部用户。
func (d *UserDao) ListAll(ctx context.Context) ([]models.User, error) {
	// 声明当前变量。
	var users []models.User
	// 判断条件并进入对应分支逻辑。
	if err := d.db.WithContext(ctx).Find(&users).Error; err != nil {
		// 返回当前处理结果。
		return nil, err
	}
	// 返回当前处理结果。
	return users, nil
}

// UpdateByID 更新用户可编辑字段。
func (d *UserDao) UpdateByID(ctx context.Context, id uint, updates map[string]interface{}) error {
	// 判断条件并进入对应分支逻辑。
	if len(updates) == 0 {
		// 返回当前处理结果。
		return fmt.Errorf("updates is empty")
	}
	// 更新指定字段
	return d.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByID 删除用户。
func (d *UserDao) DeleteByID(ctx context.Context, id uint) error {
	// 返回当前处理结果。
	return d.db.WithContext(ctx).Delete(&models.User{}, id).Error
}
