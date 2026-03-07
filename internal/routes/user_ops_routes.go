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

func UserOpsRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	ctrl := controllers.NewUserOpsController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	group := r.Group("/api/user-mgmt")
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		group.GET("/client-users", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.ListClientUsers)
		group.POST("/client-users", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.CreateClientUser)
		group.PUT("/client-users/:id", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.UpdateClientUser)
		group.DELETE("/client-users/:id", middleware.PermissionRequired(constants.PermClientUserList, userRoleSvc), ctrl.DeleteClientUser)

		group.GET("/post-articles", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.ListPostArticles)
		group.POST("/post-articles", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.CreatePostArticle)
		group.PUT("/post-articles/:id", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.UpdatePostArticle)
		group.DELETE("/post-articles/:id", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.DeletePostArticle)
		group.GET("/post-articles/:id/comments", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.ListPostComments)
		group.POST("/post-articles/:id/comments", middleware.PermissionRequired(constants.PermClientPostList, userRoleSvc), ctrl.CreatePostComment)

		group.GET("/hot-comments", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.ListHotComments)
		group.POST("/hot-comments", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.CreateHotComment)
		group.PUT("/hot-comments/:id", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.UpdateHotComment)
		group.DELETE("/hot-comments/:id", middleware.PermissionRequired(constants.PermClientHotCommentList, userRoleSvc), ctrl.DeleteHotComment)
	}
}
