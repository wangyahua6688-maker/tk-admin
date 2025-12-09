package routes

import (
    "github.com/gin-gonic/gin"
    "go-admin-full/internal/controllers"
    "go-admin-full/internal/middleware"
    "go-admin-full/internal/tokenpkg"
    "gorm.io/gorm"
)

func MenuRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenpkg.Manager) {
    mc := controllers.NewMenuController(db)
    mr := r.Group("/api/menus")
    mr.Use(middleware.NewJWTMiddleware(mgr))
    {
        mr.GET("/", mc.List)
        mr.POST("/", mc.Create)
        mr.GET("/:id", mc.Get)
        mr.DELETE("/:id", mc.Delete)
    }
}
