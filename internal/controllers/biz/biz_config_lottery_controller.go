package biz

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"go-admin-full/internal/constants"
	"go-admin-full/internal/models"
	commonresp "tk-common/utils/httpresp"

	"github.com/gin-gonic/gin"
)

// -------------------- Special Lottery --------------------

// ListSpecialLotteries 查询彩种配置列表。
func (bc *BizConfigController) ListSpecialLotteries(c *gin.Context) {
	// 固定按 sort + id 排序，确保前端按钮顺序稳定。
	items, err := bc.svc.ListSpecialLotteries(c.Request.Context(), 200)
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

// CreateSpecialLottery 新增彩种配置。
func (bc *BizConfigController) CreateSpecialLottery(c *gin.Context) {
	// 声明当前变量。
	var req struct {
		// 处理当前语句逻辑。
		Name string `json:"name"`
		// 处理当前语句逻辑。
		Code string `json:"code"`
		// 处理当前语句逻辑。
		CurrentIssue string `json:"current_issue"`
		// 处理当前语句逻辑。
		NextDrawAt string `json:"next_draw_at"`
		// 处理当前语句逻辑。
		LiveEnabled *int8 `json:"live_enabled"`
		// 处理当前语句逻辑。
		LiveStatus string `json:"live_status"`
		// 处理当前语句逻辑。
		LiveStreamURL string `json:"live_stream_url"`
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
	// 名称和编码为必填。
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Code) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "name/code required")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	liveStreamURL, err := normalizeSafeURL(req.LiveStreamURL, true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid live_stream_url")
		// 返回当前处理结果。
		return
	}
	// 下期开奖时间按“每天固定时刻”解析，前端只需要传 HH:mm:ss。
	next, parseErr := parseDailyDrawTime(req.NextDrawAt, lotteryNowInEast8())
	if parseErr != nil {
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid next_draw_at, expected HH:mm:ss")
		return
	}
	// 定义并初始化当前变量。
	item := models.WSpecialLottery{
		// 调用strings.TrimSpace完成当前处理。
		Name: strings.TrimSpace(req.Name),
		// 调用strings.TrimSpace完成当前处理。
		Code: strings.TrimSpace(req.Code),
		// 调用strings.TrimSpace完成当前处理。
		CurrentIssue: strings.TrimSpace(req.CurrentIssue),
		// 处理当前语句逻辑。
		NextDrawAt: next,
		// 处理当前语句逻辑。
		LiveEnabled: 0,
		// 调用strings.TrimSpace完成当前处理。
		LiveStatus: strings.TrimSpace(req.LiveStatus),
		// 处理当前语句逻辑。
		LiveStreamURL: liveStreamURL,
		// 处理当前语句逻辑。
		Status: 1,
		// 处理当前语句逻辑。
		Sort: 0,
	}
	// 直播状态为空时，默认 pending。
	if item.LiveStatus == "" {
		// 更新当前变量或字段值。
		item.LiveStatus = "pending"
	}
	// 判断条件并进入对应分支逻辑。
	if req.LiveEnabled != nil {
		// 更新当前变量或字段值。
		item.LiveEnabled = *req.LiveEnabled
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
	if err := bc.svc.CreateSpecialLottery(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 新增彩种后，首页彩种按钮/开奖区/开奖现场都可能立即受到影响，因此同步清理公开缓存。
	_ = invalidatePublicLotteryCaches(c.Request.Context(), item.ID)
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateSpecialLottery 更新彩种配置。
func (bc *BizConfigController) UpdateSpecialLottery(c *gin.Context) {
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
		Name *string `json:"name"`
		// 处理当前语句逻辑。
		Code *string `json:"code"`
		// 处理当前语句逻辑。
		CurrentIssue *string `json:"current_issue"`
		// 处理当前语句逻辑。
		NextDrawAt *string `json:"next_draw_at"`
		// 处理当前语句逻辑。
		LiveEnabled *int8 `json:"live_enabled"`
		// 处理当前语句逻辑。
		LiveStatus *string `json:"live_status"`
		// 处理当前语句逻辑。
		LiveStreamURL *string `json:"live_stream_url"`
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
	if req.Name != nil {
		// 更新当前变量或字段值。
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Code != nil {
		// 更新当前变量或字段值。
		updates["code"] = strings.TrimSpace(*req.Code)
	}
	// 判断条件并进入对应分支逻辑。
	if req.CurrentIssue != nil {
		// 更新当前变量或字段值。
		updates["current_issue"] = strings.TrimSpace(*req.CurrentIssue)
	}
	// 判断条件并进入对应分支逻辑。
	if req.NextDrawAt != nil {
		next, parseErr := parseDailyDrawTime(*req.NextDrawAt, lotteryNowInEast8())
		if parseErr != nil {
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid next_draw_at, expected HH:mm:ss")
			return
		}
		// 更新当前变量或字段值。
		updates["next_draw_at"] = next
	}
	// 判断条件并进入对应分支逻辑。
	if req.LiveEnabled != nil {
		// 更新当前变量或字段值。
		updates["live_enabled"] = *req.LiveEnabled
	}
	// 判断条件并进入对应分支逻辑。
	if req.LiveStatus != nil {
		// 更新当前变量或字段值。
		updates["live_status"] = strings.TrimSpace(*req.LiveStatus)
	}
	// 判断条件并进入对应分支逻辑。
	if req.LiveStreamURL != nil {
		// 定义并初始化当前变量。
		liveStreamURL, liveErr := normalizeSafeURL(*req.LiveStreamURL, true)
		// 判断条件并进入对应分支逻辑。
		if liveErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid live_stream_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		updates["live_stream_url"] = liveStreamURL
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
	if err := bc.svc.UpdateSpecialLottery(c.Request.Context(), id, updates); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 彩种配置更新后，首页概览与开奖看板缓存必须立刻失效，否则页面会继续显示旧时间。
	_ = invalidatePublicLotteryCaches(c.Request.Context(), id)
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteSpecialLottery 删除彩种配置。
func (bc *BizConfigController) DeleteSpecialLottery(c *gin.Context) {
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
	if err := bc.svc.DeleteSpecialLottery(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 删除彩种后也要清理首页/看板/开奖现场缓存，避免前端继续读取已删除彩种的旧数据。
	_ = invalidatePublicLotteryCaches(c.Request.Context(), id)
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// -------------------- Lottery Info --------------------

// lotteryInfoUpsertRequest 图库内容新增/编辑请求结构。
type lotteryInfoUpsertRequest struct {
	// 处理当前语句逻辑。
	SpecialLotteryID *uint `json:"special_lottery_id"`
	// 处理当前语句逻辑。
	CategoryID *uint `json:"category_id"`
	// 处理当前语句逻辑。
	CategoryTag *string `json:"category_tag"`
	// 处理当前语句逻辑。
	Issue *string `json:"issue"`
	// 处理当前语句逻辑。
	Year *int `json:"year"`
	// 处理当前语句逻辑。
	Title *string `json:"title"`
	// 处理当前语句逻辑。
	CoverImageURL *string `json:"cover_image_url"`
	// 处理当前语句逻辑。
	DetailImageURL *string `json:"detail_image_url"`
	// 处理当前语句逻辑。
	DrawCode *string `json:"draw_code"`
	// 处理当前语句逻辑。
	NormalDrawResult *string `json:"normal_draw_result"`
	// 处理当前语句逻辑。
	SpecialDrawResult *string `json:"special_draw_result"`
	// 处理当前语句逻辑。
	DrawResult *string `json:"draw_result"`
	// 处理当前语句逻辑。
	DrawAt *string `json:"draw_at"`
	// 处理当前语句逻辑。
	PlaybackURL *string `json:"playback_url"`
	// 处理当前语句逻辑。
	LikesCount *int64 `json:"likes_count"`
	// 处理当前语句逻辑。
	CommentCount *int64 `json:"comment_count"`
	// 处理当前语句逻辑。
	FavoriteCount *int64 `json:"favorite_count"`
	// 处理当前语句逻辑。
	ReadCount *int64 `json:"read_count"`
	// 处理当前语句逻辑。
	PollEnabled *int8 `json:"poll_enabled"`
	// 处理当前语句逻辑。
	PollDefaultExpand *int8 `json:"poll_default_expand"`
	// 处理当前语句逻辑。
	RecommendInfoIDs *string `json:"recommend_info_ids"`
	// 处理当前语句逻辑。
	OptionNames *[]string `json:"option_names"`
	// 处理当前语句逻辑。
	IsCurrent *int8 `json:"is_current"`
	// 处理当前语句逻辑。
	Status *int8 `json:"status"`
	// 处理当前语句逻辑。
	Sort *int `json:"sort"`
}

// ListLotteryInfos 查询图库内容列表。
func (bc *BizConfigController) ListLotteryInfos(c *gin.Context) {
	// 图库内容管理优先按更新时间倒序，避免运营修改后找不到记录。
	items, optionNameMap, err := bc.svc.ListLotteryInfosWithOptions(c.Request.Context(), 300)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{
		// 处理当前语句逻辑。
		"items": items,
		// 处理当前语句逻辑。
		"option_names_by_info_id": optionNameMap,
	})
}

// CreateLotteryInfo 新增图库内容。
func (bc *BizConfigController) CreateLotteryInfo(c *gin.Context) {
	// 声明当前变量。
	var req lotteryInfoUpsertRequest
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 校验基础必填项（图库内容不再强制绑定彩种）。
	if req.Issue == nil || strings.TrimSpace(*req.Issue) == "" || req.Title == nil || strings.TrimSpace(*req.Title) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "issue/title required")
		// 返回当前处理结果。
		return
	}
	// 图库内容可独立存在：special_lottery_id 允许为空，落库为 0。
	// 说明：彩种维度只在“开奖区开奖记录（tk_draw_record）”中强约束。
	specialLotteryID := uint(0)
	// 判断条件并进入对应分支逻辑。
	if req.SpecialLotteryID != nil {
		// 更新当前变量或字段值。
		specialLotteryID = *req.SpecialLotteryID
	}
	// 分类必须可落到 tk_lottery_category，禁止手工自由文本分类。
	categoryID, categoryTag, err := bc.svc.ResolveLotteryCategory(c.Request.Context(), req.CategoryID, req.CategoryTag)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, err.Error())
		// 返回当前处理结果。
		return
	}
	// 图库内容不强制录入 6+1 开奖码，仅在传值时做校验。
	normalRaw, specialRaw, mergedRaw := "", "", ""
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(safeString(req.NormalDrawResult)) != "" ||
		// 调用strings.TrimSpace完成当前处理。
		strings.TrimSpace(safeString(req.SpecialDrawResult)) != "" ||
		// 调用strings.TrimSpace完成当前处理。
		strings.TrimSpace(safeString(req.DrawResult)) != "" {
		// 声明当前变量。
		var normalizeErr error
		// 更新当前变量或字段值。
		normalRaw, specialRaw, mergedRaw, normalizeErr = normalizeAndMergeDrawNumbers(
			// 调用safeString完成当前处理。
			safeString(req.NormalDrawResult),
			// 调用safeString完成当前处理。
			safeString(req.SpecialDrawResult),
			// 调用safeString完成当前处理。
			safeString(req.DrawResult),
		)
		// 判断条件并进入对应分支逻辑。
		if normalizeErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, normalizeErr.Error())
			// 返回当前处理结果。
			return
		}
	}
	// 预处理动物竞猜选项：为空时回退到 12 生肖默认选项。
	optionNames := normalizeOptionNames(nil)
	// 判断条件并进入对应分支逻辑。
	if req.OptionNames != nil {
		// 更新当前变量或字段值。
		optionNames = normalizeOptionNames(*req.OptionNames)
	}
	// 判断条件并进入对应分支逻辑。
	if len(optionNames) == 0 {
		// 更新当前变量或字段值。
		optionNames = defaultAnimalOptionNames()
	}
	// 定义并初始化当前变量。
	coverImageURL, err := normalizeSafeURL(safeString(req.CoverImageURL), true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid cover_image_url")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	detailImageURL, err := normalizeSafeURL(safeString(req.DetailImageURL), true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid detail_image_url")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	playbackURL, err := normalizeSafeURL(safeString(req.PlaybackURL), true)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid playback_url")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	now := time.Now()
	// 定义并初始化当前变量。
	item := models.WLotteryInfo{
		// 处理当前语句逻辑。
		SpecialLotteryID: specialLotteryID,
		// 处理当前语句逻辑。
		CategoryID: categoryID,
		// 处理当前语句逻辑。
		CategoryTag: categoryTag,
		// 调用strings.TrimSpace完成当前处理。
		Issue: strings.TrimSpace(*req.Issue),
		// 调用now.Year完成当前处理。
		Year: now.Year(),
		// 调用strings.TrimSpace完成当前处理。
		Title: strings.TrimSpace(*req.Title),
		// 处理当前语句逻辑。
		CoverImageURL: coverImageURL,
		// 处理当前语句逻辑。
		DetailImageURL: detailImageURL,
		// 调用safeString完成当前处理。
		DrawCode: safeString(req.DrawCode),
		// 处理当前语句逻辑。
		NormalDrawResult: normalRaw,
		// 处理当前语句逻辑。
		SpecialDrawResult: specialRaw,
		// 处理当前语句逻辑。
		DrawResult: mergedRaw,
		// 调用parseDateTimeOrDefault完成当前处理。
		DrawAt: parseDateTimeOrDefault(safeString(req.DrawAt), now),
		// 处理当前语句逻辑。
		PlaybackURL: playbackURL,
		// 调用safeInt64完成当前处理。
		LikesCount: safeInt64(req.LikesCount, 0),
		// 调用safeInt64完成当前处理。
		CommentCount: safeInt64(req.CommentCount, 0),
		// 调用safeInt64完成当前处理。
		FavoriteCount: safeInt64(req.FavoriteCount, 0),
		// 调用safeInt64完成当前处理。
		ReadCount: safeInt64(req.ReadCount, 0),
		// 调用safeInt8完成当前处理。
		PollEnabled: safeInt8(req.PollEnabled, 1),
		// 调用safeInt8完成当前处理。
		PollDefaultExpand: safeInt8(req.PollDefaultExpand, 0),
		// 调用safeString完成当前处理。
		RecommendInfoIDs: safeString(req.RecommendInfoIDs),
		// 调用safeInt8完成当前处理。
		IsCurrent: safeInt8(req.IsCurrent, 0),
		// 调用safeInt8完成当前处理。
		Status: safeInt8(req.Status, 1),
		// 调用safeInt完成当前处理。
		Sort: safeInt(req.Sort, 0),
	}
	// 判断条件并进入对应分支逻辑。
	if req.Year != nil && *req.Year > 0 {
		// 更新当前变量或字段值。
		item.Year = *req.Year
	}

	// 当前期约束：同一彩种只能有一条 is_current=1 记录。
	if err := bc.svc.CreateLotteryInfo(c.Request.Context(), &item, optionNames); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateLotteryInfo 编辑图库内容。
func (bc *BizConfigController) UpdateLotteryInfo(c *gin.Context) {
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
	var req lotteryInfoUpsertRequest
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	current, err := bc.svc.GetLotteryInfoByID(c.Request.Context(), id)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizResourceNotFound, "lottery info not found")
		// 返回当前处理结果。
		return
	}

	// 在“当前值”基础上叠加本次更新，保证部分更新不会丢字段。
	next := *current
	// 判断条件并进入对应分支逻辑。
	if req.SpecialLotteryID != nil {
		// 图库内容允许取消彩种绑定（设置为 0）。
		// 前端“图库管理”页已移除彩种选择，这里保留兼容更新能力。
		next.SpecialLotteryID = *req.SpecialLotteryID
	}
	// 判断条件并进入对应分支逻辑。
	if req.Issue != nil {
		// 更新当前变量或字段值。
		next.Issue = strings.TrimSpace(*req.Issue)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Title != nil {
		// 更新当前变量或字段值。
		next.Title = strings.TrimSpace(*req.Title)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Year != nil && *req.Year > 0 {
		// 更新当前变量或字段值。
		next.Year = *req.Year
	}
	// 判断条件并进入对应分支逻辑。
	if req.CoverImageURL != nil {
		// 定义并初始化当前变量。
		coverImageURL, coverErr := normalizeSafeURL(*req.CoverImageURL, true)
		// 判断条件并进入对应分支逻辑。
		if coverErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid cover_image_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		next.CoverImageURL = coverImageURL
	}
	// 判断条件并进入对应分支逻辑。
	if req.DetailImageURL != nil {
		// 定义并初始化当前变量。
		detailImageURL, detailErr := normalizeSafeURL(*req.DetailImageURL, true)
		// 判断条件并进入对应分支逻辑。
		if detailErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid detail_image_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		next.DetailImageURL = detailImageURL
	}
	// 判断条件并进入对应分支逻辑。
	if req.DrawCode != nil {
		// 更新当前变量或字段值。
		next.DrawCode = strings.TrimSpace(*req.DrawCode)
	}
	// 判断条件并进入对应分支逻辑。
	if req.DrawAt != nil {
		// 更新当前变量或字段值。
		next.DrawAt = parseDateTimeOrDefault(*req.DrawAt, next.DrawAt)
	}
	// 判断条件并进入对应分支逻辑。
	if req.PlaybackURL != nil {
		// 定义并初始化当前变量。
		playbackURL, playbackErr := normalizeSafeURL(*req.PlaybackURL, true)
		// 判断条件并进入对应分支逻辑。
		if playbackErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid playback_url")
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		next.PlaybackURL = playbackURL
	}
	// 判断条件并进入对应分支逻辑。
	if req.LikesCount != nil {
		// 更新当前变量或字段值。
		next.LikesCount = *req.LikesCount
	}
	// 判断条件并进入对应分支逻辑。
	if req.CommentCount != nil {
		// 更新当前变量或字段值。
		next.CommentCount = *req.CommentCount
	}
	// 判断条件并进入对应分支逻辑。
	if req.FavoriteCount != nil {
		// 更新当前变量或字段值。
		next.FavoriteCount = *req.FavoriteCount
	}
	// 判断条件并进入对应分支逻辑。
	if req.ReadCount != nil {
		// 更新当前变量或字段值。
		next.ReadCount = *req.ReadCount
	}
	// 判断条件并进入对应分支逻辑。
	if req.PollEnabled != nil {
		// 更新当前变量或字段值。
		next.PollEnabled = *req.PollEnabled
	}
	// 判断条件并进入对应分支逻辑。
	if req.PollDefaultExpand != nil {
		// 更新当前变量或字段值。
		next.PollDefaultExpand = *req.PollDefaultExpand
	}
	// 判断条件并进入对应分支逻辑。
	if req.RecommendInfoIDs != nil {
		// 更新当前变量或字段值。
		next.RecommendInfoIDs = strings.TrimSpace(*req.RecommendInfoIDs)
	}
	// 判断条件并进入对应分支逻辑。
	if req.IsCurrent != nil {
		// 更新当前变量或字段值。
		next.IsCurrent = *req.IsCurrent
	}
	// 判断条件并进入对应分支逻辑。
	if req.Status != nil {
		// 更新当前变量或字段值。
		next.Status = *req.Status
	}
	// 判断条件并进入对应分支逻辑。
	if req.Sort != nil {
		// 更新当前变量或字段值。
		next.Sort = *req.Sort
	}

	// 分类字段任一变化时，都按分类表重算 category_id + category_tag。
	if req.CategoryID != nil || req.CategoryTag != nil {
		// 定义并初始化当前变量。
		resolvedID, resolvedTag, resolveErr := bc.svc.ResolveLotteryCategory(
			// 调用c.Request.Context完成当前处理。
			c.Request.Context(),
			// 调用valueOrCurrentUint完成当前处理。
			valueOrCurrentUint(req.CategoryID, next.CategoryID),
			// 调用valueOrCurrentString完成当前处理。
			valueOrCurrentString(req.CategoryTag, next.CategoryTag),
		)
		// 判断条件并进入对应分支逻辑。
		if resolveErr != nil {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizInvalidRequest, resolveErr.Error())
			// 返回当前处理结果。
			return
		}
		// 更新当前变量或字段值。
		next.CategoryID = resolvedID
		// 更新当前变量或字段值。
		next.CategoryTag = resolvedTag
	}

	// 号码字段任一变化时，重新校验并生成兼容字段 draw_result。
	if req.NormalDrawResult != nil || req.SpecialDrawResult != nil || req.DrawResult != nil {
		// 定义并初始化当前变量。
		normalInput := safeString(valueOrCurrentString(req.NormalDrawResult, next.NormalDrawResult))
		// 定义并初始化当前变量。
		specialInput := safeString(valueOrCurrentString(req.SpecialDrawResult, next.SpecialDrawResult))
		// 定义并初始化当前变量。
		mergedInput := safeString(valueOrCurrentString(req.DrawResult, next.DrawResult))
		// 三个号码字段都为空时，允许清空（图库内容与开奖区已解耦）。
		if strings.TrimSpace(normalInput) == "" && strings.TrimSpace(specialInput) == "" && strings.TrimSpace(mergedInput) == "" {
			// 更新当前变量或字段值。
			next.NormalDrawResult = ""
			// 更新当前变量或字段值。
			next.SpecialDrawResult = ""
			// 更新当前变量或字段值。
			next.DrawResult = ""
			// 进入新的代码块进行处理。
		} else {
			// 定义并初始化当前变量。
			normalRaw, specialRaw, mergedRaw, drawErr := normalizeAndMergeDrawNumbers(normalInput, specialInput, mergedInput)
			// 判断条件并进入对应分支逻辑。
			if drawErr != nil {
				// 调用utils.JSONError完成当前处理。
				commonresp.GinError(c, constants.AdminBizInvalidRequest, drawErr.Error())
				// 返回当前处理结果。
				return
			}
			// 更新当前变量或字段值。
			next.NormalDrawResult = normalRaw
			// 更新当前变量或字段值。
			next.SpecialDrawResult = specialRaw
			// 更新当前变量或字段值。
			next.DrawResult = mergedRaw
		}
	}
	// 更新时仅当请求显式传了 option_names 才改动物竞猜选项。
	var optionNames []string
	// 定义并初始化当前变量。
	updateOptions := false
	// 判断条件并进入对应分支逻辑。
	if req.OptionNames != nil {
		// 更新当前变量或字段值。
		updateOptions = true
		// 更新当前变量或字段值。
		optionNames = normalizeOptionNames(*req.OptionNames)
		// 判断条件并进入对应分支逻辑。
		if len(optionNames) == 0 {
			// 更新当前变量或字段值。
			optionNames = defaultAnimalOptionNames()
		}
	}

	// 再次校验基础必填。
	if strings.TrimSpace(next.Issue) == "" || strings.TrimSpace(next.Title) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "issue/title required")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	updates := map[string]interface{}{
		// 处理当前语句逻辑。
		"special_lottery_id": next.SpecialLotteryID,
		// 处理当前语句逻辑。
		"category_id": next.CategoryID,
		// 处理当前语句逻辑。
		"category_tag": next.CategoryTag,
		// 处理当前语句逻辑。
		"issue": next.Issue,
		// 处理当前语句逻辑。
		"year": next.Year,
		// 处理当前语句逻辑。
		"title": next.Title,
		// 处理当前语句逻辑。
		"cover_image_url": next.CoverImageURL,
		// 处理当前语句逻辑。
		"detail_image_url": next.DetailImageURL,
		// 处理当前语句逻辑。
		"draw_code": next.DrawCode,
		// 处理当前语句逻辑。
		"normal_draw_result": next.NormalDrawResult,
		// 处理当前语句逻辑。
		"special_draw_result": next.SpecialDrawResult,
		// 处理当前语句逻辑。
		"draw_result": next.DrawResult,
		// 处理当前语句逻辑。
		"draw_at": next.DrawAt,
		// 处理当前语句逻辑。
		"playback_url": next.PlaybackURL,
		// 处理当前语句逻辑。
		"likes_count": next.LikesCount,
		// 处理当前语句逻辑。
		"comment_count": next.CommentCount,
		// 处理当前语句逻辑。
		"favorite_count": next.FavoriteCount,
		// 处理当前语句逻辑。
		"read_count": next.ReadCount,
		// 处理当前语句逻辑。
		"poll_enabled": next.PollEnabled,
		// 处理当前语句逻辑。
		"poll_default_expand": next.PollDefaultExpand,
		// 处理当前语句逻辑。
		"recommend_info_ids": next.RecommendInfoIDs,
		// 处理当前语句逻辑。
		"is_current": next.IsCurrent,
		// 处理当前语句逻辑。
		"status": next.Status,
		// 处理当前语句逻辑。
		"sort": next.Sort,
	}

	// 更新时同样维护“每彩种唯一当前期”约束。
	if err := bc.svc.UpdateLotteryInfo(c.Request.Context(), id, updates, updateOptions, optionNames, next.SpecialLotteryID, next.IsCurrent); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteLotteryInfo 删除图库内容。
func (bc *BizConfigController) DeleteLotteryInfo(c *gin.Context) {
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
	if err := bc.svc.DeleteLotteryInfo(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// normalizeAndMergeDrawNumbers 归一化 6+1 开奖号码，并输出兼容字段 draw_result。
func normalizeAndMergeDrawNumbers(normalRaw, specialRaw, drawRaw string) (string, string, string, error) {
	// 更新当前变量或字段值。
	normalRaw = strings.TrimSpace(normalRaw)
	// 更新当前变量或字段值。
	specialRaw = strings.TrimSpace(specialRaw)
	// 更新当前变量或字段值。
	drawRaw = strings.TrimSpace(drawRaw)

	// 兼容旧表单：仅传 draw_result 时，按“前6个普通号 + 最后1个特别号”拆解。
	if normalRaw == "" && specialRaw == "" && drawRaw != "" {
		// 定义并初始化当前变量。
		all, err := parseDrawNumbers(drawRaw)
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 返回当前处理结果。
			return "", "", "", err
		}
		// 判断条件并进入对应分支逻辑。
		if len(all) != 7 {
			// 返回当前处理结果。
			return "", "", "", fmt.Errorf("draw_result must contain 7 numbers")
		}
		// 更新当前变量或字段值。
		normalRaw = joinIntCSV(all[:6])
		// 更新当前变量或字段值。
		specialRaw = strconv.Itoa(all[6])
	}

	// 定义并初始化当前变量。
	normalNums, err := parseDrawNumbers(normalRaw)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", "", err
	}
	// 判断条件并进入对应分支逻辑。
	if len(normalNums) != 6 {
		// 返回当前处理结果。
		return "", "", "", fmt.Errorf("normal_draw_result must contain 6 numbers")
	}
	// 定义并初始化当前变量。
	specialNums, err := parseDrawNumbers(specialRaw)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return "", "", "", err
	}
	// 判断条件并进入对应分支逻辑。
	if len(specialNums) != 1 {
		// 返回当前处理结果。
		return "", "", "", fmt.Errorf("special_draw_result must contain 1 number")
	}
	// 特别号不能与普通号重复，避免无效开奖数据入库。
	if containsInt(normalNums, specialNums[0]) {
		// 返回当前处理结果。
		return "", "", "", fmt.Errorf("special_draw_result cannot duplicate normal numbers")
	}
	// 定义并初始化当前变量。
	normalizedNormal := joinIntCSV(normalNums)
	// 定义并初始化当前变量。
	normalizedSpecial := strconv.Itoa(specialNums[0])
	// 定义并初始化当前变量。
	merged := normalizedNormal + "," + normalizedSpecial
	// 返回当前处理结果。
	return normalizedNormal, normalizedSpecial, merged, nil
}

