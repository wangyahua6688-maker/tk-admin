package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/constants"
	biz "go-admin/internal/controllers/biz"
	rbacdao "go-admin/internal/dao/rbac"
	"go-admin/internal/middleware"
	rbacsvc "go-admin/internal/services/rbac"
	tokenjwt "go-admin/internal/token/jwt"
	"gorm.io/gorm"
)

// UserOpsRoutes 处理UserOpsRoutes相关逻辑。
func UserOpsRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	clientUserCtrl := biz.NewClientUserController(db)
	postArticleCtrl := biz.NewPostArticleController(db)
	postCommentCtrl := biz.NewPostCommentController(db)
	hotCommentCtrl := biz.NewHotCommentController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	group := r.Group("/api/user-mgmt")
	// 调用group.Use完成当前处理。
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用group.GET完成当前处理。
		group.GET("/client-users", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc, mgr), clientUserCtrl.ListClientUsers)
		// 调用group.POST完成当前处理。
		group.POST("/client-users", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc, mgr), clientUserCtrl.CreateClientUser)
		// 调用group.PUT完成当前处理。
		group.PUT("/client-users/:id", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc, mgr), clientUserCtrl.UpdateClientUser)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/client-users/:id", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc, mgr), clientUserCtrl.DeleteClientUser)

		// 调用group.GET完成当前处理。
		group.GET("/post-articles", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc, mgr), postArticleCtrl.ListPostArticles)
		// 调用group.POST完成当前处理。
		group.POST("/post-articles", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc, mgr), postArticleCtrl.CreatePostArticle)
		// 调用group.PUT完成当前处理。
		group.PUT("/post-articles/:id", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc, mgr), postArticleCtrl.UpdatePostArticle)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/post-articles/:id", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc, mgr), postArticleCtrl.DeletePostArticle)
		// 调用group.GET完成当前处理。
		group.GET("/post-articles/:id/comments", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc, mgr), postCommentCtrl.ListPostComments)
		// 调用group.POST完成当前处理。
		group.POST("/post-articles/:id/comments", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc, mgr), postCommentCtrl.CreatePostComment)

		// 调用group.GET完成当前处理。
		group.GET("/hot-comments", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc, mgr), hotCommentCtrl.ListHotComments)
		// 调用group.POST完成当前处理。
		group.POST("/hot-comments", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc, mgr), hotCommentCtrl.CreateHotComment)
		// 调用group.PUT完成当前处理。
		group.PUT("/hot-comments/:id", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc, mgr), hotCommentCtrl.UpdateHotComment)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/hot-comments/:id", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc, mgr), hotCommentCtrl.DeleteHotComment)
	}
}
