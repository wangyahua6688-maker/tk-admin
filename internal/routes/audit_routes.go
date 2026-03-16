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

// AuditRoutes 审计相关路由。
func AuditRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	ac := rbac.NewAuditController(db)
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	ag := r.Group("/api/audit")
	ag.Use(middleware.NewJWTMiddleware(mgr))
	{
		ag.GET("/login-logs", middleware.PermissionRequired(constants.PermAuditLoginLogList, userRoleSvc), ac.ListLoginLogs)
	}
}