// parseDrawNumbers 解析号码串，支持“逗号/空格/斜杠/竖线”分隔。
func parseDrawNumbers(raw string) ([]int, error) {
	// 定义并初始化当前变量。
	tokens := strings.FieldsFunc(strings.TrimSpace(raw), func(r rune) bool {
		// 返回当前处理结果。
		return r == ',' || r == '|' || r == '/' || unicode.IsSpace(r)
	})
	// 判断条件并进入对应分支逻辑。
	if len(tokens) == 0 {
		// 返回当前处理结果。
		return []int{}, nil
	}
	// 定义并初始化当前变量。
	out := make([]int, 0, len(tokens))
	// 定义并初始化当前变量。
	seen := map[int]struct{}{}
	// 循环处理当前数据集合。
	for _, token := range tokens {
		// 定义并初始化当前变量。
		v, err := strconv.Atoi(strings.TrimSpace(token))
		// 判断条件并进入对应分支逻辑。
		if err != nil {
			// 返回当前处理结果。
			return nil, fmt.Errorf("invalid draw number: %s", token)
		}
		// 判断条件并进入对应分支逻辑。
		if v < 1 || v > 49 {
			// 返回当前处理结果。
			return nil, fmt.Errorf("draw number out of range: %d", v)
		}
		// 判断条件并进入对应分支逻辑。
		if _, ok := seen[v]; ok {
			// 返回当前处理结果。
			return nil, fmt.Errorf("duplicate draw number: %d", v)
		}
		// 更新当前变量或字段值。
		seen[v] = struct{}{}
		// 更新当前变量或字段值。
		out = append(out, v)
	}
	// 返回当前处理结果。
	return out, nil
}

