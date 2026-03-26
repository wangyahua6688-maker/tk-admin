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

// BizConfigRoutes 处理BizConfigRoutes相关逻辑。
func BizConfigRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	// 定义并初始化当前变量。
	ctrl := biz.NewBizConfigController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	group := r.Group("/api/biz-config")
	// 调用group.Use完成当前处理。
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用group.GET完成当前处理。
		group.GET("/banners", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), ctrl.ListBanners)
		// 调用group.POST完成当前处理。
		group.POST("/banners", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), ctrl.CreateBanner)
		// 调用group.PUT完成当前处理。
		group.PUT("/banners/:id", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), ctrl.UpdateBanner)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/banners/:id", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), ctrl.DeleteBanner)

		// 调用group.GET完成当前处理。
		group.GET("/broadcasts", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), ctrl.ListBroadcasts)
		// 调用group.POST完成当前处理。
		group.POST("/broadcasts", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), ctrl.CreateBroadcast)
		// 调用group.PUT完成当前处理。
		group.PUT("/broadcasts/:id", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), ctrl.UpdateBroadcast)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/broadcasts/:id", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), ctrl.DeleteBroadcast)

		// 调用group.GET完成当前处理。
		group.GET("/special-lotteries", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), ctrl.ListSpecialLotteries)
		// 调用group.POST完成当前处理。
		group.POST("/special-lotteries", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), ctrl.CreateSpecialLottery)
		// 调用group.PUT完成当前处理。
		group.PUT("/special-lotteries/:id", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), ctrl.UpdateSpecialLottery)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/special-lotteries/:id", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), ctrl.DeleteSpecialLottery)

		// 调用group.GET完成当前处理。
		group.GET("/lottery-infos", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.ListLotteryInfos)
		// 调用group.POST完成当前处理。
		group.POST("/lottery-infos", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.CreateLotteryInfo)
		// 调用group.PUT完成当前处理。
		group.PUT("/lottery-infos/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.UpdateLotteryInfo)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/lottery-infos/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.DeleteLotteryInfo)

		// 调用group.GET完成当前处理。
		group.GET("/draw-records", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.ListDrawRecords)
		// 调用group.POST完成当前处理。
		group.POST("/draw-records", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.CreateDrawRecord)
		// 调用group.PUT完成当前处理。
		group.PUT("/draw-records/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.UpdateDrawRecord)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/draw-records/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.DeleteDrawRecord)

		// 调用group.GET完成当前处理。
		group.GET("/lottery-categories", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.ListLotteryCategories)
		// 调用group.POST完成当前处理。
		group.POST("/lottery-categories", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.CreateLotteryCategory)
		// 调用group.PUT完成当前处理。
		group.PUT("/lottery-categories/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.UpdateLotteryCategory)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/lottery-categories/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), ctrl.DeleteLotteryCategory)

		// 调用group.GET完成当前处理。
		group.GET("/official-posts", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), ctrl.ListOfficialPosts)
		// 调用group.POST完成当前处理。
		group.POST("/official-posts", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), ctrl.CreateOfficialPost)
		// 调用group.PUT完成当前处理。
		group.PUT("/official-posts/:id", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), ctrl.UpdateOfficialPost)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/official-posts/:id", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), ctrl.DeleteOfficialPost)

		// 调用group.GET完成当前处理。
		group.GET("/external-links", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), ctrl.ListExternalLinks)
		// 调用group.POST完成当前处理。
		group.POST("/external-links", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), ctrl.CreateExternalLink)
		// 调用group.PUT完成当前处理。
		group.PUT("/external-links/:id", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), ctrl.UpdateExternalLink)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/external-links/:id", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), ctrl.DeleteExternalLink)

		// 调用group.GET完成当前处理。
		group.GET("/home-popups", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), ctrl.ListHomePopups)
		// 调用group.POST完成当前处理。
		group.POST("/home-popups", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), ctrl.CreateHomePopup)
		// 调用group.PUT完成当前处理。
		group.PUT("/home-popups/:id", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), ctrl.UpdateHomePopup)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/home-popups/:id", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), ctrl.DeleteHomePopup)

		// 调用group.GET完成当前处理。
		group.GET("/sms-channels", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), ctrl.ListSMSChannels)
		// 调用group.POST完成当前处理。
		group.POST("/sms-channels", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), ctrl.CreateSMSChannel)
		// 调用group.PUT完成当前处理。
		group.PUT("/sms-channels/:id", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), ctrl.UpdateSMSChannel)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/sms-channels/:id", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), ctrl.DeleteSMSChannel)
	}
}
