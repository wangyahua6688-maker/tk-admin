package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 首页首屏弹窗 --------------------

// ListHomePopups 查询首页首屏弹窗列表。
func (bc *BizConfigController) ListHomePopups(c *gin.Context) {
	// 读取位置筛选参数，默认只看首页弹窗。
	position := strings.TrimSpace(c.Query("position"))
	if position == "" {
		position = "home"
	}

	// 执行查询并输出结果。
	items, err := bc.svc.ListHomePopups(c.Request.Context(), position, 200)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

// CreateHomePopup 新增首页首屏弹窗。
func (bc *BizConfigController) CreateHomePopup(c *gin.Context) {
	// 定义请求结构，便于字段级校验。
	var req struct {
		Title      string  `json:"title"`
		Content    string  `json:"content"`
		ImageURL   string  `json:"image_url"`
		ButtonText string  `json:"button_text"`
		ButtonLink string  `json:"button_link"`
		Position   string  `json:"position"`
		ShowOnce   *int8   `json:"show_once"`
		Status     *int8   `json:"status"`
		Sort       *int    `json:"sort"`
		StartAt    *string `json:"start_at"`
		EndAt      *string `json:"end_at"`
	}

	// 绑定并校验请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		utils.JSONError(c, http.StatusBadRequest, "title required")
		return
	}

	// 组装模型并写入默认值。
	item := models.WHomePopup{
		Title:      strings.TrimSpace(req.Title),
		Content:    strings.TrimSpace(req.Content),
		ImageURL:   strings.TrimSpace(req.ImageURL),
		ButtonText: strings.TrimSpace(req.ButtonText),
		ButtonLink: strings.TrimSpace(req.ButtonLink),
		Position:   strings.TrimSpace(req.Position),
		ShowOnce:   1,
		Status:     1,
		Sort:       0,
		StartAt:    parseRFC3339Ptr(req.StartAt),
		EndAt:      parseRFC3339Ptr(req.EndAt),
	}
	if item.Position == "" {
		item.Position = "home"
	}
	if req.ShowOnce != nil {
		item.ShowOnce = *req.ShowOnce
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Sort != nil {
		item.Sort = *req.Sort
	}

	// 落库并返回新建记录。
	if err := bc.svc.CreateHomePopup(c.Request.Context(), &item); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

// UpdateHomePopup 更新首页首屏弹窗。
func (bc *BizConfigController) UpdateHomePopup(c *gin.Context) {
	// 解析路由主键ID。
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	// 定义更新请求结构，支持字段级部分更新。
	var req struct {
		Title      *string `json:"title"`
		Content    *string `json:"content"`
		ImageURL   *string `json:"image_url"`
		ButtonText *string `json:"button_text"`
		ButtonLink *string `json:"button_link"`
		Position   *string `json:"position"`
		ShowOnce   *int8   `json:"show_once"`
		Status     *int8   `json:"status"`
		Sort       *int    `json:"sort"`
		StartAt    *string `json:"start_at"`
		EndAt      *string `json:"end_at"`
	}

	// 绑定请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	// 构建更新字段映射，避免覆盖未传字段。
	updates := make(map[string]interface{})
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			utils.JSONError(c, http.StatusBadRequest, "title required")
			return
		}
		updates["title"] = title
	}
	if req.Content != nil {
		updates["content"] = strings.TrimSpace(*req.Content)
	}
	if req.ImageURL != nil {
		updates["image_url"] = strings.TrimSpace(*req.ImageURL)
	}
	if req.ButtonText != nil {
		updates["button_text"] = strings.TrimSpace(*req.ButtonText)
	}
	if req.ButtonLink != nil {
		updates["button_link"] = strings.TrimSpace(*req.ButtonLink)
	}
	if req.Position != nil {
		position := strings.TrimSpace(*req.Position)
		if position == "" {
			position = "home"
		}
		updates["position"] = position
	}
	if req.ShowOnce != nil {
		updates["show_once"] = *req.ShowOnce
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.StartAt != nil {
		updates["start_at"] = parseRFC3339Ptr(req.StartAt)
	}
	if req.EndAt != nil {
		updates["end_at"] = parseRFC3339Ptr(req.EndAt)
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}

	// 执行更新并返回ID。
	if err := bc.svc.UpdateHomePopup(c.Request.Context(), id, updates); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteHomePopup 删除首页首屏弹窗。
func (bc *BizConfigController) DeleteHomePopup(c *gin.Context) {
	// 解析路由主键ID。
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	// 执行删除。
	if err := bc.svc.DeleteHomePopup(c.Request.Context(), id); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}
