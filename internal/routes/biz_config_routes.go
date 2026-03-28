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

// BizConfigRoutes 处理BizConfigRoutes相关逻辑。
func BizConfigRoutes(r *gin.Engine, db *gorm.DB, mgr *tokenjwt.Manager) {
	bannerCtrl := biz.NewBannerController(db)
	broadcastCtrl := biz.NewBroadcastController(db)
	lotteryCtrl := biz.NewLotteryController(db)
	officialPostCtrl := biz.NewOfficialPostController(db)
	externalLinkCtrl := biz.NewExternalLinkController(db)
	homePopupCtrl := biz.NewHomePopupController(db)
	smsChannelCtrl := biz.NewSMSChannelController(db)
	// 定义并初始化当前变量。
	userRoleSvc := rbacsvc.NewUserRoleService(rbacdao.NewUserRoleDao(db))

	// 定义并初始化当前变量。
	group := r.Group("/api/biz-config")
	// 调用group.Use完成当前处理。
	group.Use(middleware.NewJWTMiddleware(mgr))
	{
		// 调用group.GET完成当前处理。
		group.GET("/banners", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), bannerCtrl.ListBanners)
		// 调用group.POST完成当前处理。
		group.POST("/banners", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), bannerCtrl.CreateBanner)
		// 调用group.PUT完成当前处理。
		group.PUT("/banners/:id", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), bannerCtrl.UpdateBanner)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/banners/:id", middleware.PermissionRequired(constants.PermBizBannerList, userRoleSvc, mgr), bannerCtrl.DeleteBanner)

		// 调用group.GET完成当前处理。
		group.GET("/broadcasts", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), broadcastCtrl.ListBroadcasts)
		// 调用group.POST完成当前处理。
		group.POST("/broadcasts", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), broadcastCtrl.CreateBroadcast)
		// 调用group.PUT完成当前处理。
		group.PUT("/broadcasts/:id", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), broadcastCtrl.UpdateBroadcast)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/broadcasts/:id", middleware.PermissionRequired(constants.PermBizBroadcastList, userRoleSvc, mgr), broadcastCtrl.DeleteBroadcast)

		// 调用group.GET完成当前处理。
		group.GET("/special-lotteries", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), lotteryCtrl.ListSpecialLotteries)
		// 调用group.POST完成当前处理。
		group.POST("/special-lotteries", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), lotteryCtrl.CreateSpecialLottery)
		// 调用group.PUT完成当前处理。
		group.PUT("/special-lotteries/:id", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), lotteryCtrl.UpdateSpecialLottery)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/special-lotteries/:id", middleware.PermissionRequired(constants.PermBizSpecialLotteryList, userRoleSvc, mgr), lotteryCtrl.DeleteSpecialLottery)

		// 调用group.GET完成当前处理。
		group.GET("/lottery-infos", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.ListLotteryInfos)
		// 调用group.POST完成当前处理。
		group.POST("/lottery-infos", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.CreateLotteryInfo)
		// 调用group.PUT完成当前处理。
		group.PUT("/lottery-infos/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.UpdateLotteryInfo)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/lottery-infos/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.DeleteLotteryInfo)

		// 调用group.GET完成当前处理。
		group.GET("/draw-records", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.ListDrawRecords)
		// 调用group.POST完成当前处理。
		group.POST("/draw-records", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.CreateDrawRecord)
		// 调用group.PUT完成当前处理。
		group.PUT("/draw-records/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.UpdateDrawRecord)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/draw-records/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.DeleteDrawRecord)

		// 调用group.GET完成当前处理。
		group.GET("/lottery-categories", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.ListLotteryCategories)
		// 调用group.POST完成当前处理。
		group.POST("/lottery-categories", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.CreateLotteryCategory)
		// 调用group.PUT完成当前处理。
		group.PUT("/lottery-categories/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.UpdateLotteryCategory)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/lottery-categories/:id", middleware.PermissionRequired(constants.PermBizLotteryInfoList, userRoleSvc, mgr), lotteryCtrl.DeleteLotteryCategory)

		// 调用group.GET完成当前处理。
		group.GET("/official-posts", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), officialPostCtrl.ListOfficialPosts)
		// 调用group.POST完成当前处理。
		group.POST("/official-posts", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), officialPostCtrl.CreateOfficialPost)
		// 调用group.PUT完成当前处理。
		group.PUT("/official-posts/:id", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), officialPostCtrl.UpdateOfficialPost)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/official-posts/:id", middleware.PermissionRequired(constants.PermBizOfficialPostList, userRoleSvc, mgr), officialPostCtrl.DeleteOfficialPost)

		// 调用group.GET完成当前处理。
		group.GET("/external-links", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), externalLinkCtrl.ListExternalLinks)
		// 调用group.POST完成当前处理。
		group.POST("/external-links", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), externalLinkCtrl.CreateExternalLink)
		// 调用group.PUT完成当前处理。
		group.PUT("/external-links/:id", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), externalLinkCtrl.UpdateExternalLink)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/external-links/:id", middleware.PermissionRequired(constants.PermBizExternalLinkList, userRoleSvc, mgr), externalLinkCtrl.DeleteExternalLink)

		// 调用group.GET完成当前处理。
		group.GET("/home-popups", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), homePopupCtrl.ListHomePopups)
		// 调用group.POST完成当前处理。
		group.POST("/home-popups", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), homePopupCtrl.CreateHomePopup)
		// 调用group.PUT完成当前处理。
		group.PUT("/home-popups/:id", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), homePopupCtrl.UpdateHomePopup)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/home-popups/:id", middleware.PermissionRequired(constants.PermBizHomePopupList, userRoleSvc, mgr), homePopupCtrl.DeleteHomePopup)

		// 调用group.GET完成当前处理。
		group.GET("/sms-channels", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), smsChannelCtrl.ListSMSChannels)
		// 调用group.POST完成当前处理。
		group.POST("/sms-channels", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), smsChannelCtrl.CreateSMSChannel)
		// 调用group.PUT完成当前处理。
		group.PUT("/sms-channels/:id", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), smsChannelCtrl.UpdateSMSChannel)
		// 调用group.DELETE完成当前处理。
		group.DELETE("/sms-channels/:id", middleware.PermissionRequired(constants.PermBizSMSChannelList, userRoleSvc, mgr), smsChannelCtrl.DeleteSMSChannel)
	}
}
