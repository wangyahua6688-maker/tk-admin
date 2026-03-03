package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/dao"
	"go-admin-full/internal/services"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
)

// SystemMessageController 系统消息控制器。
type SystemMessageController struct {
	svc *services.SystemMessageService
}

func NewSystemMessageController(db *gorm.DB) *SystemMessageController {
	return &SystemMessageController{
		svc: services.NewSystemMessageService(dao.NewSystemMessageDao(db)),
	}
}

// ListMyMessages 查询当前登录用户的系统消息列表。
func (mc *SystemMessageController) ListMyMessages(c *gin.Context) {
	uid := c.GetUint("uid")
	if uid == 0 {
		utils.JSONError(c, http.StatusUnauthorized, "用户未认证")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	onlyUnread := false
	onlyUnreadRaw := strings.TrimSpace(c.DefaultQuery("only_unread", "0"))
	if onlyUnreadRaw == "1" || strings.EqualFold(onlyUnreadRaw, "true") {
		onlyUnread = true
	}

	items, total, unread, err := mc.svc.ListUserMessages(c.Request.Context(), uid, page, pageSize, onlyUnread)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONOK(c, gin.H{
		"items":        items,
		"total":        total,
		"unread_count": unread,
		"page":         page,
		"page_size":    pageSize,
	})
}

// MarkRead 将指定系统消息标记为已读。
func (mc *SystemMessageController) MarkRead(c *gin.Context) {
	uid := c.GetUint("uid")
	if uid == 0 {
		utils.JSONError(c, http.StatusUnauthorized, "用户未认证")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.JSONError(c, http.StatusBadRequest, "invalid message id")
		return
	}

	if err := mc.svc.MarkRead(c.Request.Context(), uid, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONError(c, http.StatusNotFound, "消息不存在")
			return
		}
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONOK(c, gin.H{"msg": "ok"})
}

// MarkAllRead 将当前用户全部消息标记为已读。
func (mc *SystemMessageController) MarkAllRead(c *gin.Context) {
	uid := c.GetUint("uid")
	if uid == 0 {
		utils.JSONError(c, http.StatusUnauthorized, "用户未认证")
		return
	}

	if err := mc.svc.MarkAllRead(c.Request.Context(), uid); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONOK(c, gin.H{"msg": "ok"})
}
