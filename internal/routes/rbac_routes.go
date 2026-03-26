package routes

import (
	"github.com/gin-gonic/gin"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

// RBACRoutes 统一注册 RBAC 相关路由。
// 这样可以把鉴权/权限管理模块收敛到一个入口，便于后续维护与扩展。
func RBACRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 调用RoleRoutes完成当前处理。
	RoleRoutes(r, db, mgr)
	// 调用PermissionRoutes完成当前处理。
	PermissionRoutes(r, db, mgr)
	// 调用MenuRoutes完成当前处理。
	MenuRoutes(r, db, mgr)
	// 调用UserRoleRoutes完成当前处理。
	UserRoleRoutes(r, db, mgr)
	// 调用AuditRoutes完成当前处理。
	AuditRoutes(r, db, mgr)
	// 调用SystemMessageRoutes完成当前处理。
	SystemMessageRoutes(r, db, mgr)
	// 调用UserOpsRoutes完成当前处理。
	UserOpsRoutes(r, db, mgr)
	// 调用BizConfigRoutes完成当前处理。
	BizConfigRoutes(r, db, mgr)
}
