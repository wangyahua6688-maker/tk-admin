package biz

import (
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/models"

	"github.com/gin-gonic/gin"
)

// -------------------- 首页首屏弹窗 --------------------

// ListHomePopups 查询首页首屏弹窗列表。
func (bc *BizConfigController) ListHomePopups(c *gin.Context) {
	// 读取位置筛选参数，默认只看首页弹窗。
	position := strings.TrimSpace(c.Query("position"))
	// 判断条件并进入对应分支逻辑。
	if position == "" {
		// 更新当前变量或字段值。
		position = "home"
	}

	// 执行查询并输出结果。
	items, err := bc.svc.ListHomePopups(c.Request.Context(), position, 200)
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

// CreateHomePopup 新增首页首屏弹窗。
func (bc *BizConfigController) CreateHomePopup(c *gin.Context) {
	// 定义请求结构，便于字段级校验。
	var req struct {
		// 处理当前语句逻辑。
		Title string `json:"title"`
		// 处理当前语句逻辑。
		Content string `json:"content"`
		// 处理当前语句逻辑。
		ImageURL string `json:"image_url"`
		// 处理当前语句逻辑。
		ButtonText string `json:"button_text"`
		// 处理当前语句逻辑。
		ButtonLink string `json:"button_link"`
		// 处理当前语句逻辑。
		Position string `json:"position"`
		// 处理当前语句逻辑。
		ShowOnce *int8 `json:"show_once"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
		// 处理当前语句逻辑。
		Sort *int `json:"sort"`
		// 处理当前语句逻辑。
		StartAt *string `json:"start_at"`
		// 处理当前语句逻辑。
		EndAt *string `json:"end_at"`
	}

	// 绑定并校验请求体。
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
	// 定义并初始化当前变量。
	imageURL, err := normalizeSafeURL(req.ImageURL, true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid image_url")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	buttonLink, err := normalizeSafeURL(req.ButtonLink, true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid button_link")
		// 返回当前处理结果。
		return
	}

	// 组装模型并写入默认值。
	item := models.WHomePopup{
		// 调用strings.TrimSpace完成当前处理。
		Title: strings.TrimSpace(req.Title),
		// 调用strings.TrimSpace完成当前处理。
		Content: strings.TrimSpace(req.Content),
		// 处理当前语句逻辑。
		ImageURL: imageURL,
		// 调用strings.TrimSpace完成当前处理。
		ButtonText: strings.TrimSpace(req.ButtonText),
		// 处理当前语句逻辑。
		ButtonLink: buttonLink,
		// 调用strings.TrimSpace完成当前处理。
		Position: strings.TrimSpace(req.Position),
		// 处理当前语句逻辑。
		ShowOnce: 1,
		// 处理当前语句逻辑。
		Status: 1,
		// 处理当前语句逻辑。
		Sort: 0,
		// 调用parseRFC3339Ptr完成当前处理。
		StartAt: parseRFC3339Ptr(req.StartAt),
		// 调用parseRFC3339Ptr完成当前处理。
		EndAt: parseRFC3339Ptr(req.EndAt),
	}
	// 判断条件并进入对应分支逻辑。
	if item.Position == "" {
		// 更新当前变量或字段值。
		item.Position = "home"
	}
	// 判断条件并进入对应分支逻辑。
	if req.ShowOnce != nil {
		// 更新当前变量或字段值。
		item.ShowOnce = *req.ShowOnce
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		item.Status = *req.Status
	}
	// 判断条件并进入对应分支逻辑。
	if req.Sort != nil {
		// 更新当前变量或字段值。
		item.Sort = *req.Sort
	}

	// 落库并返回新建记录。
	if err := bc.svc.CreateHomePopup(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateHomePopup 更新首页首屏弹窗。
func (bc *BizConfigController) UpdateHomePopup(c *gin.Context) {
	// 解析路由主键ID。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid id")
		// 返回当前处理结果。
		return
	}

	// 定义更新请求结构，支持字段级部分更新。
	var req struct {
		// 处理当前语句逻辑。
		Title *string `json:"title"`
		// 处理当前语句逻辑。
		Content *string `json:"content"`
		// 处理当前语句逻辑。
		ImageURL *string `json:"image_url"`
		// 处理当前语句逻辑。
		ButtonText *string `json:"button_text"`
		// 处理当前语句逻辑。
		ButtonLink *string `json:"button_link"`
		// 处理当前语句逻辑。
		Position *string `json:"position"`
		// 处理当前语句逻辑。
		ShowOnce *int8 `json:"show_once"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
		// 处理当前语句逻辑。
		Sort *int `json:"sort"`
		// 处理当前语句逻辑。
		StartAt *string `json:"start_at"`
		// 处理当前语句逻辑。
		EndAt *string `json:"end_at"`
	}

	// 绑定请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}

	// 构建更新字段映射，避免覆盖未传字段。
	updates := make(map[string]interface{})
	// 判断条件并进入对应分支逻辑。
	if req.Title != nil {
		// 定义并初始化当前变量。
		title := strings.TrimSpace(*req.Title)
		// 判断条件并进入对应分支逻辑。
		if title == "" {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "title required")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["title"] = title
	}
	// 判断条件并进入对应分支逻辑。
	if req.Content != nil {
		// 更新当前变量或字段值。
		updates["content"] = strings.TrimSpace(*req.Content)
	}
	// 判断条件并进入对应分支逻辑。
	if req.ImageURL != nil {
		// 定义并初始化当前变量。
		imageURL, imageErr := normalizeSafeURL(*req.ImageURL, true)
		// 判断条件并进入对应分支逻辑。
		if imageErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid image_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["image_url"] = imageURL
	}
	// 判断条件并进入对应分支逻辑。
	if req.ButtonText != nil {
		// 更新当前变量或字段值。
		updates["button_text"] = strings.TrimSpace(*req.ButtonText)
	}
	// 判断条件并进入对应分支逻辑。
	if req.ButtonLink != nil {
		// 定义并初始化当前变量。
		buttonLink, linkErr := normalizeSafeURL(*req.ButtonLink, true)
		// 判断条件并进入对应分支逻辑。
		if linkErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid button_link")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["button_link"] = buttonLink
	}
	// 判断条件并进入对应分支逻辑。
	if req.Position != nil {
		// 定义并初始化当前变量。
		position := strings.TrimSpace(*req.Position)
		// 判断条件并进入对应分支逻辑。
		if position == "" {
			// 更新当前变量或字段值。
			position = "home"
		}
		// 更新当前变量或字段值。
		updates["position"] = position
	}
	// 判断条件并进入对应分支逻辑。
	if req.ShowOnce != nil {
		// 更新当前变量或字段值。
		updates["show_once"] = *req.ShowOnce
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		updates["status"] = *req.Status
	}
	// 判断条件并进入对应分支逻辑。
	if req.Sort != nil {
		// 更新当前变量或字段值。
		updates["sort"] = *req.Sort
	}
	// 判断条件并进入对应分支逻辑。
	if req.StartAt != nil {
		// 更新当前变量或字段值。
		updates["start_at"] = parseRFC3339Ptr(req.StartAt)
	}
	// 判断条件并进入对应分支逻辑。
	if req.EndAt != nil {
		// 更新当前变量或字段值。
		updates["end_at"] = parseRFC3339Ptr(req.EndAt)
	}
	// 判断条件并进入对应分支逻辑。
	if len(updates) == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "empty updates")
		// 返回当前处理结果。
		return
	}

	// 执行更新并返回ID。
	if err := bc.svc.UpdateHomePopup(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteHomePopup 删除首页首屏弹窗。
func (bc *BizConfigController) DeleteHomePopup(c *gin.Context) {
	// 解析路由主键ID。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid id")
		// 返回当前处理结果。
		return
	}

	// 执行删除。
	if err := bc.svc.DeleteHomePopup(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}
