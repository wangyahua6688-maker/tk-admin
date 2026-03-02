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

// AuditRoutes 审计相关路由。
func AuditRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	ac := controllers.NewAuditController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	ag := r.Group("/api/audit")
	ag.Use(middleware.NewJWTMiddleware(mgr))
	{
		ag.GET("/login-logs", middleware.PermissionRequired(constants.PermAuditLoginLogList, userRoleSvc), ac.ListLoginLogs)
	}
}
