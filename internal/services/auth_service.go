package services

import (
	"context"
	"errors"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	authDao *dao.AuthDao
}

func NewAuthService(authDao *dao.AuthDao) *AuthService {
	return &AuthService{authDao: authDao}
}

// Register 处理用户注册业务逻辑
func (s *AuthService) Register(ctx context.Context, username, password, email string) error {
	// 检查用户名是否存在
	existingUser, err := s.authDao.GetUserByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// 密码加密
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 构建用户模型
	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPwd),
		Email:        email,
		Status:       1, // 正常状态
	}

	// 保存用户到数据库
	return s.authDao.CreateUser(ctx, user)
}

// Login 处理用户登录业务逻辑
func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, error) {
	// 查询用户
	user, err := s.authDao.GetUserByUsername(ctx, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}

// RefreshToken 处理刷新Token业务逻辑
func (s *AuthService) RefreshToken(ctx context.Context, userID uint, refreshToken string) (*models.User, error) {
	// 查询用户
	user, err := s.authDao.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 验证刷新Token（此处根据实际业务逻辑调整，比如对比数据库中的refresh_token）
	if user.RefreshToken != refreshToken {
		return nil, errors.New("无效的刷新Token")
	}

	return user, nil
}

func (s *AuthService) UpdateUserToken(ctx context.Context, userID uint, refreshToken string) error {
	return s.authDao.UpdateUserToken(ctx, userID, refreshToken)
}

// Logout 处理用户登出业务逻辑
func (s *AuthService) Logout(ctx context.Context, userID uint) error {
	// 清空刷新Token（此处根据实际业务逻辑调整）
	return s.authDao.UpdateUserToken(ctx, userID, "")
}
