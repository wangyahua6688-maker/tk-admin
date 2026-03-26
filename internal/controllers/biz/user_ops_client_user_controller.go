package biz

import (
	"strings"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin-full/internal/constants"
	"go-admin-full/internal/models"

	"github.com/gin-gonic/gin"
)

// -------------------- 客户端用户 --------------------

func (uc *UserOpsController) ListClientUsers(c *gin.Context) {
	// 定义并初始化当前变量。
	userType := strings.TrimSpace(c.Query("user_type"))
	// 定义并初始化当前变量。
	items, err := uc.svc.ListClientUsers(c.Request.Context(), userType, 300)
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

// CreateClientUser 创建ClientUser。
func (uc *UserOpsController) CreateClientUser(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Username string `json:"username"`
		// 处理当前语句逻辑。
		Nickname string `json:"nickname"`
		// 处理当前语句逻辑。
		Avatar string `json:"avatar"`
		// 处理当前语句逻辑。
		UserType string `json:"user_type"`
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
	if strings.TrimSpace(req.Username) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "username required")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	userType := normalizeUserType(req.UserType)

	// 定义并初始化当前变量。
	item := models.WUser{
		// 调用strings.TrimSpace完成当前处理。
		Username: strings.TrimSpace(req.Username),
		// 调用strings.TrimSpace完成当前处理。
		Nickname: strings.TrimSpace(req.Nickname),
		// 调用strings.TrimSpace完成当前处理。
		Avatar: strings.TrimSpace(req.Avatar),
		// 处理当前语句逻辑。
		UserType: userType,
		// 处理当前语句逻辑。
		Status: 1,
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		item.Status = *req.Status
	}

	// 判断条件并进入对应分支逻辑。
	if err := uc.svc.CreateClientUser(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateClientUser 更新ClientUser。
func (uc *UserOpsController) UpdateClientUser(c *gin.Context) {
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
		Nickname *string `json:"nickname"`
		// 处理当前语句逻辑。
		Avatar *string `json:"avatar"`
		// 处理当前语句逻辑。
		UserType *string `json:"user_type"`
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
	if req.Nickname != nil {
		// 更新当前变量或字段值。
		updates["nickname"] = strings.TrimSpace(*req.Nickname)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Avatar != nil {
		// 更新当前变量或字段值。
		updates["avatar"] = strings.TrimSpace(*req.Avatar)
	}
	// 判断条件并进入对应分支逻辑。
	if req.UserType != nil {
		// 更新当前变量或字段值。
		updates["user_type"] = normalizeUserType(*req.UserType)
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
	if err := uc.svc.UpdateClientUser(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteClientUser 删除ClientUser。
func (uc *UserOpsController) DeleteClientUser(c *gin.Context) {
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
	if err := uc.svc.DeleteClientUser(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// normalizeUserType 处理normalizeUserType相关逻辑。
func normalizeUserType(v string) string {
	// 根据表达式进入多分支处理。
	switch strings.TrimSpace(v) {
	case "official", "robot", "natural":
		// 返回当前处理结果。
		return strings.TrimSpace(v)
	default:
		return "natural"
	}
}
