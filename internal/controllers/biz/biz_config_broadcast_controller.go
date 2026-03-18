package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- Broadcast --------------------

func (bc *BizConfigController) ListBroadcasts(c *gin.Context) {
	// 定义并初始化当前变量。
	items, err := bc.svc.ListBroadcasts(c.Request.Context(), 200)
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

// CreateBroadcast 创建Broadcast。
func (bc *BizConfigController) CreateBroadcast(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Title string `json:"title"`
		// 处理当前语句逻辑。
		Content string `json:"content"`
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
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "title/content required")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	item := models.WBroadcast{Title: strings.TrimSpace(req.Title), Content: strings.TrimSpace(req.Content), Status: 1, Sort: 0}
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
	if err := bc.svc.CreateBroadcast(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, item)
}

// UpdateBroadcast 更新Broadcast。
func (bc *BizConfigController) UpdateBroadcast(c *gin.Context) {
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
	var req struct {
		// 处理当前语句逻辑。
		Title *string `json:"title"`
		// 处理当前语句逻辑。
		Content *string `json:"content"`
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
	// 定义并初始化当前变量。
	updates := map[string]interface{}{}
	// 判断条件并进入对应分支逻辑。
	if req.Title != nil {
		// 更新当前变量或字段值。
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Content != nil {
		// 更新当前变量或字段值。
		updates["content"] = strings.TrimSpace(*req.Content)
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
	if len(updates) == 0 {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.UpdateBroadcast(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteBroadcast 删除Broadcast。
func (bc *BizConfigController) DeleteBroadcast(c *gin.Context) {
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
	if err := bc.svc.DeleteBroadcast(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{"id": id})
}
