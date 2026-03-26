package biz

import (
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/models"

	"github.com/gin-gonic/gin"
)

// -------------------- 发帖管理 --------------------

func (uc *UserOpsController) ListPostArticles(c *gin.Context) {
	// 定义并初始化当前变量。
	items, err := uc.svc.ListPostArticles(c.Request.Context(), false, 200)
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

// CreatePostArticle 创建PostArticle。
func (uc *UserOpsController) CreatePostArticle(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		UserID uint `json:"user_id"`
		// 处理当前语句逻辑。
		Title string `json:"title"`
		// 处理当前语句逻辑。
		CoverImage string `json:"cover_image"`
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
	if strings.TrimSpace(req.Title) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "title required")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if req.UserID == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "user_id required")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if !uc.svc.IsUserTypes(c.Request.Context(), req.UserID, "robot") {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "user_id must be robot account")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	item := models.WPostArticle{
		// 处理当前语句逻辑。
		UserID: req.UserID,
		// 调用strings.TrimSpace完成当前处理。
		Title: strings.TrimSpace(req.Title),
		// 调用strings.TrimSpace完成当前处理。
		CoverImage: strings.TrimSpace(req.CoverImage),
		// 处理当前语句逻辑。
		Content: req.Content,
		// 处理当前语句逻辑。
		IsOfficial: 0,
		// 处理当前语句逻辑。
		Status: 1,
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		item.Status = *req.Status
	}

	// 判断条件并进入对应分支逻辑。
	if err := uc.svc.CreatePostArticle(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdatePostArticle 更新PostArticle。
func (uc *UserOpsController) UpdatePostArticle(c *gin.Context) {
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
		UserID *uint `json:"user_id"`
		// 处理当前语句逻辑。
		Title *string `json:"title"`
		// 处理当前语句逻辑。
		CoverImage *string `json:"cover_image"`
		// 处理当前语句逻辑。
		Content *string `json:"content"`
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
	if req.UserID != nil {
		// 判断条件并进入对应分支逻辑。
		if *req.UserID == 0 || !uc.svc.IsUserTypes(c.Request.Context(), *req.UserID, "robot") {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "user_id must be robot account")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["user_id"] = *req.UserID
	}
	// 判断条件并进入对应分支逻辑。
	if req.Title != nil {
		// 更新当前变量或字段值。
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	// 判断条件并进入对应分支逻辑。
	if req.CoverImage != nil {
		// 更新当前变量或字段值。
		updates["cover_image"] = strings.TrimSpace(*req.CoverImage)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Content != nil {
		// 更新当前变量或字段值。
		updates["content"] = *req.Content
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		updates["status"] = *req.Status
	}
	// 更新当前变量或字段值。
	updates["is_official"] = int8(0)
	// 判断条件并进入对应分支逻辑。
	if len(updates) == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "empty updates")
		// 返回当前处理结果。
		return
	}

	// 判断条件并进入对应分支逻辑。
	if err := uc.svc.UpdatePostArticle(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeletePostArticle 删除PostArticle。
func (uc *UserOpsController) DeletePostArticle(c *gin.Context) {
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
	if err := uc.svc.DeletePostArticle(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}
