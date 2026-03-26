package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/constants"
	rbac "go-admin/internal/controllers/rbac"
	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/middleware"
	rbacsvc "go-admin/internal/services/rbac"
	tokenjwt "go-admin/internal/token/jwt"
	"gorm.io/gorm"
)

// AuditRoutes 审计相关路由。
func AuditRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 定义并初始化当前变量。
	ac := rbac.NewAuditController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	ag := r.Group("/api/audit")
	// 调用ag.Use完成当前处理。
	ag.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用ag.GET完成当前处理。
		ag.GET("/login-logs", middleware.PermissionRequired(constants.PermAuditLoginLogList, userRoleSvc, mgr), ac.ListLoginLogs)
	}
}
