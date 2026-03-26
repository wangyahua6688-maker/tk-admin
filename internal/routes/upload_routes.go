package routes

import (
	"github.com/gin-gonic/gin"
	biz "go-admin/internal/controllers/biz"
	"go-admin/internal/middleware"
	tokenjwt "go-admin/internal/token/jwt"
)

// UploadRoutes 处理UploadRoutes相关逻辑。
func UploadRoutes(r *gin.Engine, mgr *tokenjwt.Manager) {
	// 定义并初始化当前变量。
	ctrl := biz.NewUploadController()

	// 定义并初始化当前变量。
	group := r.Group("/api")
	// 调用group.Use完成当前处理。
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用group.POST完成当前处理。
		group.POST("/upload/image", ctrl.UploadImage)
	}
}
