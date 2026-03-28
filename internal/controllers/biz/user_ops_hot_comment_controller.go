package biz

import (
	"go-admin/internal/constants"
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	bizsvc "go-admin/internal/services/biz"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// -------------------- 热点评论 --------------------

// HotCommentController 热点评论控制器。
type HotCommentController struct {
	hotCommentSvc *bizsvc.HotCommentService
}

// NewHotCommentController 创建热点评论控制器。
func NewHotCommentController(db *gorm.DB) *HotCommentController {
	return &HotCommentController{hotCommentSvc: bizsvc.NewHotCommentService(db)}
}

func (uc *HotCommentController) ListHotComments(c *gin.Context) {
	// 定义并初始化当前变量。
	items, err := uc.hotCommentSvc.ListHotComments(c.Request.Context(), 200)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"items": items})
}

// CreateHotComment 创建HotComment。
func (uc *HotCommentController) CreateHotComment(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		PostID uint `json:"post_id"`
		// 处理当前语句逻辑。
		UserID uint `json:"user_id"`
		// 处理当前语句逻辑。
		ParentID *uint `json:"parent_id"`
		// 处理当前语句逻辑。
		Content string `json:"content"`
		// 处理当前语句逻辑。
		Likes *int64 `json:"likes"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
	}
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if req.PostID == 0 || req.UserID == 0 || strings.TrimSpace(req.Content) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "post_id/user_id/content required")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	item, err := uc.hotCommentSvc.CreateHotComment(c.Request.Context(), req.PostID, req.UserID, req.ParentID, strings.TrimSpace(req.Content), req.Likes, req.Status)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateHotComment 更新HotComment。
func (uc *HotCommentController) UpdateHotComment(c *gin.Context) {
	// 定义并初始化当前变量。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid id")
		// 返回当前处理结果。
		return
	}
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Content *string `json:"content"`
		// 处理当前语句逻辑。
		Likes *int64 `json:"likes"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
	}
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	updates := map[string]interface{}{}
	// 判断条件并进入对应分支逻辑。
	if req.Content != nil {
		// 更新当前变量或字段值。
		updates["content"] = strings.TrimSpace(*req.Content)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Likes != nil {
		// 更新当前变量或字段值。
		updates["likes"] = *req.Likes
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		updates["status"] = *req.Status
	}
	// 判断条件并进入对应分支逻辑。
	if len(updates) == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "empty updates")
		// 返回当前处理结果。
		return
	}

	// 判断条件并进入对应分支逻辑。
	if err := uc.hotCommentSvc.UpdateHotComment(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteHotComment 删除HotComment。
func (uc *HotCommentController) DeleteHotComment(c *gin.Context) {
	// 定义并初始化当前变量。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid id")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if err := uc.hotCommentSvc.DeleteHotComment(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}
