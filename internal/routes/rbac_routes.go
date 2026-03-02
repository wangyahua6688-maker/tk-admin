package routes

import (
	"github.com/gin-gonic/gin"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

// RBACRoutes 统一注册 RBAC 相关路由。
// 这样可以把鉴权/权限管理模块收敛到一个入口，便于后续维护与扩展。
func RBACRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	RoleRoutes(r, db, mgr)
	PermissionRoutes(r, db, mgr)
	MenuRoutes(r, db, mgr)
	UserRoleRoutes(r, db, mgr)
	AuditRoutes(r, db, mgr)
}
