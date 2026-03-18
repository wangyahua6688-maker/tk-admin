package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/constants"
	rbac "go-admin-full/internal/controllers/rbac"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/middleware"
	rbacsvc "go-admin-full/internal/services/rbac"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

// MenuRoutes 处理MenuRoutes相关逻辑。
func MenuRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 定义并初始化当前变量。
	mc := rbac.NewMenuController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 菜单管理：菜单 CRUD、前端菜单树、菜单-权限绑定。
	mr := r.Group("/api/menus")
	// 调用mr.Use完成当前处理。
	mr.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用mr.GET完成当前处理。
		mr.GET("/", middleware.PermissionRequired(constants.PermMenuList, userRoleSvc), mc.List)
		// 调用mr.POST完成当前处理。
		mr.POST("/", middleware.PermissionRequired(constants.PermMenuCreate, userRoleSvc), mc.Create)
		// 调用mr.PUT完成当前处理。
		mr.PUT("/:id", middleware.PermissionRequired(constants.PermMenuUpdate, userRoleSvc), mc.Update)
		// 调用mr.GET完成当前处理。
		mr.GET("/frontend/tree", middleware.PermissionRequired(constants.PermMenuFrontendTree, userRoleSvc), mc.FrontendTree)
		// 调用mr.GET完成当前处理。
		mr.GET("/:id/permissions", middleware.PermissionRequired(constants.PermMenuPermissionView, userRoleSvc), mc.GetPermissions)
		// 调用mr.PUT完成当前处理。
		mr.PUT("/:id/permissions", middleware.PermissionRequired(constants.PermMenuPermissionBind, userRoleSvc), mc.BindPermissions)
		// 调用mr.GET完成当前处理。
		mr.GET("/:id", middleware.PermissionRequired(constants.PermMenuView, userRoleSvc), mc.Get)
		// 调用mr.DELETE完成当前处理。
		mr.DELETE("/:id", middleware.PermissionRequired(constants.PermMenuDelete, userRoleSvc), mc.Delete)
	}
}
