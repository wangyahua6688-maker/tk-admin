package routes

import (
	"github.com/gin-gonic/gin"
	"go-admin-full/internal/controllers"
	"go-admin-full/internal/tokenpkg"
	"gorm.io/gorm"
)

func RoleRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenpkg.Manager) {
	rc := controllers.NewRoleController(db)
	rg := r.Group("/api/roles")
	{
		rg.GET("/", rc.List)
		rg.POST("/", rc.Create)
	}
}
