package rbac

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"go-admin-full/config"
	"go-admin-full/internal/auth/sessioncookie"
	"go-admin-full/internal/constants"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	rbacsvc "go-admin-full/internal/services/rbac"
	tokenjwt "go-admin-full/internal/token/jwt"
	commonresp "tk-common/utils/httpresp"

	commonlogx "tk-common/utils/logx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthController 定义AuthController相关结构。
type AuthController struct {
	authService *rbacsvc.AuthService     // 认证业务服务
	loginLogSvc *rbacsvc.LoginLogService // 登录日志服务
	tokenMgr    *tokenjwt.Manager        // JWT 管理器
	cookieOpt   sessioncookie.Options    // 认证 Cookie 选项
}

// 声明当前常量。
const (
	// 登录失败限制：15 分钟最多 5 次
	loginFailMax = 5
	// 更新当前变量或字段值。
	loginFailWindow = 15 * time.Minute
)

// 修改构造函数，同时传入db和token管理器
func NewAuthController(db *gorm.DB, mgr *tokenjwt.Manager, cfg config.Config) *AuthController {
	authDao := rbacdao.NewAuthDao(db)                                     // 初始化认证 DAO
	authService := rbacsvc.NewAuthService(authDao)                        // 初始化认证服务
	loginLogSvc := rbacsvc.NewLoginLogService(rbacdao.NewLoginLogDao(db)) // 初始化登录日志服务
	// 返回当前处理结果。
	return &AuthController{
		authService: authService,                   // 注入认证服务
		loginLogSvc: loginLogSvc,                   // 注入登录日志服务
		tokenMgr:    mgr,                           // 注入 Token 管理器
		cookieOpt:   sessioncookie.FromConfig(cfg), // 注入认证 Cookie 配置
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Username string `json:"username" binding:"required"`
		// 处理当前语句逻辑。
		Password string `json:"password" binding:"required"`
	}

	// 绑定并校验请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "参数错误："+err.Error())
		// 返回当前处理结果。
		return
	}

	// 获取上下文中的logger
	logger := commonlogx.LoggerFromContext(ctx.Request.Context())
	req.Username = strings.TrimSpace(req.Username) // 规范化用户名
	clientIP := ctx.ClientIP()                     // 取客户端 IP
	deviceID := resolveDeviceID(ctx)               // 解析设备标识
	// 调用logger.Info完成当前处理。
	logger.Info("用户登录尝试: %s", req.Username)

	// 安全防护：登录失败次数限制
	if blocked := c.isLoginBlocked(req.Username, clientIP); blocked {
		c.writeLoginLog(ctx, 0, req.Username, clientIP, deviceID, 0) // 记录失败日志
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthRateLimited, "登录失败次数过多，请稍后重试")
		// 返回当前处理结果。
		return
	}

	// 调用服务层登录
	user, err := c.authService.Login(ctx.Request.Context(), req.Username, req.Password)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		c.recordLoginFailure(req.Username, clientIP)                 // 记录失败次数
		c.writeLoginLog(ctx, 0, req.Username, clientIP, deviceID, 0) // 记录失败日志
		// 调用logger.Error完成当前处理。
		logger.Error("登录失败: %v", err)
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthUnauthorized, err.Error())
		// 返回当前处理结果。
		return
	}
	c.clearLoginFailure(req.Username, clientIP) // 清空失败计数

	// 生成JWT Token
	accessToken, refreshToken, err := c.tokenMgr.GenerateTokensWithDevice(user.ID, deviceID)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用logger.Error完成当前处理。
		logger.Error("生成Token失败: %v", err)
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminSysInternalError, "生成Token失败")
		// 返回当前处理结果。
		return
	}
	c.writeLoginLog(ctx, user.ID, user.Username, clientIP, deviceID, 1) // 记录成功日志
	c.setAuthCookies(ctx, accessToken, refreshToken)                    // 设置 HttpOnly 认证 Cookie

	// 调用logger.Info完成当前处理。
	logger.Info("用户登录成功: %s (ID: %d)", req.Username, user.ID)

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, gin.H{
		// 处理当前语句逻辑。
		"access_token": accessToken,
		// 处理当前语句逻辑。
		"refresh_token": refreshToken,
		// 处理当前语句逻辑。
		"user_id": user.ID,
		// 处理当前语句逻辑。
		"username": user.Username,
		// 调用int完成当前处理。
		"expires_in": int(c.tokenMgr.Config.AccessExpire.Seconds()),
	})
}

