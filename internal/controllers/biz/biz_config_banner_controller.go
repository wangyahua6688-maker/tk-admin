package biz

import (
	"sort"
	"strings"

	"go-admin-full/internal/constants"
	"go-admin-full/internal/models"
	commonresp "tk-common/utils/httpresp"

	"github.com/gin-gonic/gin"
)

// -------------------- Banner --------------------

func (bc *BizConfigController) ListBanners(c *gin.Context) {
	// 定义并初始化当前变量。
	bannerType := strings.TrimSpace(c.Query("type"))
	// 定义并初始化当前变量。
	items, err := bc.svc.ListBanners(c.Request.Context(), bannerType, 300)
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

// CreateBanner 创建Banner。
func (bc *BizConfigController) CreateBanner(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Title string `json:"title"`
		// 处理当前语句逻辑。
		ImageURL string `json:"image_url"`
		// 处理当前语句逻辑。
		LinkURL string `json:"link_url"`
		// 处理当前语句逻辑。
		Type string `json:"type"`
		// 处理当前语句逻辑。
		Position string `json:"position"`
		// 处理当前语句逻辑。
		Positions []string `json:"positions"`
		// 处理当前语句逻辑。
		JumpType string `json:"jump_type"`
		// 处理当前语句逻辑。
		JumpPostID uint `json:"jump_post_id"`
		// 处理当前语句逻辑。
		JumpURL string `json:"jump_url"`
		// 处理当前语句逻辑。
		ContentHTML string `json:"content_html"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
		// 处理当前语句逻辑。
		Sort *int `json:"sort"`
		// 处理当前语句逻辑。
		StartAt *string `json:"start_at"`
		// 处理当前语句逻辑。
		EndAt *string `json:"end_at"`
	}
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	positions := normalizePositions(req.Positions, req.Position)
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.ImageURL) == "" || strings.TrimSpace(req.Type) == "" || len(positions) == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "title/image_url/type/positions required")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	linkURL, err := normalizeSafeURL(req.LinkURL, true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid link_url")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	jumpURL, err := normalizeSafeURL(req.JumpURL, true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid jump_url")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	jumpType := normalizeJumpType(req.JumpType)
	// 判断条件并进入对应分支逻辑。
	if jumpType == "post" && req.JumpPostID == 0 {
		// 更新当前变量或字段值。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "jump_post_id required when jump_type=post")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if jumpType == "external" && strings.TrimSpace(req.JumpURL) == "" {
		// 更新当前变量或字段值。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "jump_url required when jump_type=external")
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if jumpType == "post" {
		// 定义并初始化当前变量。
		ok, err := bc.svc.IsPostExists(c.Request.Context(), req.JumpPostID)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
			// 返回当前处理结果。
			return
		}
		// 判断条件并进入对应分支逻辑。
		if !ok {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "jump_post_id not found")
			// 返回当前处理结果。
			return
		}
	}

	// 定义并初始化当前变量。
	item := models.WBanner{
		// 调用strings.TrimSpace完成当前处理。
		Title: strings.TrimSpace(req.Title),
		// 调用strings.TrimSpace完成当前处理。
		ImageURL: strings.TrimSpace(req.ImageURL),
		// 处理当前语句逻辑。
		LinkURL: linkURL,
		// 调用strings.TrimSpace完成当前处理。
		Type: strings.TrimSpace(req.Type),
		// 处理当前语句逻辑。
		Position: positions[0],
		// 调用strings.Join完成当前处理。
		Positions: strings.Join(positions, ","),
		// 处理当前语句逻辑。
		JumpType: jumpType,
		// 处理当前语句逻辑。
		JumpPostID: req.JumpPostID,
		// 处理当前语句逻辑。
		JumpURL: jumpURL,
		// 处理当前语句逻辑。
		ContentHTML: req.ContentHTML,
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
	if jumpType == "external" {
		// 更新当前变量或字段值。
		item.LinkURL = jumpURL
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
	if err := bc.svc.CreateBanner(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateBanner 更新Banner。
func (bc *BizConfigController) UpdateBanner(c *gin.Context) {
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
		Title *string `json:"title"`
		// 处理当前语句逻辑。
		ImageURL *string `json:"image_url"`
		// 处理当前语句逻辑。
		LinkURL *string `json:"link_url"`
		// 处理当前语句逻辑。
		Type *string `json:"type"`
		// 处理当前语句逻辑。
		Position *string `json:"position"`
		// 处理当前语句逻辑。
		Positions []string `json:"positions"`
		// 处理当前语句逻辑。
		JumpType *string `json:"jump_type"`
		// 处理当前语句逻辑。
		JumpPostID *uint `json:"jump_post_id"`
		// 处理当前语句逻辑。
		JumpURL *string `json:"jump_url"`
		// 处理当前语句逻辑。
		ContentHTML *string `json:"content_html"`
		// 处理当前语句逻辑。
		Status *int8 `json:"status"`
		// 处理当前语句逻辑。
		Sort *int `json:"sort"`
		// 处理当前语句逻辑。
		StartAt *string `json:"start_at"`
		// 处理当前语句逻辑。
		EndAt *string `json:"end_at"`
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
	if req.Title != nil {
		// 更新当前变量或字段值。
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	// 判断条件并进入对应分支逻辑。
	if req.ImageURL != nil {
		// 更新当前变量或字段值。
		updates["image_url"] = strings.TrimSpace(*req.ImageURL)
	}
	// 判断条件并进入对应分支逻辑。
	if req.LinkURL != nil {
		// 定义并初始化当前变量。
		linkURL, linkErr := normalizeSafeURL(*req.LinkURL, true)
		// 判断条件并进入对应分支逻辑。
		if linkErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid link_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["link_url"] = linkURL
	}
	// 判断条件并进入对应分支逻辑。
	if req.Type != nil {
		// 更新当前变量或字段值。
		updates["type"] = strings.TrimSpace(*req.Type)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Position != nil {
		// 更新当前变量或字段值。
		updates["position"] = strings.TrimSpace(*req.Position)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Positions != nil {
		// 定义并初始化当前变量。
		positions := normalizePositions(req.Positions, "")
		// 判断条件并进入对应分支逻辑。
		if len(positions) > 0 {
			// 更新当前变量或字段值。
			updates["positions"] = strings.Join(positions, ",")
			// 更新当前变量或字段值。
			updates["position"] = positions[0]
		}
	}
	// 判断条件并进入对应分支逻辑。
	if req.JumpType != nil {
		// 更新当前变量或字段值。
		updates["jump_type"] = normalizeJumpType(*req.JumpType)
	}
	// 判断条件并进入对应分支逻辑。
	if req.JumpPostID != nil {
		// 更新当前变量或字段值。
		updates["jump_post_id"] = *req.JumpPostID
	}
	// 判断条件并进入对应分支逻辑。
	if req.JumpURL != nil {
		// 定义并初始化当前变量。
		jumpURL, jumpErr := normalizeSafeURL(*req.JumpURL, true)
		// 判断条件并进入对应分支逻辑。
		if jumpErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid jump_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["jump_url"] = jumpURL
	}
	// 判断条件并进入对应分支逻辑。
	if req.ContentHTML != nil {
		// 更新当前变量或字段值。
		updates["content_html"] = *req.ContentHTML
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
	if link, ok := updates["jump_url"]; ok {
		// 判断条件并进入对应分支逻辑。
		if jt, hasJT := updates["jump_type"]; hasJT && jt == "external" {
			// 更新当前变量或字段值。
			updates["link_url"] = link
		}
	}
	// 判断条件并进入对应分支逻辑。
	if len(updates) == 0 {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "empty updates")
		// 返回当前处理结果。
		return
	}

	// 判断条件并进入对应分支逻辑。
	if err := bc.svc.UpdateBanner(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteBanner 删除Banner。
func (bc *BizConfigController) DeleteBanner(c *gin.Context) {
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
	if err := bc.svc.DeleteBanner(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// normalizePositions 处理normalizePositions相关逻辑。
func normalizePositions(positions []string, fallback string) []string {
	// 定义并初始化当前变量。
	out := make([]string, 0)
	// 定义并初始化当前变量。
	seen := map[string]struct{}{}
	// 定义并初始化当前变量。
	appendPos := func(v string) {
		// 定义并初始化当前变量。
		pos := strings.TrimSpace(v)
		// 判断条件并进入对应分支逻辑。
		if pos == "" {
			// 返回当前处理结果。
			return
		}
		// 判断条件并进入对应分支逻辑。
		if _, ok := seen[pos]; ok {
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		seen[pos] = struct{}{}
		// 更新当前变量或字段值。
		out = append(out, pos)
	}
	// 循环处理当前数据集合。
	for _, p := range positions {
		// 调用appendPos完成当前处理。
		appendPos(p)
	}
	// 判断条件并进入对应分支逻辑。
	if len(out) == 0 {
		// 调用appendPos完成当前处理。
		appendPos(fallback)
	}
	// 调用sort.Strings完成当前处理。
	sort.Strings(out)
	// 返回当前处理结果。
	return out
}

// normalizeJumpType 处理normalizeJumpType相关逻辑。
func normalizeJumpType(v string) string {
	// 根据表达式进入多分支处理。
	switch strings.TrimSpace(v) {
	case "post", "external", "custom", "none":
		// 返回当前处理结果。
		return strings.TrimSpace(v)
	default:
		return "none"
	}
}
