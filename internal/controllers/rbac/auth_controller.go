package rbac

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-admin-full/internal/constants"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	rbacsvc "go-admin-full/internal/services/rbac"
	tokenjwt "go-admin-full/internal/token/jwt"
	"go-admin-full/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	authService *rbacsvc.AuthService     // 认证业务服务
	loginLogSvc *rbacsvc.LoginLogService // 登录日志服务
	tokenMgr    *tokenjwt.Manager        // JWT 管理器
}

const (
	// 登录失败限制：15 分钟最多 5 次
	loginFailMax    = 5
	loginFailWindow = 15 * time.Minute
)

// 修改构造函数，同时传入db和token管理器
func NewAuthController(db *gorm.DB, mgr *tokenjwt.Manager) *AuthController {
	authDao := rbacdao.NewAuthDao(db)                                     // 初始化认证 DAO
	authService := rbacsvc.NewAuthService(authDao)                        // 初始化认证服务
	loginLogSvc := rbacsvc.NewLoginLogService(rbacdao.NewLoginLogDao(db)) // 初始化登录日志服务
	return &AuthController{
		authService: authService, // 注入认证服务
		loginLogSvc: loginLogSvc, // 注入登录日志服务
		tokenMgr:    mgr,         // 注入 Token 管理器
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定并校验请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	// 获取上下文中的logger
	logger := utils.LoggerFromContext(ctx.Request.Context())
	req.Username = strings.TrimSpace(req.Username) // 规范化用户名
	clientIP := ctx.ClientIP()                     // 取客户端 IP
	deviceID := resolveDeviceID(ctx)               // 解析设备标识
	logger.Info("用户登录尝试: %s", req.Username)

	// 安全防护：登录失败次数限制
	if blocked := c.isLoginBlocked(req.Username, clientIP); blocked {
		c.writeLoginLog(ctx, 0, req.Username, clientIP, deviceID, 0) // 记录失败日志
		utils.JSONError(ctx, http.StatusTooManyRequests, "登录失败次数过多，请稍后重试")
		return
	}

	// 调用服务层登录
	user, err := c.authService.Login(ctx.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.recordLoginFailure(req.Username, clientIP)                 // 记录失败次数
		c.writeLoginLog(ctx, 0, req.Username, clientIP, deviceID, 0) // 记录失败日志
		logger.Error("登录失败: %v", err)
		utils.JSONError(ctx, http.StatusUnauthorized, err.Error())
		return
	}
	c.clearLoginFailure(req.Username, clientIP) // 清空失败计数

	// 生成JWT Token
	accessToken, refreshToken, err := c.tokenMgr.GenerateTokensWithDevice(user.ID, deviceID)
	if err != nil {
		logger.Error("生成Token失败: %v", err)
		utils.JSONError(ctx, http.StatusInternalServerError, "生成Token失败")
		return
	}
	c.writeLoginLog(ctx, user.ID, user.Username, clientIP, deviceID, 1) // 记录成功日志

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

	// 绑定并校验请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	logger := utils.LoggerFromContext(ctx.Request.Context())
	logger.Info("刷新Token请求")

	// 安全防护：刷新前先校验 refresh token 基础有效性和用户状态。
	req.RefreshToken = strings.TrimSpace(req.RefreshToken)                                   // 规范化 token
	claims, err := tokenjwt.ParseTokenClaims(req.RefreshToken, c.tokenMgr.Config.SigningKey) // 解析 claims
	if err != nil {
		status := http.StatusUnauthorized
		if err == constants.ErrExpiredToken {
			status = http.StatusForbidden
		}
		utils.JSONError(ctx, status, err.Error())
		return
	}
	if claims.TokenType != tokenjwt.TokenTypeRefresh {
		utils.JSONError(ctx, http.StatusUnauthorized, constants.ErrInvalidToken.Error())
		return
	}
	// 校验用户状态, 封禁用户禁止刷新
	user, err := c.authService.GetUserByID(ctx.Request.Context(), claims.UserID) // 查询用户状态
	if err != nil || user.Status != 1 {
		utils.JSONError(ctx, http.StatusUnauthorized, "账号不可用")
		return
	}

	// 使用token管理器刷新token
	newAccessToken, newRefreshToken, err := c.tokenMgr.RefreshTokenPair(req.RefreshToken) // 生成新 token
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
		"access_token":  newAccessToken,                                // 新 access token
		"refresh_token": newRefreshToken,                               // 新 refresh token
		"expires_in":    int(c.tokenMgr.Config.AccessExpire.Seconds()), // access token 有效期
	})
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		Email    string `json:"email"`
	}

	// 绑定并校验请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误："+err.Error())
		return
	}

	logger := utils.LoggerFromContext(ctx.Request.Context())
	req.Username = strings.TrimSpace(req.Username) // 规范化用户名
	req.Email = strings.TrimSpace(req.Email)       // 规范化邮箱
	logger.Info("用户注册尝试: %s", req.Username)

	// 调用注册服务
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

	userID := uid.(uint) // 转换为 uint
	logger := utils.LoggerFromContext(ctx.Request.Context())
	logger.Info("用户登出: ID=%d", userID)

	deviceID := resolveDeviceID(ctx) // 解析设备标识
	if claimsRaw, ok := ctx.Get("claims"); ok {
		if claims, ok := claimsRaw.(*tokenjwt.Claims); ok && strings.TrimSpace(claims.DeviceID) != "" {
			deviceID = claims.DeviceID
		}
	}

	// 使 refresh token 无效（当前设备）
	if err := c.tokenMgr.InvalidateRefresh(userID, deviceID); err != nil {
		logger.Error("登出失败，无法使 refresh token 无效: %v", err)
	}

	// 将当前 access token 放入黑名单
	accessToken, _ := ctx.Get("access_token")
	if tokenStr, ok := accessToken.(string); ok && tokenStr != "" {
		if err := c.tokenMgr.RevokeAccessToken(tokenStr); err != nil {
			logger.Error("登出失败，无法撤销 access token: %v", err)
		}
	}

	logger.Info("用户登出成功: ID=%d", userID)

	utils.JSONOK(ctx, gin.H{"message": "登出成功"})
}

