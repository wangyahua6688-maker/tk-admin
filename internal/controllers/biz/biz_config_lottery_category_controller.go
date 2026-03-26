package biz

import (
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/models"

	"github.com/gin-gonic/gin"
)

// -------------------- Lottery Category --------------------

func (bc *BizConfigController) ListLotteryCategories(c *gin.Context) {
	// 定义并初始化当前变量。
	keyword := strings.TrimSpace(c.Query("keyword"))
	// 定义并初始化当前变量。
	items, err := bc.svc.ListLotteryCategories(c.Request.Context(), keyword)
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

// CreateLotteryCategory 创建LotteryCategory。
func (bc *BizConfigController) CreateLotteryCategory(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		CategoryKey string `json:"category_key"`
		// 处理当前语句逻辑。
		Name string `json:"name"`
		// 处理当前语句逻辑。
		SearchKeywords string `json:"search_keywords"`
		// 处理当前语句逻辑。
		ShowOnHome *int8 `json:"show_on_home"`
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

	// 定义并初始化当前变量。
	item := models.WLotteryCategory{
		// 调用strings.TrimSpace完成当前处理。
		CategoryKey: strings.TrimSpace(req.CategoryKey),
		// 调用strings.TrimSpace完成当前处理。
		Name: strings.TrimSpace(req.Name),
		// 调用strings.TrimSpace完成当前处理。
		SearchKeywords: strings.TrimSpace(req.SearchKeywords),
		// 处理当前语句逻辑。
		ShowOnHome: 1,
		// 处理当前语句逻辑。
		Status: 1,
		// 处理当前语句逻辑。
		Sort: 0,
	}
	// 判断条件并进入对应分支逻辑。
	if item.CategoryKey == "" || item.Name == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "category_key/name required")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if req.ShowOnHome != nil {
		// 更新当前变量或字段值。
		item.ShowOnHome = *req.ShowOnHome
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
	if err := bc.svc.CreateLotteryCategory(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateLotteryCategory 更新LotteryCategory。
func (bc *BizConfigController) UpdateLotteryCategory(c *gin.Context) {
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
		CategoryKey *string `json:"category_key"`
		// 处理当前语句逻辑。
		Name *string `json:"name"`
		// 处理当前语句逻辑。
		SearchKeywords *string `json:"search_keywords"`
		// 处理当前语句逻辑。
		ShowOnHome *int8 `json:"show_on_home"`
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

	// 定义并初始化当前变量。
	updates := map[string]interface{}{}
	// 判断条件并进入对应分支逻辑。
	if req.CategoryKey != nil {
		// 更新当前变量或字段值。
		updates["category_key"] = strings.TrimSpace(*req.CategoryKey)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Name != nil {
		// 更新当前变量或字段值。
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	// 判断条件并进入对应分支逻辑。
	if req.SearchKeywords != nil {
		// 更新当前变量或字段值。
		updates["search_keywords"] = strings.TrimSpace(*req.SearchKeywords)
	}
	// 判断条件并进入对应分支逻辑。
	if req.ShowOnHome != nil {
		// 更新当前变量或字段值。
		updates["show_on_home"] = *req.ShowOnHome
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
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "empty updates")
		// 返回当前处理结果。
		return
	}

	// 判断条件并进入对应分支逻辑。
	if keyRaw, ok := updates["category_key"]; ok && strings.TrimSpace(keyRaw.(string)) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "category_key cannot be empty")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if nameRaw, ok := updates["name"]; ok && strings.TrimSpace(nameRaw.(string)) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "name cannot be empty")
		// 返回当前处理结果。
		return
	}

	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.UpdateLotteryCategory(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteLotteryCategory 删除LotteryCategory。
func (bc *BizConfigController) DeleteLotteryCategory(c *gin.Context) {
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
	if err := bc.svc.DeleteLotteryCategory(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}
