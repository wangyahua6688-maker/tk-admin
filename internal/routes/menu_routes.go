package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/services"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

func MenuRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	mc := controllers.NewMenuController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	// 菜单管理：菜单 CRUD、前端菜单树、菜单-权限绑定。
	mr := r.Group("/api/menus")
	mr.Use(middleware.NewJWTMiddleware(mgr))
	{
		mr.GET("/", middleware.PermissionRequired(constants.PermMenuList, userRoleSvc), mc.List)
		mr.POST("/", middleware.PermissionRequired(constants.PermMenuCreate, userRoleSvc), mc.Create)
		mr.PUT("/:id", middleware.PermissionRequired(constants.PermMenuUpdate, userRoleSvc), mc.Update)
		mr.GET("/frontend/tree", middleware.PermissionRequired(constants.PermMenuFrontendTree, userRoleSvc), mc.FrontendTree)
		mr.GET("/:id/permissions", middleware.PermissionRequired(constants.PermMenuPermissionView, userRoleSvc), mc.GetPermissions)
		mr.PUT("/:id/permissions", middleware.PermissionRequired(constants.PermMenuPermissionBind, userRoleSvc), mc.BindPermissions)
		mr.GET("/:id", middleware.PermissionRequired(constants.PermMenuView, userRoleSvc), mc.Get)
		mr.DELETE("/:id", middleware.PermissionRequired(constants.PermMenuDelete, userRoleSvc), mc.Delete)
	}
}
