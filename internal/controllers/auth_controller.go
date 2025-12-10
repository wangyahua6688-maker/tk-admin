package controllers

import (
	"go-admin-full/internal/constants"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/services"
	"go-admin-full/internal/tokenpkg"
	"go-admin-full/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	authService *services.AuthService
	tokenMgr    *tokenpkg.Manager
}

// 修改构造函数，同时传入db和token管理器
func NewAuthController(db *gorm.DB, mgr *tokenpkg.Manager) *AuthController {
	authDao := dao.NewAuthDao(db)
	authService := services.NewAuthService(authDao)
	return &AuthController{
		authService: authService,
		tokenMgr:    mgr,
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// 获取上下文中的logger
	logger := utils.LoggerFromContext(ctx.Request.Context())
	logger.Info("用户登录尝试: %s", req.Username)

	// 调用服务层登录
	user, err := c.authService.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		logger.Error("登录失败: %v", err)
		utils.JSONError(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	// 生成JWT Token
	accessToken, refreshToken, err := c.tokenMgr.GenerateTokens(user.ID)
	if err != nil {
		logger.Error("生成Token失败: %v", err)
		utils.JSONError(ctx, http.StatusInternalServerError, "生成Token失败")
		return
	}

	logger.Info("用户登录成功: %s (ID: %d)", req.Username, user.ID)

	utils.JSONOK(ctx, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_id":       user.ID,
		"username":      user.Username,
		"expires_in":    int(c.tokenMgr.Config.AccessExpire.Seconds()),
	})
}

// Refresh 刷新Token
func (c *AuthController) Refresh(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	logger := utils.LoggerFromContext(ctx.Request.Context())
	logger.Info("刷新Token请求")

	// 使用token管理器刷新token
	newAccessToken, err := c.tokenMgr.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.Error("刷新Token失败: %v", err)
		status := http.StatusUnauthorized
		if err == constants.ErrExpiredToken {
			status = http.StatusForbidden
		}
		utils.JSONError(ctx, status, err.Error())
		return
	}

	logger.Info("Token刷新成功")

	utils.JSONOK(ctx, gin.H{
		"access_token": newAccessToken,
		"expires_in":   int(c.tokenMgr.Config.AccessExpire.Seconds()),
	})
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Email    string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	logger := utils.LoggerFromContext(ctx.Request.Context())
	logger.Info("用户注册尝试: %s", req.Username)

	if err := c.authService.Register(ctx.Request.Context(), req.Username, req.Password, req.Email); err != nil {
		logger.Error("注册失败: %v", err)
		utils.JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	logger.Info("用户注册成功: %s", req.Username)

	utils.JSONOK(ctx, gin.H{"message": "注册成功"})
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// 从JWT中间件中获取用户ID
	uid, exists := ctx.Get("uid")
	if !exists {
		utils.JSONError(ctx, http.StatusUnauthorized, "用户未认证")
		return
	}

	userID := uid.(uint)
	logger := utils.LoggerFromContext(ctx.Request.Context())
	logger.Info("用户登出: ID=%d", userID)

	// 使刷新令牌无效
	if err := c.tokenMgr.InvalidateRefresh(userID); err != nil {
		logger.Error("登出失败，无法使令牌无效: %v", err)
		// 这里可以选择是否返回错误，通常登出即使失败也返回成功
	}

	logger.Info("用户登出成功: ID=%d", userID)

	utils.JSONOK(ctx, gin.H{"message": "登出成功"})
}

// Profile 获取用户信息（示例）
func (c *AuthController) Profile(ctx *gin.Context) {
	//// 从JWT中间件中获取用户ID
	//uid, exists := ctx.Get("uid")
	//if !exists {
	//	utils.JSONError(ctx, http.StatusUnauthorized, "用户未认证")
	//	return
	//}
	//
	//userID := uid.(uint)
	//logger := utils.LoggerFromContext(ctx.Request.Context())
	//logger.Info("获取用户信息: ID=%d", userID)
	//
	//user, err := c.authService.GetUserByID(ctx.Request.Context(), userID)
	//if err != nil {
	//	logger.Error("获取用户信息失败: %v", err)
	//	utils.JSONError(ctx, http.StatusNotFound, "用户不存在")
	//	return
	//}
	//
	//// 返回用户信息（注意不要返回敏感信息）
	//utils.JSONOK(ctx, gin.H{
	//	"user_id":    user.ID,
	//	"username":   user.Username,
	//	"email":      user.Email,
	//	"status":     user.Status,
	//	"created_at": user.CreatedAt,
	//})
}
