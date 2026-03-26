package rbac

import (
	"context"
	"errors"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// UserService 定义UserService相关结构。
type UserService struct {
	userDao *rbacdao.UserDao // 用户 DAO
}

// NewUserService 创建用户服务。
func NewUserService(userDao *rbacdao.UserDao) *UserService {
	// 返回当前处理结果。
	return &UserService{userDao: userDao}
}

// CreateUser 创建User。
func (s *UserService) CreateUser(ctx context.Context, username, password, email, avatar string) (*models.User, error) {
	username = strings.TrimSpace(username) // 规范化用户名
	email = strings.TrimSpace(email)       // 规范化邮箱
	avatar = strings.TrimSpace(avatar)     // 规范化头像
	// 判断条件并进入对应分支逻辑。
	if username == "" {
		// 返回当前处理结果。
		return nil, errors.New("用户名不能为空")
	}

	// 检查用户名是否存在
	existing, err := s.userDao.GetByUsername(ctx, username)
	// 判断条件并进入对应分支逻辑。
	if err == nil && existing != nil {
		// 返回当前处理结果。
		return nil, errors.New("用户名已存在")
	}

	// 密码加密
	if err := ValidatePasswordStrength(password); err != nil {
		// 返回当前处理结果。
		return nil, err
	}

	// 定义并初始化当前变量。
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return nil, err
	}

	// 创建用户
	user := &models.User{
		// 处理当前语句逻辑。
		Username: username,
		// 调用string完成当前处理。
		PasswordHash: string(hashedPassword),
		// 处理当前语句逻辑。
		Email: email,
		// 处理当前语句逻辑。
		Avatar: avatar,
		// 处理当前语句逻辑。
		Status: 1,
	}

	// 判断条件并进入对应分支逻辑。
	if err := s.userDao.Create(ctx, user); err != nil {
		// 返回当前处理结果。
		return nil, err
	}

	// 返回当前处理结果。
	return user, nil
}

// ListAllUsers 查询AllUsers列表。
func (s *UserService) ListAllUsers(ctx context.Context) ([]models.User, error) {
	// 返回当前处理结果。
	return s.userDao.ListAll(ctx)
}

// GetUserByID 查询用户基础信息。
func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	// 返回当前处理结果。
	return s.userDao.GetByID(ctx, userID)
}

// UpdateUser 更新用户资料。
// 安全规则：
// 1. 用户名不可通过该接口修改；
// 2. 密码非空时会重新哈希后存储。
func (s *UserService) UpdateUser(ctx context.Context, userID uint, email string, status *int, password string, avatar *string) error {
	// 定义并初始化当前变量。
	updates := make(map[string]interface{})

	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(email) != "" {
		// 更新当前变量或字段值。
		updates["email"] = strings.TrimSpace(email)
	}

	// 判断条件并进入对应分支逻辑。
	if status != nil {
		// 更新当前变量或字段值。
		updates["status"] = *status
	}

	// avatar 允许显式置空，因此以指针判断“是否传入”。
	if avatar != nil {
		// 更新当前变量或字段值。
		updates["avatar"] = strings.TrimSpace(*avatar)
	}

	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(password) != "" {
		// 判断条件并进入对应分支逻辑。
		if err := ValidatePasswordStrength(password); err != nil {
			// 返回当前处理结果。
			return err
		}
		// 定义并初始化当前变量。
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 返回当前处理结果。
			return err
		}
		// 更新当前变量或字段值。
		updates["password_hash"] = string(hashedPassword)
	}

	// 判断条件并进入对应分支逻辑。
	if len(updates) == 0 {
		// 返回当前处理结果。
		return errors.New("无可更新字段")
	}

	// 返回当前处理结果。
	return s.userDao.UpdateByID(ctx, userID, updates)
}

// DeleteUser 删除用户。
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	// 返回当前处理结果。
	return s.userDao.DeleteByID(ctx, userID)
}
