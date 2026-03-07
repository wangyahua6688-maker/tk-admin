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

func BizConfigRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	ctrl := controllers.NewBizConfigController(db)
	userRoleSvc := services.NewUserRoleService(dao.NewUserRoleDao(db))

	group := r.Group("/api/biz-config")
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		group.GET("/banners", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc), ctrl.ListBanners)
		group.POST("/banners", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc), ctrl.CreateBanner)
		group.PUT("/banners/:id", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc), ctrl.UpdateBanner)
		group.DELETE("/banners/:id", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc), ctrl.DeleteBanner)

		group.GET("/broadcasts", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc), ctrl.ListBroadcasts)
		group.POST("/broadcasts", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc), ctrl.CreateBroadcast)
		group.PUT("/broadcasts/:id", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc), ctrl.UpdateBroadcast)
		group.DELETE("/broadcasts/:id", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc), ctrl.DeleteBroadcast)

		group.GET("/special-lotteries", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc), ctrl.ListSpecialLotteries)
		group.POST("/special-lotteries", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc), ctrl.CreateSpecialLottery)
		group.PUT("/special-lotteries/:id", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc), ctrl.UpdateSpecialLottery)
		group.DELETE("/special-lotteries/:id", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc), ctrl.DeleteSpecialLottery)

		group.GET("/lottery-infos", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.ListLotteryInfos)
		group.POST("/lottery-infos", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.CreateLotteryInfo)
		group.PUT("/lottery-infos/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.UpdateLotteryInfo)
		group.DELETE("/lottery-infos/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.DeleteLotteryInfo)

		group.GET("/draw-records", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.ListDrawRecords)
		group.POST("/draw-records", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.CreateDrawRecord)
		group.PUT("/draw-records/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.UpdateDrawRecord)
		group.DELETE("/draw-records/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.DeleteDrawRecord)

		group.GET("/lottery-categories", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.ListLotteryCategories)
		group.POST("/lottery-categories", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.CreateLotteryCategory)
		group.PUT("/lottery-categories/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.UpdateLotteryCategory)
		group.DELETE("/lottery-categories/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc), ctrl.DeleteLotteryCategory)

		group.GET("/official-posts", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc), ctrl.ListOfficialPosts)
		group.POST("/official-posts", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc), ctrl.CreateOfficialPost)
		group.PUT("/official-posts/:id", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc), ctrl.UpdateOfficialPost)
		group.DELETE("/official-posts/:id", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc), ctrl.DeleteOfficialPost)

		group.GET("/external-links", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc), ctrl.ListExternalLinks)
		group.POST("/external-links", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc), ctrl.CreateExternalLink)
		group.PUT("/external-links/:id", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc), ctrl.UpdateExternalLink)
		group.DELETE("/external-links/:id", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc), ctrl.DeleteExternalLink)

		group.GET("/home-popups", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc), ctrl.ListHomePopups)
		group.POST("/home-popups", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc), ctrl.CreateHomePopup)
		group.PUT("/home-popups/:id", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc), ctrl.UpdateHomePopup)
		group.DELETE("/home-popups/:id", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc), ctrl.DeleteHomePopup)

		group.GET("/sms-channels", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc), ctrl.ListSMSChannels)
		group.POST("/sms-channels", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc), ctrl.CreateSMSChannel)
		group.PUT("/sms-channels/:id", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc), ctrl.UpdateSMSChannel)
		group.DELETE("/sms-channels/:id", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc), ctrl.DeleteSMSChannel)
	}
}
