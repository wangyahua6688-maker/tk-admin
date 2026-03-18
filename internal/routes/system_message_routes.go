package routes

import (
	"github.com/gin-gonic/gin"
	rbac "go-admin-full/internal/controllers/rbac"
	"go-admin-full/internal/middleware"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

// SystemMessageRoutes 系统消息路由。
func SystemMessageRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 定义并初始化当前变量。
	mc := rbac.NewSystemMessageController(db)

	// 定义并初始化当前变量。
	group := r.Group("/api/messages")
	// 调用group.Use完成当前处理。
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 系统消息属于“登录用户私有资源”，不依赖角色权限码。
		// 这样可保证：当管理员调整用户/角色/权限后，被影响用户一定能读取通知结果。
		group.GET("/", mc.ListMyMessages)
		// 调用group.POST完成当前处理。
		group.POST("/:id/read", mc.MarkRead)
		// 调用group.POST完成当前处理。
		group.POST("/read-all", mc.MarkAllRead)
	}
}