// joinIntCSV 将整型数组按逗号拼接为字符串。
func joinIntCSV(nums []int) string {
	// 判断条件并进入对应分支逻辑。
	if len(nums) == 0 {
		// 返回当前处理结果。
		return ""
	}
	// 定义并初始化当前变量。
	parts := make([]string, 0, len(nums))
	// 循环处理当前数据集合。
	for _, n := range nums {
		// 更新当前变量或字段值。
		parts = append(parts, strconv.Itoa(n))
	}
	// 返回当前处理结果。
	return strings.Join(parts, ",")
}

// containsInt 判断切片是否包含指定值。
func containsInt(nums []int, target int) bool {
	// 循环处理当前数据集合。
	for _, n := range nums {
		// 判断条件并进入对应分支逻辑。
		if n == target {
			// 返回当前处理结果。
			return true
		}
	}
	// 返回当前处理结果。
	return false
}

// parseDateTimeOrDefault 解析时间字符串，失败时回退默认值。
func parseDateTimeOrDefault(raw string, fallback time.Time) time.Time {
	// 更新当前变量或字段值。
	raw = strings.TrimSpace(raw)
	// 判断条件并进入对应分支逻辑。
	if raw == "" {
		// 返回当前处理结果。
		return fallback
	}
	// 定义并初始化当前变量。
	layouts := []string{
		// 处理当前语句逻辑。
		time.RFC3339,
		// 处理当前语句逻辑。
		"2006-01-02 15:04:05",
		// 处理当前语句逻辑。
		"2006-01-02 15:04",
		// 处理当前语句逻辑。
		"2006-01-02T15:04:05",
	}
	// 循环处理当前数据集合。
	for _, layout := range layouts {
		// 判断条件并进入对应分支逻辑。
		if t, err := time.Parse(layout, raw); err == nil {
			// 返回当前处理结果。
			return t
		}
	}
	// 返回当前处理结果。
	return fallback
}