func resolveDeviceID(ctx *gin.Context) string {
	deviceID := strings.TrimSpace(ctx.GetHeader("X-Device-ID"))
	if deviceID != "" {
		return deviceID
	}
	ua := strings.TrimSpace(ctx.Request.UserAgent())
	if ua != "" {
		return ua
	}
	return "default"
}

func (c *AuthController) loginFailKey(username, ip string) string {
	return "auth:login_fail:" + strings.ToLower(strings.TrimSpace(username)) + ":" + strings.TrimSpace(ip)
}

// isLoginBlocked 判断是否触发登录失败限制。
func (c *AuthController) isLoginBlocked(username, ip string) bool {
	if c.tokenMgr == nil || c.tokenMgr.Store == nil {
		return false
	}

	key := c.loginFailKey(username, ip)
	v, err := c.tokenMgr.Store.Get(key)
	if err != nil {
		return false
	}

	count, _ := strconv.Atoi(v)
	return count >= loginFailMax
}

// recordLoginFailure 记录登录失败次数。
func (c *AuthController) recordLoginFailure(username, ip string) {
	if c.tokenMgr == nil || c.tokenMgr.Store == nil {
		return
	}

	key := c.loginFailKey(username, ip)
	current := 0
	if v, err := c.tokenMgr.Store.Get(key); err == nil {
		current, _ = strconv.Atoi(v)
	}
	_ = c.tokenMgr.Store.Set(key, strconv.Itoa(current+1), loginFailWindow)
}

// clearLoginFailure 清空登录失败计数。
func (c *AuthController) clearLoginFailure(username, ip string) {
	if c.tokenMgr == nil || c.tokenMgr.Store == nil {
		return
	}
	_ = c.tokenMgr.Store.Delete(c.loginFailKey(username, ip))
}

// writeLoginLog 写入登录审计日志，失败不影响主流程。
func (c *AuthController) writeLoginLog(ctx *gin.Context, userID uint, username, ip, device string, status int) {
	if c.loginLogSvc == nil {
		return
	}
	_ = c.loginLogSvc.CreateLoginLog(ctx.Request.Context(), &models.LoginLog{
		UserID:   userID,
		Username: strings.TrimSpace(username),
		IP:       strings.TrimSpace(ip),
		Device:   strings.TrimSpace(device),
		Status:   status,
	})
}