// Refresh 刷新Token
func (c *AuthController) Refresh(ctx *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		RefreshToken string `json:"refresh_token"`
	}

	// 绑定并校验请求体
	if err := ctx.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "参数错误："+err.Error())
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	logger := commonlogx.LoggerFromContext(ctx.Request.Context())
	// 调用logger.Info完成当前处理。
	logger.Info("刷新Token请求")

	// 安全防护：刷新前先校验 refresh token 基础有效性和用户状态。
	req.RefreshToken = strings.TrimSpace(req.RefreshToken) // 规范化 token
	if req.RefreshToken == "" {
		req.RefreshToken = c.readRefreshTokenFromCookie(ctx) // 从 Cookie 回退读取 refresh token
	}
	if req.RefreshToken == "" {
		c.clearAuthCookies(ctx) // 缺少 refresh token 时清理残留 Cookie
		commonresp.GinError(ctx, constants.AdminAuthUnauthorized, constants.ErrInvalidToken.Error())
		return
	}
	claims, err := tokenjwt.ParseTokenClaims(req.RefreshToken, c.tokenMgr.Config.SigningKey) // 解析 claims
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		c.clearAuthCookies(ctx) // 非法/过期 refresh token 直接清空 Cookie
		// 定义并初始化当前变量。
		code := constants.AdminAuthTokenInvalid
		// 判断条件并进入对应分支逻辑。
		if err == constants.ErrExpiredToken {
			// 更新当前变量或字段值。
			code = constants.AdminAuthForbidden
		}
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, code, err.Error())
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if claims.TokenType != tokenjwt.TokenTypeRefresh {
		c.clearAuthCookies(ctx) // token 类型错误时清理 Cookie
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthUnauthorized, constants.ErrInvalidToken.Error())
		// 返回当前处理结果。
		return
	}
	// 校验用户状态, 封禁用户禁止刷新
	user, err := c.authService.GetUserByID(ctx.Request.Context(), claims.UserID) // 查询用户状态
	// 判断条件并进入对应分支逻辑。
	if err != nil || user.Status != 1 {
		c.clearAuthCookies(ctx) // 账号不可用时清理 Cookie
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthUnauthorized, "账号不可用")
		// 返回当前处理结果。
		return
	}

	// 使用token管理器刷新token
	newAccessToken, newRefreshToken, err := c.tokenMgr.RefreshTokenPair(req.RefreshToken) // 生成新 token
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		c.clearAuthCookies(ctx) // 刷新失败时清理 Cookie，避免前端重复携带失效凭证
		// 调用logger.Error完成当前处理。
		logger.Error("刷新Token失败: %v", err)
		// 定义并初始化当前变量。
		code := constants.AdminAuthTokenInvalid
		// 判断条件并进入对应分支逻辑。
		if err == constants.ErrExpiredToken {
			// 更新当前变量或字段值。
			code = constants.AdminAuthForbidden
		}
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, code, err.Error())
		// 返回当前处理结果。
		return
	}

	// 调用logger.Info完成当前处理。
	logger.Info("Token刷新成功")
	c.setAuthCookies(ctx, newAccessToken, newRefreshToken) // 回写轮换后的认证 Cookie

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, gin.H{
		"access_token":  newAccessToken,                                // 新 access token
		"refresh_token": newRefreshToken,                               // 新 refresh token
		"expires_in":    int(c.tokenMgr.Config.AccessExpire.Seconds()), // access token 有效期
	})
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Username string `json:"username" binding:"required"`
		// 更新当前变量或字段值。
		Password string `json:"password" binding:"required,min=8"`
		// 处理当前语句逻辑。
		Email string `json:"email"`
	}

	// 绑定并校验请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "参数错误："+err.Error())
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	logger := commonlogx.LoggerFromContext(ctx.Request.Context())
	req.Username = strings.TrimSpace(req.Username) // 规范化用户名
	req.Email = strings.TrimSpace(req.Email)       // 规范化邮箱
	// 调用logger.Info完成当前处理。
	logger.Info("用户注册尝试: %s", req.Username)

	// 调用注册服务
	if err := c.authService.Register(ctx.Request.Context(), req.Username, req.Password, req.Email); err != nil {
		// 调用logger.Error完成当前处理。
		logger.Error("注册失败: %v", err)
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, err.Error())
		// 返回当前处理结果。
		return
	}

	// 调用logger.Info完成当前处理。
	logger.Info("用户注册成功: %s", req.Username)

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, gin.H{"message": "注册成功"})
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	c.clearAuthCookies(ctx) // 无论鉴权结果如何，优先清理客户端 Cookie
	// 从JWT中间件中获取用户ID
	uid, exists := ctx.Get("uid")
	// 判断条件并进入对应分支逻辑。
	if !exists {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthUnauthorized, "用户未认证")
		// 返回当前处理结果。
		return
	}

	userID := uid.(uint) // 转换为 uint
	// 定义并初始化当前变量。
	logger := commonlogx.LoggerFromContext(ctx.Request.Context())
	// 更新当前变量或字段值。
	logger.Info("用户登出: ID=%d", userID)

	deviceID := resolveDeviceID(ctx) // 解析设备标识
	// 判断条件并进入对应分支逻辑。
	if claimsRaw, ok := ctx.Get("claims"); ok {
		// 判断条件并进入对应分支逻辑。
		if claims, ok := claimsRaw.(*tokenjwt.Claims); ok && strings.TrimSpace(claims.DeviceID) != "" {
			// 更新当前变量或字段值。
			deviceID = claims.DeviceID
		}
	}

	// 使 refresh token 无效（当前设备）
	if err := c.tokenMgr.InvalidateRefresh(userID, deviceID); err != nil {
		// 调用logger.Error完成当前处理。
		logger.Error("登出失败，无法使 refresh token 无效: %v", err)
	}

	// 将当前 access token 放入黑名单
	accessToken, _ := ctx.Get("access_token")
	// 判断条件并进入对应分支逻辑。
	if tokenStr, ok := accessToken.(string); ok && tokenStr != "" {
		// 判断条件并进入对应分支逻辑。
		if err := c.tokenMgr.RevokeAccessToken(tokenStr); err != nil {
			// 调用logger.Error完成当前处理。
			logger.Error("登出失败，无法撤销 access token: %v", err)
		}
	}

	// 更新当前变量或字段值。
	logger.Info("用户登出成功: ID=%d", userID)

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, gin.H{"message": "登出成功"})
}

