package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 外链 --------------------

func (bc *BizConfigController) ListExternalLinks(c *gin.Context) {
	// 定义并初始化当前变量。
	items, err := bc.svc.ListExternalLinks(c.Request.Context(), 200)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{"items": items})
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
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.URL) == "" || strings.TrimSpace(req.Position) == "" {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "name/url/position required")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	item := models.WExternalLink{
		// 调用strings.TrimSpace完成当前处理。
		Name: strings.TrimSpace(req.Name),
		// 调用strings.TrimSpace完成当前处理。
		URL: strings.TrimSpace(req.URL),
		// 调用strings.TrimSpace完成当前处理。
		Position: strings.TrimSpace(req.Position),
		// 调用strings.TrimSpace完成当前处理。
		IconURL: strings.TrimSpace(req.IconURL),
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
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, item)
}

// UpdateExternalLink 更新ExternalLink。
func (bc *BizConfigController) UpdateExternalLink(c *gin.Context) {
	// 定义并初始化当前变量。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		// 返回当前处理结果。
		return
	}
	// 声明当前变量。
	var req map[string]interface{}
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if len(req) == 0 {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.UpdateExternalLink(c.Request.Context(), id, req); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteExternalLink 删除ExternalLink。
func (bc *BizConfigController) DeleteExternalLink(c *gin.Context) {
	// 定义并初始化当前变量。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.DeleteExternalLink(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{"id": id})
}
