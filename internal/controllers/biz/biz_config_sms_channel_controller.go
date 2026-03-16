package biz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- 短信通道配置 --------------------

// ListSMSChannels 查询短信通道配置列表。
func (bc *BizConfigController) ListSMSChannels(c *gin.Context) {
	// 按状态筛选时仅接收 0/1，其它值按“全部”处理。
	status := strings.TrimSpace(c.Query("status"))

	// 执行查询并返回。
	items, err := bc.svc.ListSMSChannels(c.Request.Context(), status, 200)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

// CreateSMSChannel 新增短信通道配置。
func (bc *BizConfigController) CreateSMSChannel(c *gin.Context) {
	// 声明请求结构，避免 map 直接更新带来的脏字段风险。
	var req struct {
		Provider             string `json:"provider"`
		ChannelName          string `json:"channel_name"`
		AccessKey            string `json:"access_key"`
		AccessSecret         string `json:"access_secret"`
		Endpoint             string `json:"endpoint"`
		SignName             string `json:"sign_name"`
		TemplateCodeLogin    string `json:"template_code_login"`
		TemplateCodeRegister string `json:"template_code_register"`
		DailyLimit           *int   `json:"daily_limit"`
		MinuteLimit          *int   `json:"minute_limit"`
		CodeTTLSeconds       *int   `json:"code_ttl_seconds"`
		MockMode             *int8  `json:"mock_mode"`
		Status               *int8  `json:"status"`
	}

	// 绑定请求体并进行基础校验。
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	// 归一化必填字段。
	provider := strings.TrimSpace(req.Provider)
	channelName := strings.TrimSpace(req.ChannelName)
	if provider == "" {
		provider = "custom"
	}
	if channelName == "" {
		utils.JSONError(c, http.StatusBadRequest, "channel_name required")
		return
	}

	// 组装模型并设置默认值。
	item := models.WSMSChannel{
		Provider:             provider,
		ChannelName:          channelName,
		AccessKey:            strings.TrimSpace(req.AccessKey),
		AccessSecret:         strings.TrimSpace(req.AccessSecret),
		Endpoint:             strings.TrimSpace(req.Endpoint),
		SignName:             strings.TrimSpace(req.SignName),
		TemplateCodeLogin:    strings.TrimSpace(req.TemplateCodeLogin),
		TemplateCodeRegister: strings.TrimSpace(req.TemplateCodeRegister),
		DailyLimit:           20,
		MinuteLimit:          1,
		CodeTTLSeconds:       300,
		MockMode:             1,
		Status:               1,
	}
	if req.DailyLimit != nil && *req.DailyLimit > 0 {
		item.DailyLimit = *req.DailyLimit
	}
	if req.MinuteLimit != nil && *req.MinuteLimit > 0 {
		item.MinuteLimit = *req.MinuteLimit
	}
	if req.CodeTTLSeconds != nil && *req.CodeTTLSeconds > 0 {
		item.CodeTTLSeconds = *req.CodeTTLSeconds
	}
	if req.MockMode != nil {
		item.MockMode = *req.MockMode
	}
	if req.Status != nil {
		item.Status = *req.Status
	}

	// 落库并返回。
	if err := bc.svc.CreateSMSChannel(c.Request.Context(), &item); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

// UpdateSMSChannel 更新短信通道配置。
func (bc *BizConfigController) UpdateSMSChannel(c *gin.Context) {
	// 解析路由 ID。
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	// 定义可更新字段，禁止前端透传任意字段。
	var req struct {
		Provider             *string `json:"provider"`
		ChannelName          *string `json:"channel_name"`
		AccessKey            *string `json:"access_key"`
		AccessSecret         *string `json:"access_secret"`
		Endpoint             *string `json:"endpoint"`
		SignName             *string `json:"sign_name"`
		TemplateCodeLogin    *string `json:"template_code_login"`
		TemplateCodeRegister *string `json:"template_code_register"`
		DailyLimit           *int    `json:"daily_limit"`
		MinuteLimit          *int    `json:"minute_limit"`
		CodeTTLSeconds       *int    `json:"code_ttl_seconds"`
		MockMode             *int8   `json:"mock_mode"`
		Status               *int8   `json:"status"`
	}

	// 解析请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	// 构建更新映射。
	updates := make(map[string]interface{})
	if req.Provider != nil {
		provider := strings.TrimSpace(*req.Provider)
		if provider == "" {
			provider = "custom"
		}
		updates["provider"] = provider
	}
	if req.ChannelName != nil {
		channelName := strings.TrimSpace(*req.ChannelName)
		if channelName == "" {
			utils.JSONError(c, http.StatusBadRequest, "channel_name required")
			return
		}
		updates["channel_name"] = channelName
	}
	if req.AccessKey != nil {
		updates["access_key"] = strings.TrimSpace(*req.AccessKey)
	}
	if req.AccessSecret != nil {
		updates["access_secret"] = strings.TrimSpace(*req.AccessSecret)
	}
	if req.Endpoint != nil {
		updates["endpoint"] = strings.TrimSpace(*req.Endpoint)
	}
	if req.SignName != nil {
		updates["sign_name"] = strings.TrimSpace(*req.SignName)
	}
	if req.TemplateCodeLogin != nil {
		updates["template_code_login"] = strings.TrimSpace(*req.TemplateCodeLogin)
	}
	if req.TemplateCodeRegister != nil {
		updates["template_code_register"] = strings.TrimSpace(*req.TemplateCodeRegister)
	}
	if req.DailyLimit != nil {
		if *req.DailyLimit <= 0 {
			utils.JSONError(c, http.StatusBadRequest, "daily_limit must > 0")
			return
		}
		updates["daily_limit"] = *req.DailyLimit
	}
	if req.MinuteLimit != nil {
		if *req.MinuteLimit <= 0 {
			utils.JSONError(c, http.StatusBadRequest, "minute_limit must > 0")
			return
		}
		updates["minute_limit"] = *req.MinuteLimit
	}
	if req.CodeTTLSeconds != nil {
		if *req.CodeTTLSeconds <= 0 {
			utils.JSONError(c, http.StatusBadRequest, "code_ttl_seconds must > 0")
			return
		}
		updates["code_ttl_seconds"] = *req.CodeTTLSeconds
	}
	if req.MockMode != nil {
		updates["mock_mode"] = *req.MockMode
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}

	// 执行更新并返回。
	if err := bc.svc.UpdateSMSChannel(c.Request.Context(), id, updates); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteSMSChannel 删除短信通道配置。
func (bc *BizConfigController) DeleteSMSChannel(c *gin.Context) {
	// 解析 ID 并执行删除。
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.svc.DeleteSMSChannel(c.Request.Context(), id); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}
