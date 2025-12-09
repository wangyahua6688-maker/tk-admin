package routes

import (
    "github.com/gin-gonic/gin"
    "go-admin-full/internal/controllers"
    "go-admin-full/internal/middleware"
    "go-admin-full/internal/tokenpkg"
    "gorm.io/gorm"
)

func PermissionRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenpkg.Manager) {
    pc := controllers.NewPermissionController(db)
    pr := r.Group("/api/permissions")
    pr.Use(middleware.NewJWTMiddleware(mgr))
    {
        pr.GET("/", pc.List)
        pr.POST("/", pc.Create)
        pr.GET("/:id", pc.Get)
        pr.DELETE("/:id", pc.Delete)
    }
}
