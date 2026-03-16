package rbac

import (
	"context"
	"errors"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type AuthService struct {
	authDao *rbacdao.AuthDao // 认证 DAO
}

// NewAuthService 创建认证服务。
func NewAuthService(authDao *rbacdao.AuthDao) *AuthService {
	return &AuthService{authDao: authDao}
}

// Register 处理用户注册业务逻辑
func (s *AuthService) Register(ctx context.Context, username, password, email string) error {
	// 1) 输入标准化：去掉首尾空格，避免“视觉相同账号”造成重复或绕过校验。
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	if username == "" {
		return errors.New("用户名不能为空")
	}

	// 2) 检查用户名唯一性。
	existingUser, err := s.authDao.GetUserByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// 3) 密码强度校验 + bcrypt 哈希存储（禁止明文落库）。
	if err := ValidatePasswordStrength(password); err != nil {
		return err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4) 构建并持久化用户。
	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPwd),
		Email:        email,
		Status:       1, // 正常状态
	}

	// 保存用户到数据库
	return s.authDao.CreateUser(ctx, user)
}

// Login 处理用户登录业务逻辑。
// 安全策略：
// - 用户名不存在与密码错误统一返回相同提示，避免账号枚举；
// - 禁用用户直接拒绝登录。
func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, error) {
	username = strings.TrimSpace(username)

	// 1) 通过用户名查询用户。
	user, err := s.authDao.GetUserByUsername(ctx, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 2) 比较密码哈希。
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 3) 禁用状态用户不允许登录。
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	return user, nil
}

// GetUserByID 根据 ID 查询用户，用于鉴权场景下的状态校验。
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	return s.authDao.GetUserByID(ctx, userID)
}
