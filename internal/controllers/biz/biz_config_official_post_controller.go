package biz

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 官方发帖 --------------------

func (bc *BizConfigController) ListOfficialPosts(c *gin.Context) {
	// 定义并初始化当前变量。
	items, err := bc.svc.ListOfficialPosts(c.Request.Context(), 200)
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

// CreateOfficialPost 创建OfficialPost。
func (bc *BizConfigController) CreateOfficialPost(c *gin.Context) {
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
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(req.Title) == "" {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "title required")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if req.UserID == 0 {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "user_id required")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if !bc.svc.IsUserTypes(c.Request.Context(), req.UserID, "official") {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, http.StatusBadRequest, "user_id must be official account")
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
		IsOfficial: 1,
		// 处理当前语句逻辑。
		Status: 1,
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		item.Status = *req.Status
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.CreateOfficialPost(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, item)
}

// UpdateOfficialPost 更新OfficialPost。
func (bc *BizConfigController) UpdateOfficialPost(c *gin.Context) {
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
	// 更新当前变量或字段值。
	req["is_official"] = int8(1)
	// 判断条件并进入对应分支逻辑。
	if rawUserID, ok := req["user_id"]; ok {
		// 定义并初始化当前变量。
		userID := toUint(rawUserID)
		// 判断条件并进入对应分支逻辑。
		if userID == 0 || !bc.svc.IsUserTypes(c.Request.Context(), userID, "official") {
			// 调用utils.JSONError完成当前处理。
			utils.JSONError(c, http.StatusBadRequest, "user_id must be official account")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		req["user_id"] = userID
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.UpdateOfficialPost(c.Request.Context(), id, req); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteOfficialPost 删除OfficialPost。
func (bc *BizConfigController) DeleteOfficialPost(c *gin.Context) {
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
	if err := bc.svc.DeleteOfficialPost(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		utils.JSONError(c, 500, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	utils.JSONOK(c, gin.H{"id": id})
}

// toUint 处理toUint相关逻辑。
func toUint(raw interface{}) uint {
	// 根据表达式进入多分支处理。
	switch v := raw.(type) {
	case float64:
		// 判断条件并进入对应分支逻辑。
		if v > 0 {
			// 返回当前处理结果。
			return uint(v)
		}
	case int:
		// 判断条件并进入对应分支逻辑。
		if v > 0 {
			// 返回当前处理结果。
			return uint(v)
		}
	case int64:
		// 判断条件并进入对应分支逻辑。
		if v > 0 {
			// 返回当前处理结果。
			return uint(v)
		}
	case string:
		// 更新当前变量或字段值。
		v = strings.TrimSpace(v)
		// 判断条件并进入对应分支逻辑。
		if v == "" {
			// 返回当前处理结果。
			return 0
		}
		// 声明当前变量。
		var n uint
		// 更新当前变量或字段值。
		_, _ = fmt.Sscanf(v, "%d", &n)
		// 返回当前处理结果。
		return n
	}
	// 返回当前处理结果。
	return 0
}
