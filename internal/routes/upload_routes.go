package routes

import (
	"github.com/gin-gonic/gin"
	biz "go-admin-full/internal/controllers/biz"
	"go-admin-full/internal/middleware"
	tokenjwt "go-admin-full/internal/token/jwt"
)

func UploadRoutes(r *gin.Engine, mgr *tokenjwt.Manager) {
	ctrl := biz.NewUploadController()

	group := r.Group("/api")
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		group.POST("/upload/image", ctrl.UploadImage)
	}
}