// setAuthCookies 写入认证 Cookie（HttpOnly + SameSite）。
func (c *AuthController) setAuthCookies(ctx *gin.Context, accessToken, refreshToken string) {
	ctx.SetSameSite(c.cookieOpt.SameSite) // 应用 SameSite 策略
	ctx.SetCookie(
		c.cookieOpt.AccessTokenName,
		strings.TrimSpace(accessToken),
		int(c.tokenMgr.Config.AccessExpire.Seconds()),
		c.cookieOpt.Path,
		c.cookieOpt.Domain,
		c.cookieOpt.Secure,
		c.cookieOpt.HTTPOnly,
	)
	ctx.SetCookie(
		c.cookieOpt.RefreshTokenName,
		strings.TrimSpace(refreshToken),
		int(c.tokenMgr.Config.RefreshExpire.Seconds()),
		c.cookieOpt.Path,
		c.cookieOpt.Domain,
		c.cookieOpt.Secure,
		c.cookieOpt.HTTPOnly,
	)
}

// clearAuthCookies 清理认证 Cookie。
func (c *AuthController) clearAuthCookies(ctx *gin.Context) {
	ctx.SetSameSite(c.cookieOpt.SameSite) // 清理时保持同一 SameSite 策略
	ctx.SetCookie(c.cookieOpt.AccessTokenName, "", -1, c.cookieOpt.Path, c.cookieOpt.Domain, c.cookieOpt.Secure, c.cookieOpt.HTTPOnly)
	ctx.SetCookie(c.cookieOpt.RefreshTokenName, "", -1, c.cookieOpt.Path, c.cookieOpt.Domain, c.cookieOpt.Secure, c.cookieOpt.HTTPOnly)
}

