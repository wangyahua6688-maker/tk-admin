package biz

import (
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin/internal/constants"
	"go-admin/internal/models"

	"github.com/gin-gonic/gin"
)

// -------------------- 外链 --------------------

func (bc *BizConfigController) ListExternalLinks(c *gin.Context) {
	// 定义并初始化当前变量。
	items, err := bc.svc.ListExternalLinks(c.Request.Context(), 200)
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

// CreateExternalLink 创建ExternalLink。
func (bc *BizConfigController) CreateExternalLink(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Name string `json:"name"`
		// 处理当前语句逻辑。
		URL string `json:"url"`
		// 处理当前语句逻辑。
		Position string `json:"position"`
		// 处理当前语句逻辑。
		IconURL string `json:"icon_url"`
		// 处理当前语句逻辑。
		GroupKey string `json:"group_key"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
		// 处理当前语句逻辑。
		Sort *int `json:"sort"`
	}
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.URL) == "" || strings.TrimSpace(req.Position) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "name/url/position required")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	linkURL, err := normalizeSafeURL(req.URL, true)
	// 判断条件并进入对应分支逻辑。
	if err != nil || strings.TrimSpace(linkURL) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid url")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	iconURL, err := normalizeSafeURL(req.IconURL, true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid icon_url")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	item := models.WExternalLink{
		// 调用strings.TrimSpace完成当前处理。
		Name: strings.TrimSpace(req.Name),
		// 处理当前语句逻辑。
		URL: linkURL,
		// 调用strings.TrimSpace完成当前处理。
		Position: strings.TrimSpace(req.Position),
		// 处理当前语句逻辑。
		IconURL: iconURL,
		// 调用strings.TrimSpace完成当前处理。
		GroupKey: strings.TrimSpace(req.GroupKey),
		// 处理当前语句逻辑。
		Status: 1,
		// 处理当前语句逻辑。
		Sort: 0,
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
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.CreateExternalLink(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateExternalLink 更新ExternalLink。
func (bc *BizConfigController) UpdateExternalLink(c *gin.Context) {
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
	var req map[string]interface{}
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if len(req) == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "empty updates")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if raw, ok := req["url"]; ok {
		// 定义并初始化当前变量。
		urlRaw, ok := raw.(string)
		// 判断条件并进入对应分支逻辑。
		if !ok {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid url")
			// 返回当前处理结果。
			return
		}
		// 定义并初始化当前变量。
		linkURL, urlErr := normalizeSafeURL(urlRaw, true)
		// 判断条件并进入对应分支逻辑。
		if urlErr != nil || strings.TrimSpace(linkURL) == "" {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		req["url"] = linkURL
	}
	// 判断条件并进入对应分支逻辑。
	if raw, ok := req["icon_url"]; ok {
		// 定义并初始化当前变量。
		iconRaw, ok := raw.(string)
		// 判断条件并进入对应分支逻辑。
		if !ok {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid icon_url")
			// 返回当前处理结果。
			return
		}
		// 定义并初始化当前变量。
		iconURL, iconErr := normalizeSafeURL(iconRaw, true)
		// 判断条件并进入对应分支逻辑。
		if iconErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid icon_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		req["icon_url"] = iconURL
	}
	// 判断条件并进入对应分支逻辑。
	if raw, ok := req["name"]; ok {
		// 定义并初始化当前变量。
		nameRaw, ok := raw.(string)
		// 判断条件并进入对应分支逻辑。
		if !ok {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid name")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		req["name"] = strings.TrimSpace(nameRaw)
	}
	// 判断条件并进入对应分支逻辑。
	if raw, ok := req["position"]; ok {
		// 定义并初始化当前变量。
		positionRaw, ok := raw.(string)
		// 判断条件并进入对应分支逻辑。
		if !ok {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid position")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		req["position"] = strings.TrimSpace(positionRaw)
	}
	// 判断条件并进入对应分支逻辑。
	if raw, ok := req["group_key"]; ok {
		// 定义并初始化当前变量。
		groupRaw, ok := raw.(string)
		// 判断条件并进入对应分支逻辑。
		if !ok {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid group_key")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		req["group_key"] = strings.TrimSpace(groupRaw)
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.UpdateExternalLink(c.Request.Context(), id, req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteExternalLink 删除ExternalLink。
func (bc *BizConfigController) DeleteExternalLink(c *gin.Context) {
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
	if err := bc.svc.DeleteExternalLink(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}
