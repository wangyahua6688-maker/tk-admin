package controllers

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/dao"
	service "go-admin-full/internal/services"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
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
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONError(ctx, 400, "参数错误")
		return
	}

	user, err := c.service.CreateUser(ctx.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		utils.JSONError(ctx, 400, err.Error())
		return
	}
	utils.JSONOK(ctx, gin.H{"user_id": user.ID})
}

func (c *UserController) ListUsers(ctx *gin.Context) {
	users, err := c.service.ListAllUsers(ctx.Request.Context())
	if err != nil {
		utils.JSONError(ctx, 500, err.Error())
		return
	}
	utils.JSONOK(ctx, users)
}

func (c *UserController) Profile(ctx *gin.Context) {

	return
}
