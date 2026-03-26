package rbac

import (
	"fmt"
	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin-full/internal/constants"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/models"
	rbacsvc "go-admin-full/internal/services/rbac"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 依赖注入
func NewUserController(db *gorm.DB) *UserController {
	userDao := rbacdao.NewUserDao(db)                                          // 初始化用户 DAO
	userService := rbacsvc.NewUserService(userDao)                             // 初始化用户服务
	msgSvc := rbacsvc.NewSystemMessageService(rbacdao.NewSystemMessageDao(db)) // 初始化系统消息服务
	// 返回当前处理结果。
	return &UserController{service: userService, msgSvc: msgSvc}
}

// UserController 定义UserController相关结构。
type UserController struct {
	service *rbacsvc.UserService          // 用户业务服务
	msgSvc  *rbacsvc.SystemMessageService // 系统消息服务
}

// CreateUserReq 定义CreateUserReq相关结构。
type CreateUserReq struct {
	// 处理当前语句逻辑。
	Username string `json:"username" binding:"required"`
	// 更新当前变量或字段值。
	Password string `json:"password" binding:"required,min=8"`
	// 处理当前语句逻辑。
	Email string `json:"email"`
	// 处理当前语句逻辑。
	Avatar string `json:"avatar"`
	// 处理当前语句逻辑。
	Status *int `json:"status"`
}

// UpdateUserReq 定义UpdateUserReq相关结构。
type UpdateUserReq struct {
	// 处理当前语句逻辑。
	Email string `json:"email"`
	// 处理当前语句逻辑。
	Avatar *string `json:"avatar"`
	// 处理当前语句逻辑。
	Password string `json:"password"`
	// 处理当前语句逻辑。
	Status *int `json:"status"`
}

// UserResp 定义UserResp相关结构。
type UserResp struct {
	// 处理当前语句逻辑。
	ID uint `json:"id"`
	// 处理当前语句逻辑。
	Username string `json:"username"`
	// 处理当前语句逻辑。
	Email string `json:"email"`
	// 处理当前语句逻辑。
	Avatar string `json:"avatar"`
	// 处理当前语句逻辑。
	Status int `json:"status"`
	// 处理当前语句逻辑。
	CreatedAt time.Time `json:"created_at"`
	// 处理当前语句逻辑。
	UpdatedAt time.Time `json:"updated_at"`
}

// List 用户列表。
func (c *UserController) List(ctx *gin.Context) {
	// 拉取全部用户列表
	users, err := c.service.ListAllUsers(ctx.Request.Context())
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}

	// 转换为响应结构
	resp := make([]UserResp, 0, len(users))
	// 循环处理当前数据集合。
	for _, u := range users {
		// 更新当前变量或字段值。
		resp = append(resp, toUserResp(u))
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, resp)
}

// Get 查询用户详情。
func (c *UserController) Get(ctx *gin.Context) {
	// 解析并校验用户 ID
	id, err := strconv.Atoi(ctx.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "invalid user id")
		// 返回当前处理结果。
		return
	}

	// 查询用户详情
	user, err := c.service.GetUserByID(ctx.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizResourceNotFound, "用户不存在")
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, toUserResp(*user))
}

// Create 创建用户。
func (c *UserController) Create(ctx *gin.Context) {
	// 声明当前变量。
	var req CreateUserReq
	// 判断条件并进入对应分支逻辑。
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "参数错误")
		// 返回当前处理结果。
		return
	}

	// 更新当前变量或字段值。
	req.Username = strings.TrimSpace(req.Username)
	// 更新当前变量或字段值。
	req.Email = strings.TrimSpace(req.Email)
	// 更新当前变量或字段值。
	req.Avatar = strings.TrimSpace(req.Avatar)

	// 定义并初始化当前变量。
	user, err := c.service.CreateUser(ctx.Request.Context(), req.Username, req.Password, req.Email, req.Avatar)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, err.Error())
		// 返回当前处理结果。
		return
	}

	// 可选设置状态
	if req.Status != nil {
		// 判断条件并进入对应分支逻辑。
		if err := c.service.UpdateUser(ctx.Request.Context(), user.ID, req.Email, req.Status, "", nil); err != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(ctx, constants.AdminSysInternalError, err.Error())
			// 返回当前处理结果。
			return
		}
	}

	// 定义并初始化当前变量。
	created, err := c.service.GetUserByID(ctx.Request.Context(), user.ID)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}

	// 创建用户后向目标用户投递系统消息，便于其登录后感知账号变更。
	_ = c.msgSvc.PushToUser(
		// 调用ctx.Request.Context完成当前处理。
		ctx.Request.Context(),
		// 处理当前语句逻辑。
		created.ID,
		// 处理当前语句逻辑。
		"账号创建通知",
		// 调用fmt.Sprintf完成当前处理。
		fmt.Sprintf("管理员已创建你的账号（用户名：%s），请及时确认角色与权限配置。", created.Username),
		// 处理当前语句逻辑。
		"success",
		// 处理当前语句逻辑。
		"user",
		// 处理当前语句逻辑。
		created.ID,
		// 调用ctx.GetUint完成当前处理。
		ctx.GetUint("uid"),
	)

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, toUserResp(*created))
}

