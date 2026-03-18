package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/constants"
	biz "go-admin-full/internal/controllers/biz"
	rbacdao "go-admin-full/internal/dao/rbac"
	"go-admin-full/internal/middleware"
	rbacsvc "go-admin-full/internal/services/rbac"
	tokenjwt "go-admin-full/internal/token/jwt"
	"gorm.io/gorm"
)

// UserOpsRoutes 处理UserOpsRoutes相关逻辑。
func UserOpsRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 定义并初始化当前变量。
	ctrl := biz.NewUserOpsController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	group := r.Group("/api/user-mgmt")
	// 调用group.Use完成当前处理。
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用group.GET完成当前处理。
		group.GET("/client-users", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.ListClientUsers)
		// 调用group.POST完成当前处理。
		group.POST("/client-users", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.CreateClientUser)
		// 调用group.PUT完成当前处理。
		group.PUT("/client-users/:id", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.UpdateClientUser)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/client-users/:id", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.DeleteClientUser)

		// 调用group.GET完成当前处理。
		group.GET("/post-articles", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.ListPostArticles)
		// 调用group.POST完成当前处理。
		group.POST("/post-articles", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.CreatePostArticle)
		// 调用group.PUT完成当前处理。
		group.PUT("/post-articles/:id", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.UpdatePostArticle)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/post-articles/:id", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.DeletePostArticle)
		// 调用group.GET完成当前处理。
		group.GET("/post-articles/:id/comments", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.ListPostComments)
		// 调用group.POST完成当前处理。
		group.POST("/post-articles/:id/comments", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.CreatePostComment)

		// 调用group.GET完成当前处理。
		group.GET("/hot-comments", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.ListHotComments)
		// 调用group.POST完成当前处理。
		group.POST("/hot-comments", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.CreateHotComment)
		// 调用group.PUT完成当前处理。
		group.PUT("/hot-comments/:id", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.UpdateHotComment)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/hot-comments/:id", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.DeleteHotComment)
	}
}