// readRefreshTokenFromCookie 从认证 Cookie 读取 refresh token。
func (c *AuthController) readRefreshTokenFromCookie(ctx *gin.Context) string {
	raw, err := ctx.Cookie(c.cookieOpt.RefreshTokenName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(raw)
}

// resolveDeviceID 处理resolveDeviceID相关逻辑。
func resolveDeviceID(ctx *gin.Context) string {
	// 定义并初始化当前变量。
	deviceID := strings.TrimSpace(ctx.GetHeader("X-Device-ID"))
	// 判断条件并进入对应分支逻辑。
	if deviceID != "" {
		// 返回当前处理结果。
		return deviceID
	}
	// 定义并初始化当前变量。
	ua := strings.TrimSpace(ctx.Request.UserAgent())
	// 判断条件并进入对应分支逻辑。
	if ua != "" {
		// 返回当前处理结果。
		return ua
	}
	return "default"
}

// loginFailKey 处理loginFailKey相关逻辑。
func (c *AuthController) loginFailKey(username, ip string) string {
	// 返回当前处理结果。
	return "auth:login_fail:" + strings.ToLower(strings.TrimSpace(username)) + ":" + strings.TrimSpace(ip)
}

// isLoginBlocked 判断是否触发登录失败限制。
func (c *AuthController) isLoginBlocked(username, ip string) bool {
	// 判断条件并进入对应分支逻辑。
	if c.tokenMgr == nil || c.tokenMgr.Store == nil {
		// 返回当前处理结果。
		return false
	}

	// 定义并初始化当前变量。
	key := c.loginFailKey(username, ip)
	// 定义并初始化当前变量。
	v, err := c.tokenMgr.Store.Get(key)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return false
	}

	// 定义并初始化当前变量。
	count, _ := strconv.Atoi(v)
	// 返回当前处理结果。
	return count >= loginFailMax
}

// recordLoginFailure 记录登录失败次数。
func (c *AuthController) recordLoginFailure(username, ip string) {
	// 判断条件并进入对应分支逻辑。
	if c.tokenMgr == nil || c.tokenMgr.Store == nil {
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	key := c.loginFailKey(username, ip)
	// 定义并初始化当前变量。
	current := 0
	// 判断条件并进入对应分支逻辑。
	if v, err := c.tokenMgr.Store.Get(key); err == nil {
		// 更新当前变量或字段值。
		current, _ = strconv.Atoi(v)
	}
	// 更新当前变量或字段值。
	_ = c.tokenMgr.Store.Set(key, strconv.Itoa(current+1), loginFailWindow)
}

// clearLoginFailure 清空登录失败计数。
func (c *AuthController) clearLoginFailure(username, ip string) {
	// 判断条件并进入对应分支逻辑。
	if c.tokenMgr == nil || c.tokenMgr.Store == nil {
		// 返回当前处理结果。
		return
	}
	// 更新当前变量或字段值。
	_ = c.tokenMgr.Store.Delete(c.loginFailKey(username, ip))
}

// writeLoginLog 写入登录审计日志，失败不影响主流程。
func (c *AuthController) writeLoginLog(ctx *gin.Context, userID uint, username, ip, device string, status int) {
	// 判断条件并进入对应分支逻辑。
	if c.loginLogSvc == nil {
		// 返回当前处理结果。
		return
	}
	// 更新当前变量或字段值。
	_ = c.loginLogSvc.CreateLoginLog(ctx.Request.Context(), &models.LoginLog{
		// 处理当前语句逻辑。
		UserID: userID,
		// 调用strings.TrimSpace完成当前处理。
		Username: strings.TrimSpace(username),
		// 调用strings.TrimSpace完成当前处理。
		IP: strings.TrimSpace(ip),
		// 调用strings.TrimSpace完成当前处理。
		Device: strings.TrimSpace(device),
		// 处理当前语句逻辑。
		Status: status,
	})
}
