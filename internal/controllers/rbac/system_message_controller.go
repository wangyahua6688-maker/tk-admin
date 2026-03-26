package rbac

import (
	"errors"
	"strconv"
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin-full/internal/constants"
	rbacdao "go-admin-full/internal/dao/rbac"
	rbacsvc "go-admin-full/internal/services/rbac"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SystemMessageController 系统消息控制器。
type SystemMessageController struct {
	svc *rbacsvc.SystemMessageService // 系统消息服务
}

// NewSystemMessageController 创建SystemMessageController实例。
func NewSystemMessageController(db *gorm.DB) *SystemMessageController {
	// 返回当前处理结果。
	return &SystemMessageController{
		svc: rbacsvc.NewSystemMessageService(rbacdao.NewSystemMessageDao(db)), // 注入系统消息服务
	}
}

// ListMyMessages 查询当前登录用户的系统消息列表。
func (mc *SystemMessageController) ListMyMessages(c *gin.Context) {
	uid := c.GetUint("uid") // 从 JWT 中读取用户 ID
	// 判断条件并进入对应分支逻辑。
	if uid == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminAuthUnauthorized, "用户未认证")
		// 返回当前处理结果。
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))           // 页码
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20")) // 每页数量
	// 判断条件并进入对应分支逻辑。
	if page <= 0 {
		// 更新当前变量或字段值。
		page = 1
	}
	// 判断条件并进入对应分支逻辑。
	if pageSize <= 0 {
		// 更新当前变量或字段值。
		pageSize = 20
	}
	// 判断条件并进入对应分支逻辑。
	if pageSize > 100 {
		// 更新当前变量或字段值。
		pageSize = 100
	}

	onlyUnread := false                                                    // 是否仅未读
	onlyUnreadRaw := strings.TrimSpace(c.DefaultQuery("only_unread", "0")) // 读取查询参数
	// 判断条件并进入对应分支逻辑。
	if onlyUnreadRaw == "1" || strings.EqualFold(onlyUnreadRaw, "true") {
		// 更新当前变量或字段值。
		onlyUnread = true
	}

	// 查询消息列表与未读数
	items, total, unread, err := mc.svc.ListUserMessages(c.Request.Context(), uid, page, pageSize, onlyUnread)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}

	// 返回分页数据
	commonresp.GinOK(c, gin.H{
		// 处理当前语句逻辑。
		"items": items,
		// 处理当前语句逻辑。
		"total": total,
		// 处理当前语句逻辑。
		"unread_count": unread,
		// 处理当前语句逻辑。
		"page": page,
		// 处理当前语句逻辑。
		"page_size": pageSize,
	})
}

// MarkRead 将指定系统消息标记为已读。
func (mc *SystemMessageController) MarkRead(c *gin.Context) {
	uid := c.GetUint("uid") // 从 JWT 中读取用户 ID
	// 判断条件并进入对应分支逻辑。
	if uid == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminAuthUnauthorized, "用户未认证")
		// 返回当前处理结果。
		return
	}

	// 解析并校验消息 ID
	id, err := strconv.Atoi(c.Param("id"))
	// 判断条件并进入对应分支逻辑。
	if err != nil || id <= 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid message id")
		// 返回当前处理结果。
		return
	}

	// 标记已读
	if err := mc.svc.MarkRead(c.Request.Context(), uid, uint(id)); err != nil {
		// 判断条件并进入对应分支逻辑。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizResourceNotFound, "消息不存在")
			// 返回当前处理结果。
			return
		}
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"msg": "ok"})
}

// MarkAllRead 将当前用户全部消息标记为已读。
func (mc *SystemMessageController) MarkAllRead(c *gin.Context) {
	uid := c.GetUint("uid") // 从 JWT 中读取用户 ID
	// 判断条件并进入对应分支逻辑。
	if uid == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminAuthUnauthorized, "用户未认证")
		// 返回当前处理结果。
		return
	}

	// 标记全部已读
	if err := mc.svc.MarkAllRead(c.Request.Context(), uid); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"msg": "ok"})
}
