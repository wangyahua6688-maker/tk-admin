package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/middleware"
	"go-admin-full/internal/services"
	"go-admin-full/internal/tokenpkg"
	"gorm.io/gorm"
)

func PermissionRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenpkg.Manager) {
	pc := controllers.NewPermissionController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	pr := r.Group("/api/permissions")
	pr.Use(middleware.NewJWTMiddleware(mgr))
	{
		pr.GET("/", middleware.PermissionRequired(constants.PermPermissionList, userRoleSvc), pc.List)
		pr.POST("/", middleware.PermissionRequired(constants.PermPermissionCreate, userRoleSvc), pc.Create)
		pr.PUT("/:id", middleware.PermissionRequired(constants.PermPermissionUpdate, userRoleSvc), pc.Update)
		pr.GET("/:id", middleware.PermissionRequired(constants.PermPermissionView, userRoleSvc), pc.Get)
		pr.DELETE("/:id", middleware.PermissionRequired(constants.PermPermissionDelete, userRoleSvc), pc.Delete)
	}
}