// Update 更新用户。
func (c *UserController) Update(ctx *gin.Context) {
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(ctx.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "invalid user id")
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var req UpdateUserReq
	// 判断条件并进入对应分支逻辑。
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "参数错误")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	target, err := c.service.GetUserByID(ctx.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizResourceNotFound, "用户不存在")
		// 返回当前处理结果。
		return
	}

	// 判断条件并进入对应分支逻辑。
	if err := c.service.UpdateUser(ctx.Request.Context(), uint(id), req.Email, req.Status, req.Password, req.Avatar); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, err.Error())
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	updated, err := c.service.GetUserByID(ctx.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}

	// 更新当前变量或字段值。
	_ = c.msgSvc.PushToUser(
		// 调用ctx.Request.Context完成当前处理。
		ctx.Request.Context(),
		// 处理当前语句逻辑。
		updated.ID,
		// 处理当前语句逻辑。
		"账号资料更新通知",
		// 调用fmt.Sprintf完成当前处理。
		fmt.Sprintf("管理员已更新你的账号资料（用户名：%s）。如有异常请联系系统管理员。", target.Username),
		// 处理当前语句逻辑。
		"info",
		// 处理当前语句逻辑。
		"user",
		// 处理当前语句逻辑。
		updated.ID,
		// 调用ctx.GetUint完成当前处理。
		ctx.GetUint("uid"),
	)

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, toUserResp(*updated))
}

// Delete 删除用户。
func (c *UserController) Delete(ctx *gin.Context) {
	// 定义并初始化当前变量。
	id, err := strconv.Atoi(ctx.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizInvalidRequest, "invalid user id")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	target, err := c.service.GetUserByID(ctx.Request.Context(), uint(id))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizResourceNotFound, "用户不存在")
		// 返回当前处理结果。
		return
	}

	// 安全防护：禁止删除内置管理员账号
	if strings.EqualFold(target.Username, "admin") {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthForbidden, "admin账号不可删除")
		// 返回当前处理结果。
		return
	}

	// 安全防护：禁止删除当前登录用户
	uid := ctx.GetUint("uid")
	// 判断条件并进入对应分支逻辑。
	if uid == uint(id) {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthForbidden, "不可删除当前登录账号")
		// 返回当前处理结果。
		return
	}

	// 判断条件并进入对应分支逻辑。
	if err := c.service.DeleteUser(ctx.Request.Context(), uint(id)); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, gin.H{"msg": "deleted"})
}

// Profile 处理Profile相关逻辑。
func (c *UserController) Profile(ctx *gin.Context) {
	// 定义并初始化当前变量。
	uid := ctx.GetUint("uid")
	// 判断条件并进入对应分支逻辑。
	if uid == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminAuthUnauthorized, "用户未认证")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	user, err := c.service.GetUserByID(ctx.Request.Context(), uid)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(ctx, constants.AdminBizResourceNotFound, "用户不存在")
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(ctx, toUserResp(*user))
}

// toUserResp 处理toUserResp相关逻辑。
func toUserResp(u models.User) UserResp {
	// 返回当前处理结果。
	return UserResp{
		// 处理当前语句逻辑。
		ID: u.ID,
		// 处理当前语句逻辑。
		Username: u.Username,
		// 处理当前语句逻辑。
		Email: u.Email,
		// 处理当前语句逻辑。
		Avatar: u.Avatar,
		// 处理当前语句逻辑。
		Status: u.Status,
		// 处理当前语句逻辑。
		CreatedAt: u.CreatedAt,
		// 处理当前语句逻辑。
		UpdatedAt: u.UpdatedAt,
	}
}