// parseDailyDrawTime 解析“每天开奖时刻”。
// 1) 支持 HH:mm / HH:mm:ss；
// 2) 兼容历史 datetime/RFC3339 输入；
// 3) 最终仅保留“时分秒”，日期统一使用 fallback 当天，避免每天手动改日期。
// 4) 输入无效时返回错误，避免错误值静默回退到当前时间。
func parseDailyDrawTime(raw string, fallback time.Time) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty next_draw_at")
	}
	loc := fallback.Location()
	// 先直接提取时分秒，避免 datetime 时区转换导致“东京填 21:30 被转成 20:30”。
	if h, m, s, ok := extractClockFromText(raw); ok {
		return time.Date(
			fallback.Year(),
			fallback.Month(),
			fallback.Day(),
			h,
			m,
			s,
			0,
			loc,
		), nil
	}
	for _, layout := range []string{"15:04:05", "15:04"} {
		if t, err := time.ParseInLocation(layout, raw, loc); err == nil {
			return time.Date(
				fallback.Year(),
				fallback.Month(),
				fallback.Day(),
				t.Hour(),
				t.Minute(),
				t.Second(),
				0,
				loc,
			), nil
		}
	}
	parsed := time.Time{}
	parsedOK := false
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		parsed = t
		parsedOK = true
	} else {
		for _, layout := range []string{
			time.RFC3339Nano,
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04:05.000Z07:00",
			"2006-01-02T15:04:05.000Z",
		} {
			if dt, dtErr := time.ParseInLocation(layout, raw, loc); dtErr == nil {
				parsed = dt
				parsedOK = true
				break
			}
		}
	}
	if !parsedOK {
		return time.Time{}, fmt.Errorf("invalid next_draw_at format")
	}
	return time.Date(
		fallback.Year(),
		fallback.Month(),
		fallback.Day(),
		parsed.Hour(),
		parsed.Minute(),
		parsed.Second(),
		0,
		loc,
	), nil
}

