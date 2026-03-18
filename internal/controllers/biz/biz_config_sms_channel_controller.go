package biz

import (
	"strings"

	"go-admin-full/internal/constants"
	"go-admin-full/internal/models"
	commonresp "tk-common/utils/httpresp"

	"github.com/gin-gonic/gin"
)

// -------------------- 短信通道配置 --------------------

// ListSMSChannels 查询短信通道配置列表。
func (bc *BizConfigController) ListSMSChannels(c *gin.Context) {
	// 按状态筛选时仅接收 0/1，其它值按“全部”处理。
	status := strings.TrimSpace(c.Query("status"))

	// 执行查询并返回。
	items, err := bc.svc.ListSMSChannels(c.Request.Context(), status, 200)
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

// CreateSMSChannel 新增短信通道配置。
func (bc *BizConfigController) CreateSMSChannel(c *gin.Context) {
	// 声明请求结构，避免 map 直接更新带来的脏字段风险。
	var req struct {
		// 处理当前语句逻辑。
		Provider string `json:"provider"`
		// 处理当前语句逻辑。
		ChannelName string `json:"channel_name"`
		// 处理当前语句逻辑。
		AccessKey string `json:"access_key"`
		// 处理当前语句逻辑。
		AccessSecret string `json:"access_secret"`
		// 处理当前语句逻辑。
		Endpoint string `json:"endpoint"`
		// 处理当前语句逻辑。
		SignName string `json:"sign_name"`
		// 处理当前语句逻辑。
		TemplateCodeLogin string `json:"template_code_login"`
		// 处理当前语句逻辑。
		TemplateCodeRegister string `json:"template_code_register"`
		// 处理当前语句逻辑。
		DailyLimit *int `json:"daily_limit"`
		// 处理当前语句逻辑。
		MinuteLimit *int `json:"minute_limit"`
		// 处理当前语句逻辑。
		CodeTTLSeconds *int `json:"code_ttl_seconds"`
		// 处理当前语句逻辑。
		MockMode *int8 `json:"mock_mode"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
	}

	// 绑定请求体并进行基础校验。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}

	// 归一化必填字段。
	provider := strings.TrimSpace(req.Provider)
	// 定义并初始化当前变量。
	channelName := strings.TrimSpace(req.ChannelName)
	// 判断条件并进入对应分支逻辑。
	if provider == "" {
		// 更新当前变量或字段值。
		provider = "custom"
	}
	// 判断条件并进入对应分支逻辑。
	if channelName == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "channel_name required")
		// 返回当前处理结果。
		return
	}

	// 组装模型并设置默认值。
	item := models.WSMSChannel{
		// 处理当前语句逻辑。
		Provider: provider,
		// 处理当前语句逻辑。
		ChannelName: channelName,
		// 调用strings.TrimSpace完成当前处理。
		AccessKey: strings.TrimSpace(req.AccessKey),
		// 调用strings.TrimSpace完成当前处理。
		AccessSecret: strings.TrimSpace(req.AccessSecret),
		// 调用strings.TrimSpace完成当前处理。
		Endpoint: strings.TrimSpace(req.Endpoint),
		// 调用strings.TrimSpace完成当前处理。
		SignName: strings.TrimSpace(req.SignName),
		// 调用strings.TrimSpace完成当前处理。
		TemplateCodeLogin: strings.TrimSpace(req.TemplateCodeLogin),
		// 调用strings.TrimSpace完成当前处理。
		TemplateCodeRegister: strings.TrimSpace(req.TemplateCodeRegister),
		// 处理当前语句逻辑。
		DailyLimit: 20,
		// 处理当前语句逻辑。
		MinuteLimit: 1,
		// 处理当前语句逻辑。
		CodeTTLSeconds: 300,
		// 处理当前语句逻辑。
		MockMode: 1,
		// 处理当前语句逻辑。
		Status: 1,
	}
	// 判断条件并进入对应分支逻辑。
	if req.DailyLimit != nil && *req.DailyLimit > 0 {
		// 更新当前变量或字段值。
		item.DailyLimit = *req.DailyLimit
	}
	// 判断条件并进入对应分支逻辑。
	if req.MinuteLimit != nil && *req.MinuteLimit > 0 {
		// 更新当前变量或字段值。
		item.MinuteLimit = *req.MinuteLimit
	}
	// 判断条件并进入对应分支逻辑。
	if req.CodeTTLSeconds != nil && *req.CodeTTLSeconds > 0 {
		// 更新当前变量或字段值。
		item.CodeTTLSeconds = *req.CodeTTLSeconds
	}
	// 判断条件并进入对应分支逻辑。
	if req.MockMode != nil {
		// 更新当前变量或字段值。
		item.MockMode = *req.MockMode
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		item.Status = *req.Status
	}

	// 落库并返回。
	if err := bc.svc.CreateSMSChannel(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateSMSChannel 更新短信通道配置。
func (bc *BizConfigController) UpdateSMSChannel(c *gin.Context) {
	// 解析路由 ID。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid id")
		// 返回当前处理结果。
		return
	}

	// 定义可更新字段，禁止前端透传任意字段。
	var req struct {
		// 处理当前语句逻辑。
		Provider *string `json:"provider"`
		// 处理当前语句逻辑。
		ChannelName *string `json:"channel_name"`
		// 处理当前语句逻辑。
		AccessKey *string `json:"access_key"`
		// 处理当前语句逻辑。
		AccessSecret *string `json:"access_secret"`
		// 处理当前语句逻辑。
		Endpoint *string `json:"endpoint"`
		// 处理当前语句逻辑。
		SignName *string `json:"sign_name"`
		// 处理当前语句逻辑。
		TemplateCodeLogin *string `json:"template_code_login"`
		// 处理当前语句逻辑。
		TemplateCodeRegister *string `json:"template_code_register"`
		// 处理当前语句逻辑。
		DailyLimit *int `json:"daily_limit"`
		// 处理当前语句逻辑。
		MinuteLimit *int `json:"minute_limit"`
		// 处理当前语句逻辑。
		CodeTTLSeconds *int `json:"code_ttl_seconds"`
		// 处理当前语句逻辑。
		MockMode *int8 `json:"mock_mode"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
	}

	// 解析请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}

	// 构建更新映射。
	updates := make(map[string]interface{})
	// 判断条件并进入对应分支逻辑。
	if req.Provider != nil {
		// 定义并初始化当前变量。
		provider := strings.TrimSpace(*req.Provider)
		// 判断条件并进入对应分支逻辑。
		if provider == "" {
			// 更新当前变量或字段值。
			provider = "custom"
		}
		// 更新当前变量或字段值。
		updates["provider"] = provider
	}
	// 判断条件并进入对应分支逻辑。
	if req.ChannelName != nil {
		// 定义并初始化当前变量。
		channelName := strings.TrimSpace(*req.ChannelName)
		// 判断条件并进入对应分支逻辑。
		if channelName == "" {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "channel_name required")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["channel_name"] = channelName
	}
	// 判断条件并进入对应分支逻辑。
	if req.AccessKey != nil {
		// 更新当前变量或字段值。
		updates["access_key"] = strings.TrimSpace(*req.AccessKey)
	}
	// 判断条件并进入对应分支逻辑。
	if req.AccessSecret != nil {
		// 更新当前变量或字段值。
		updates["access_secret"] = strings.TrimSpace(*req.AccessSecret)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Endpoint != nil {
		// 更新当前变量或字段值。
		updates["endpoint"] = strings.TrimSpace(*req.Endpoint)
	}
	// 判断条件并进入对应分支逻辑。
	if req.SignName != nil {
		// 更新当前变量或字段值。
		updates["sign_name"] = strings.TrimSpace(*req.SignName)
	}
	// 判断条件并进入对应分支逻辑。
	if req.TemplateCodeLogin != nil {
		// 更新当前变量或字段值。
		updates["template_code_login"] = strings.TrimSpace(*req.TemplateCodeLogin)
	}
	// 判断条件并进入对应分支逻辑。
	if req.TemplateCodeRegister != nil {
		// 更新当前变量或字段值。
		updates["template_code_register"] = strings.TrimSpace(*req.TemplateCodeRegister)
	}
	// 判断条件并进入对应分支逻辑。
	if req.DailyLimit != nil {
		// 判断条件并进入对应分支逻辑。
		if *req.DailyLimit <= 0 {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "daily_limit must > 0")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["daily_limit"] = *req.DailyLimit
	}
	// 判断条件并进入对应分支逻辑。
	if req.MinuteLimit != nil {
		// 判断条件并进入对应分支逻辑。
		if *req.MinuteLimit <= 0 {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "minute_limit must > 0")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["minute_limit"] = *req.MinuteLimit
	}
	// 判断条件并进入对应分支逻辑。
	if req.CodeTTLSeconds != nil {
		// 判断条件并进入对应分支逻辑。
		if *req.CodeTTLSeconds <= 0 {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "code_ttl_seconds must > 0")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["code_ttl_seconds"] = *req.CodeTTLSeconds
	}
	// 判断条件并进入对应分支逻辑。
	if req.MockMode != nil {
		// 更新当前变量或字段值。
		updates["mock_mode"] = *req.MockMode
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

	// 执行更新并返回。
	if err := bc.svc.UpdateSMSChannel(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteSMSChannel 删除短信通道配置。
func (bc *BizConfigController) DeleteSMSChannel(c *gin.Context) {
	// 解析 ID 并执行删除。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid id")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.DeleteSMSChannel(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}
