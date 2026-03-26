package biz

import (
	"go-admin/internal/constants"
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"

	"github.com/gin-gonic/gin"
)

// -------------------- 帖子评论管理（按帖子维度） --------------------

func (uc *UserOpsController) ListPostComments(c *gin.Context) {
	// 定义并初始化当前变量。
	postID, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid post id")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	items, err := uc.svc.ListPostComments(c.Request.Context(), postID)
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

// CreatePostComment 创建PostComment。
func (uc *UserOpsController) CreatePostComment(c *gin.Context) {
	// 定义并初始化当前变量。
	postID, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid post id")
		// 返回当前处理结果。
		return
	}

	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		UserID uint `json:"user_id"`
		// 处理当前语句逻辑。
		ParentID uint `json:"parent_id"`
		// 处理当前语句逻辑。
		Content string `json:"content"`
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
	if req.UserID == 0 || strings.TrimSpace(req.Content) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "user_id/content required")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	item, err := uc.svc.CreatePostComment(c.Request.Context(), postID, req.UserID, req.ParentID, strings.TrimSpace(req.Content), req.Status)
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