var dailyClockMatcher = regexp.MustCompile(`(\d{1,2}):(\d{2})(?::(\d{2}))?`)

// extractClockFromText 从任意时间文本中提取时分秒（忽略日期与时区部分）。
func extractClockFromText(raw string) (int, int, int, bool) {
	matched := dailyClockMatcher.FindStringSubmatch(strings.TrimSpace(raw))
	if len(matched) < 3 {
		return 0, 0, 0, false
	}
	h, hErr := strconv.Atoi(matched[1])
	m, mErr := strconv.Atoi(matched[2])
	s := 0
	var sErr error
	if len(matched) >= 4 && strings.TrimSpace(matched[3]) != "" {
		s, sErr = strconv.Atoi(matched[3])
	}
	if hErr != nil || mErr != nil || sErr != nil {
		return 0, 0, 0, false
	}
	if h < 0 || h > 23 || m < 0 || m > 59 || s < 0 || s > 59 {
		return 0, 0, 0, false
	}
	return h, m, s, true
}

// lotteryNowInEast8 返回当前东八区时间，作为开奖配置统一时区基准。
func lotteryNowInEast8() time.Time {
	return time.Now().In(lotteryLocationEast8())
}

// lotteryLocationEast8 返回开奖业务使用的固定时区（东八区）。
func lotteryLocationEast8() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		return loc
	}
	return time.FixedZone("UTC+8", 8*3600)
}

