package controllers

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/models"
	service "go-admin-full/internal/services"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 依赖注入
func NewUserController(db *gorm.DB) *UserController {
	userDao := dao.NewUserDao(db)
	userService := service.NewUserService(userDao)
	return &UserController{service: userService}
}

type UserController struct {
	service *service.UserService
}

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email"`
	Status   *int   `json:"status"`
}

type UpdateUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   *int   `json:"status"`
}

type UserResp struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// List 用户列表。
func (c *UserController) List(ctx *gin.Context) {
	users, err := c.service.ListAllUsers(ctx.Request.Context())
	if err != nil {
		utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]UserResp, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResp(u))
	}
	utils.JSONOK(ctx, resp)
}

// Get 查询用户详情。
func (c *UserController) Get(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		utils.JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := c.service.GetUserByID(ctx.Request.Context(), uint(id))
	if err != nil {
		utils.JSONError(ctx, http.StatusNotFound, "用户不存在")
		return
	}
	utils.JSONOK(ctx, toUserResp(*user))
}

// Create 创建用户。
func (c *UserController) Create(ctx *gin.Context) {
	var req CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	user, err := c.service.CreateUser(ctx.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 可选设置状态
	if req.Status != nil {
		if err := c.service.UpdateUser(ctx.Request.Context(), user.ID, req.Email, req.Status, ""); err != nil {
			utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	created, err := c.service.GetUserByID(ctx.Request.Context(), user.ID)
	if err != nil {
		utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(ctx, toUserResp(*created))
}

// Update 更新用户。
func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		utils.JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	if err := c.service.UpdateUser(ctx.Request.Context(), uint(id), req.Email, req.Status, req.Password); err != nil {
		utils.JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := c.service.GetUserByID(ctx.Request.Context(), uint(id))
	if err != nil {
		utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(ctx, toUserResp(*updated))
}

// Delete 删除用户。
func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		utils.JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	target, err := c.service.GetUserByID(ctx.Request.Context(), uint(id))
	if err != nil {
		utils.JSONError(ctx, http.StatusNotFound, "用户不存在")
		return
	}

	// 安全防护：禁止删除内置管理员账号
	if strings.EqualFold(target.Username, "admin") {
		utils.JSONError(ctx, http.StatusForbidden, "admin账号不可删除")
		return
	}

	// 安全防护：禁止删除当前登录用户
	uid := ctx.GetUint("uid")
	if uid == uint(id) {
		utils.JSONError(ctx, http.StatusForbidden, "不可删除当前登录账号")
		return
	}

	if err := c.service.DeleteUser(ctx.Request.Context(), uint(id)); err != nil {
		utils.JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(ctx, gin.H{"msg": "deleted"})
}

func (c *UserController) Profile(ctx *gin.Context) {
	uid := ctx.GetUint("uid")
	if uid == 0 {
		utils.JSONError(ctx, http.StatusUnauthorized, "用户未认证")
		return
	}

	user, err := c.service.GetUserByID(ctx.Request.Context(), uid)
	if err != nil {
		utils.JSONError(ctx, http.StatusNotFound, "用户不存在")
		return
	}
	utils.JSONOK(ctx, toUserResp(*user))
}

func toUserResp(u models.User) UserResp {
	return UserResp{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