// safeString 读取可空字符串指针并做空格裁剪。
func safeString(v *string) string {
	// 判断条件并进入对应分支逻辑。
	if v == nil {
		// 返回当前处理结果。
		return ""
	}
	// 返回当前处理结果。
	return strings.TrimSpace(*v)
}

// safeInt8 读取可空 int8 指针并提供默认值。
func safeInt8(v *int8, def int8) int8 {
	// 判断条件并进入对应分支逻辑。
	if v == nil {
		// 返回当前处理结果。
		return def
	}
	// 返回当前处理结果。
	return *v
}

// safeInt 读取可空 int 指针并提供默认值。
func safeInt(v *int, def int) int {
	// 判断条件并进入对应分支逻辑。
	if v == nil {
		// 返回当前处理结果。
		return def
	}
	// 返回当前处理结果。
	return *v
}

// safeInt64 读取可空 int64 指针并提供默认值。
func safeInt64(v *int64, def int64) int64 {
	// 判断条件并进入对应分支逻辑。
	if v == nil {
		// 返回当前处理结果。
		return def
	}
	// 返回当前处理结果。
	return *v
}

// valueOrCurrentUint 当请求值为空时回落到当前值。
func valueOrCurrentUint(v *uint, current uint) *uint {
	// 判断条件并进入对应分支逻辑。
	if v != nil {
		// 返回当前处理结果。
		return v
	}
	// 返回当前处理结果。
	return &current
}

// valueOrCurrentString 当请求值为空时回落到当前值。
func valueOrCurrentString(v *string, current string) *string {
	// 判断条件并进入对应分支逻辑。
	if v != nil {
		// 返回当前处理结果。
		return v
	}
	// 返回当前处理结果。
	return &current
}

// normalizeOptionNames 去重清洗动物竞猜选项，保持输入顺序。
func normalizeOptionNames(input []string) []string {
	// 定义并初始化当前变量。
	out := make([]string, 0, len(input))
	// 定义并初始化当前变量。
	seen := make(map[string]struct{}, len(input))
	// 循环处理当前数据集合。
	for _, raw := range input {
		// 定义并初始化当前变量。
		name := strings.TrimSpace(raw)
		// 判断条件并进入对应分支逻辑。
		if name == "" {
			// 处理当前语句逻辑。
			continue
		}
		// 判断条件并进入对应分支逻辑。
		if _, ok := seen[name]; ok {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		seen[name] = struct{}{}
		// 更新当前变量或字段值。
		out = append(out, name)
	}
	// 返回当前处理结果。
	return out
}

// defaultAnimalOptionNames 返回默认 12 生肖竞猜选项。
func defaultAnimalOptionNames() []string {
	// 返回当前处理结果。
	return []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
}
