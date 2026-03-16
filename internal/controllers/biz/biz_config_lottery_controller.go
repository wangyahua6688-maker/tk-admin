package biz

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
)

// -------------------- Special Lottery --------------------

// ListSpecialLotteries 查询彩种配置列表。
func (bc *BizConfigController) ListSpecialLotteries(c *gin.Context) {
	// 固定按 sort + id 排序，确保前端按钮顺序稳定。
	items, err := bc.svc.ListSpecialLotteries(c.Request.Context(), 200)
	if err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

// CreateSpecialLottery 新增彩种配置。
func (bc *BizConfigController) CreateSpecialLottery(c *gin.Context) {
	var req struct {
		Name          string `json:"name"`
		Code          string `json:"code"`
		CurrentIssue  string `json:"current_issue"`
		NextDrawAt    string `json:"next_draw_at"`
		LiveEnabled   *int8  `json:"live_enabled"`
		LiveStatus    string `json:"live_status"`
		LiveStreamURL string `json:"live_stream_url"`
		Status        *int8  `json:"status"`
		Sort          *int   `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	// 名称和编码为必填。
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Code) == "" {
		utils.JSONError(c, http.StatusBadRequest, "name/code required")
		return
	}
	// 下期开奖时间无法解析时，兜底当前时间，避免写入空值。
	next := parseDateTimeOrDefault(req.NextDrawAt, time.Now())
	item := models.WSpecialLottery{
		Name:          strings.TrimSpace(req.Name),
		Code:          strings.TrimSpace(req.Code),
		CurrentIssue:  strings.TrimSpace(req.CurrentIssue),
		NextDrawAt:    next,
		LiveEnabled:   0,
		LiveStatus:    strings.TrimSpace(req.LiveStatus),
		LiveStreamURL: strings.TrimSpace(req.LiveStreamURL),
		Status:        1,
		Sort:          0,
	}
	// 直播状态为空时，默认 pending。
	if item.LiveStatus == "" {
		item.LiveStatus = "pending"
	}
	if req.LiveEnabled != nil {
		item.LiveEnabled = *req.LiveEnabled
	}
	if req.Status != nil {
		item.Status = *req.Status
	}
	if req.Sort != nil {
		item.Sort = *req.Sort
	}
	if err := bc.svc.CreateSpecialLottery(c.Request.Context(), &item); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

// UpdateSpecialLottery 更新彩种配置。
func (bc *BizConfigController) UpdateSpecialLottery(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Name          *string `json:"name"`
		Code          *string `json:"code"`
		CurrentIssue  *string `json:"current_issue"`
		NextDrawAt    *string `json:"next_draw_at"`
		LiveEnabled   *int8   `json:"live_enabled"`
		LiveStatus    *string `json:"live_status"`
		LiveStreamURL *string `json:"live_stream_url"`
		Status        *int8   `json:"status"`
		Sort          *int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Code != nil {
		updates["code"] = strings.TrimSpace(*req.Code)
	}
	if req.CurrentIssue != nil {
		updates["current_issue"] = strings.TrimSpace(*req.CurrentIssue)
	}
	if req.NextDrawAt != nil {
		updates["next_draw_at"] = parseDateTimeOrDefault(*req.NextDrawAt, time.Now())
	}
	if req.LiveEnabled != nil {
		updates["live_enabled"] = *req.LiveEnabled
	}
	if req.LiveStatus != nil {
		updates["live_status"] = strings.TrimSpace(*req.LiveStatus)
	}
	if req.LiveStreamURL != nil {
		updates["live_stream_url"] = strings.TrimSpace(*req.LiveStreamURL)
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if len(updates) == 0 {
		utils.JSONError(c, http.StatusBadRequest, "empty updates")
		return
	}
	if err := bc.svc.UpdateSpecialLottery(c.Request.Context(), id, updates); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteSpecialLottery 删除彩种配置。
func (bc *BizConfigController) DeleteSpecialLottery(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.svc.DeleteSpecialLottery(c.Request.Context(), id); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// -------------------- Lottery Info --------------------

// lotteryInfoUpsertRequest 图库内容新增/编辑请求结构。
type lotteryInfoUpsertRequest struct {
	SpecialLotteryID  *uint     `json:"special_lottery_id"`
	CategoryID        *uint     `json:"category_id"`
	CategoryTag       *string   `json:"category_tag"`
	Issue             *string   `json:"issue"`
	Year              *int      `json:"year"`
	Title             *string   `json:"title"`
	CoverImageURL     *string   `json:"cover_image_url"`
	DetailImageURL    *string   `json:"detail_image_url"`
	DrawCode          *string   `json:"draw_code"`
	NormalDrawResult  *string   `json:"normal_draw_result"`
	SpecialDrawResult *string   `json:"special_draw_result"`
	DrawResult        *string   `json:"draw_result"`
	DrawAt            *string   `json:"draw_at"`
	PlaybackURL       *string   `json:"playback_url"`
	LikesCount        *int64    `json:"likes_count"`
	CommentCount      *int64    `json:"comment_count"`
	FavoriteCount     *int64    `json:"favorite_count"`
	ReadCount         *int64    `json:"read_count"`
	PollEnabled       *int8     `json:"poll_enabled"`
	PollDefaultExpand *int8     `json:"poll_default_expand"`
	RecommendInfoIDs  *string   `json:"recommend_info_ids"`
	OptionNames       *[]string `json:"option_names"`
	IsCurrent         *int8     `json:"is_current"`
	Status            *int8     `json:"status"`
	Sort              *int      `json:"sort"`
}

// ListLotteryInfos 查询图库内容列表。
func (bc *BizConfigController) ListLotteryInfos(c *gin.Context) {
	// 图库内容管理优先按更新时间倒序，避免运营修改后找不到记录。
	items, optionNameMap, err := bc.svc.ListLotteryInfosWithOptions(c.Request.Context(), 300)
	if err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{
		"items":                   items,
		"option_names_by_info_id": optionNameMap,
	})
}

// CreateLotteryInfo 新增图库内容。
func (bc *BizConfigController) CreateLotteryInfo(c *gin.Context) {
	var req lotteryInfoUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	// 校验基础必填项（图库内容不再强制绑定彩种）。
	if req.Issue == nil || strings.TrimSpace(*req.Issue) == "" || req.Title == nil || strings.TrimSpace(*req.Title) == "" {
		utils.JSONError(c, http.StatusBadRequest, "issue/title required")
		return
	}
	// 图库内容可独立存在：special_lottery_id 允许为空，落库为 0。
	// 说明：彩种维度只在“开奖区开奖记录（tk_draw_record）”中强约束。
	specialLotteryID := uint(0)
	if req.SpecialLotteryID != nil {
		specialLotteryID = *req.SpecialLotteryID
	}
	// 分类必须可落到 tk_lottery_category，禁止手工自由文本分类。
	categoryID, categoryTag, err := bc.svc.ResolveLotteryCategory(c.Request.Context(), req.CategoryID, req.CategoryTag)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	// 图库内容不强制录入 6+1 开奖码，仅在传值时做校验。
	normalRaw, specialRaw, mergedRaw := "", "", ""
	if strings.TrimSpace(safeString(req.NormalDrawResult)) != "" ||
		strings.TrimSpace(safeString(req.SpecialDrawResult)) != "" ||
		strings.TrimSpace(safeString(req.DrawResult)) != "" {
		var normalizeErr error
		normalRaw, specialRaw, mergedRaw, normalizeErr = normalizeAndMergeDrawNumbers(
			safeString(req.NormalDrawResult),
			safeString(req.SpecialDrawResult),
			safeString(req.DrawResult),
		)
		if normalizeErr != nil {
			utils.JSONError(c, http.StatusBadRequest, normalizeErr.Error())
			return
		}
	}
	// 预处理动物竞猜选项：为空时回退到 12 生肖默认选项。
	optionNames := normalizeOptionNames(nil)
	if req.OptionNames != nil {
		optionNames = normalizeOptionNames(*req.OptionNames)
	}
	if len(optionNames) == 0 {
		optionNames = defaultAnimalOptionNames()
	}
	now := time.Now()
	item := models.WLotteryInfo{
		SpecialLotteryID:  specialLotteryID,
		CategoryID:        categoryID,
		CategoryTag:       categoryTag,
		Issue:             strings.TrimSpace(*req.Issue),
		Year:              now.Year(),
		Title:             strings.TrimSpace(*req.Title),
		CoverImageURL:     safeString(req.CoverImageURL),
		DetailImageURL:    safeString(req.DetailImageURL),
		DrawCode:          safeString(req.DrawCode),
		NormalDrawResult:  normalRaw,
		SpecialDrawResult: specialRaw,
		DrawResult:        mergedRaw,
		DrawAt:            parseDateTimeOrDefault(safeString(req.DrawAt), now),
		PlaybackURL:       safeString(req.PlaybackURL),
		LikesCount:        safeInt64(req.LikesCount, 0),
		CommentCount:      safeInt64(req.CommentCount, 0),
		FavoriteCount:     safeInt64(req.FavoriteCount, 0),
		ReadCount:         safeInt64(req.ReadCount, 0),
		PollEnabled:       safeInt8(req.PollEnabled, 1),
		PollDefaultExpand: safeInt8(req.PollDefaultExpand, 0),
		RecommendInfoIDs:  safeString(req.RecommendInfoIDs),
		IsCurrent:         safeInt8(req.IsCurrent, 0),
		Status:            safeInt8(req.Status, 1),
		Sort:              safeInt(req.Sort, 0),
	}
	if req.Year != nil && *req.Year > 0 {
		item.Year = *req.Year
	}

	// 当前期约束：同一彩种只能有一条 is_current=1 记录。
	if err := bc.svc.CreateLotteryInfo(c.Request.Context(), &item, optionNames); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, item)
}

// UpdateLotteryInfo 编辑图库内容。
func (bc *BizConfigController) UpdateLotteryInfo(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req lotteryInfoUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	current, err := bc.svc.GetLotteryInfoByID(c.Request.Context(), id)
	if err != nil {
		utils.JSONError(c, http.StatusNotFound, "lottery info not found")
		return
	}

	// 在“当前值”基础上叠加本次更新，保证部分更新不会丢字段。
	next := *current
	if req.SpecialLotteryID != nil {
		// 图库内容允许取消彩种绑定（设置为 0）。
		// 前端“图库管理”页已移除彩种选择，这里保留兼容更新能力。
		next.SpecialLotteryID = *req.SpecialLotteryID
	}
	if req.Issue != nil {
		next.Issue = strings.TrimSpace(*req.Issue)
	}
	if req.Title != nil {
		next.Title = strings.TrimSpace(*req.Title)
	}
	if req.Year != nil && *req.Year > 0 {
		next.Year = *req.Year
	}
	if req.CoverImageURL != nil {
		next.CoverImageURL = strings.TrimSpace(*req.CoverImageURL)
	}
	if req.DetailImageURL != nil {
		next.DetailImageURL = strings.TrimSpace(*req.DetailImageURL)
	}
	if req.DrawCode != nil {
		next.DrawCode = strings.TrimSpace(*req.DrawCode)
	}
	if req.DrawAt != nil {
		next.DrawAt = parseDateTimeOrDefault(*req.DrawAt, next.DrawAt)
	}
	if req.PlaybackURL != nil {
		next.PlaybackURL = strings.TrimSpace(*req.PlaybackURL)
	}
	if req.LikesCount != nil {
		next.LikesCount = *req.LikesCount
	}
	if req.CommentCount != nil {
		next.CommentCount = *req.CommentCount
	}
	if req.FavoriteCount != nil {
		next.FavoriteCount = *req.FavoriteCount
	}
	if req.ReadCount != nil {
		next.ReadCount = *req.ReadCount
	}
	if req.PollEnabled != nil {
		next.PollEnabled = *req.PollEnabled
	}
	if req.PollDefaultExpand != nil {
		next.PollDefaultExpand = *req.PollDefaultExpand
	}
	if req.RecommendInfoIDs != nil {
		next.RecommendInfoIDs = strings.TrimSpace(*req.RecommendInfoIDs)
	}
	if req.IsCurrent != nil {
		next.IsCurrent = *req.IsCurrent
	}
	if req.Status != nil {
		next.Status = *req.Status
	}
	if req.Sort != nil {
		next.Sort = *req.Sort
	}

	// 分类字段任一变化时，都按分类表重算 category_id + category_tag。
	if req.CategoryID != nil || req.CategoryTag != nil {
		resolvedID, resolvedTag, resolveErr := bc.svc.ResolveLotteryCategory(
			c.Request.Context(),
			valueOrCurrentUint(req.CategoryID, next.CategoryID),
			valueOrCurrentString(req.CategoryTag, next.CategoryTag),
		)
		if resolveErr != nil {
			utils.JSONError(c, http.StatusBadRequest, resolveErr.Error())
			return
		}
		next.CategoryID = resolvedID
		next.CategoryTag = resolvedTag
	}

	// 号码字段任一变化时，重新校验并生成兼容字段 draw_result。
	if req.NormalDrawResult != nil || req.SpecialDrawResult != nil || req.DrawResult != nil {
		normalInput := safeString(valueOrCurrentString(req.NormalDrawResult, next.NormalDrawResult))
		specialInput := safeString(valueOrCurrentString(req.SpecialDrawResult, next.SpecialDrawResult))
		mergedInput := safeString(valueOrCurrentString(req.DrawResult, next.DrawResult))
		// 三个号码字段都为空时，允许清空（图库内容与开奖区已解耦）。
		if strings.TrimSpace(normalInput) == "" && strings.TrimSpace(specialInput) == "" && strings.TrimSpace(mergedInput) == "" {
			next.NormalDrawResult = ""
			next.SpecialDrawResult = ""
			next.DrawResult = ""
		} else {
			normalRaw, specialRaw, mergedRaw, drawErr := normalizeAndMergeDrawNumbers(normalInput, specialInput, mergedInput)
			if drawErr != nil {
				utils.JSONError(c, http.StatusBadRequest, drawErr.Error())
				return
			}
			next.NormalDrawResult = normalRaw
			next.SpecialDrawResult = specialRaw
			next.DrawResult = mergedRaw
		}
	}
	// 更新时仅当请求显式传了 option_names 才改动物竞猜选项。
	var optionNames []string
	updateOptions := false
	if req.OptionNames != nil {
		updateOptions = true
		optionNames = normalizeOptionNames(*req.OptionNames)
		if len(optionNames) == 0 {
			optionNames = defaultAnimalOptionNames()
		}
	}

	// 再次校验基础必填。
	if strings.TrimSpace(next.Issue) == "" || strings.TrimSpace(next.Title) == "" {
		utils.JSONError(c, http.StatusBadRequest, "issue/title required")
		return
	}

	updates := map[string]interface{}{
		"special_lottery_id":  next.SpecialLotteryID,
		"category_id":         next.CategoryID,
		"category_tag":        next.CategoryTag,
		"issue":               next.Issue,
		"year":                next.Year,
		"title":               next.Title,
		"cover_image_url":     next.CoverImageURL,
		"detail_image_url":    next.DetailImageURL,
		"draw_code":           next.DrawCode,
		"normal_draw_result":  next.NormalDrawResult,
		"special_draw_result": next.SpecialDrawResult,
		"draw_result":         next.DrawResult,
		"draw_at":             next.DrawAt,
		"playback_url":        next.PlaybackURL,
		"likes_count":         next.LikesCount,
		"comment_count":       next.CommentCount,
		"favorite_count":      next.FavoriteCount,
		"read_count":          next.ReadCount,
		"poll_enabled":        next.PollEnabled,
		"poll_default_expand": next.PollDefaultExpand,
		"recommend_info_ids":  next.RecommendInfoIDs,
		"is_current":          next.IsCurrent,
		"status":              next.Status,
		"sort":                next.Sort,
	}

	// 更新时同样维护“每彩种唯一当前期”约束。
	if err := bc.svc.UpdateLotteryInfo(c.Request.Context(), id, updates, updateOptions, optionNames, next.SpecialLotteryID, next.IsCurrent); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteLotteryInfo 删除图库内容。
func (bc *BizConfigController) DeleteLotteryInfo(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.svc.DeleteLotteryInfo(c.Request.Context(), id); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// normalizeAndMergeDrawNumbers 归一化 6+1 开奖号码，并输出兼容字段 draw_result。
func normalizeAndMergeDrawNumbers(normalRaw, specialRaw, drawRaw string) (string, string, string, error) {
	normalRaw = strings.TrimSpace(normalRaw)
	specialRaw = strings.TrimSpace(specialRaw)
	drawRaw = strings.TrimSpace(drawRaw)

	// 兼容旧表单：仅传 draw_result 时，按“前6个普通号 + 最后1个特别号”拆解。
	if normalRaw == "" && specialRaw == "" && drawRaw != "" {
		all, err := parseDrawNumbers(drawRaw)
		if err != nil {
			return "", "", "", err
		}
		if len(all) != 7 {
			return "", "", "", fmt.Errorf("draw_result must contain 7 numbers")
		}
		normalRaw = joinIntCSV(all[:6])
		specialRaw = strconv.Itoa(all[6])
	}

	normalNums, err := parseDrawNumbers(normalRaw)
	if err != nil {
		return "", "", "", err
	}
	if len(normalNums) != 6 {
		return "", "", "", fmt.Errorf("normal_draw_result must contain 6 numbers")
	}
	specialNums, err := parseDrawNumbers(specialRaw)
	if err != nil {
		return "", "", "", err
	}
	if len(specialNums) != 1 {
		return "", "", "", fmt.Errorf("special_draw_result must contain 1 number")
	}
	// 特别号不能与普通号重复，避免无效开奖数据入库。
	if containsInt(normalNums, specialNums[0]) {
		return "", "", "", fmt.Errorf("special_draw_result cannot duplicate normal numbers")
	}
	normalizedNormal := joinIntCSV(normalNums)
	normalizedSpecial := strconv.Itoa(specialNums[0])
	merged := normalizedNormal + "," + normalizedSpecial
	return normalizedNormal, normalizedSpecial, merged, nil
}

// parseDrawNumbers 解析号码串，支持“逗号/空格/斜杠/竖线”分隔。
func parseDrawNumbers(raw string) ([]int, error) {
	tokens := strings.FieldsFunc(strings.TrimSpace(raw), func(r rune) bool {
		return r == ',' || r == '|' || r == '/' || unicode.IsSpace(r)
	})
	if len(tokens) == 0 {
		return []int{}, nil
	}
	out := make([]int, 0, len(tokens))
	seen := map[int]struct{}{}
	for _, token := range tokens {
		v, err := strconv.Atoi(strings.TrimSpace(token))
		if err != nil {
			return nil, fmt.Errorf("invalid draw number: %s", token)
		}
		if v < 1 || v > 49 {
			return nil, fmt.Errorf("draw number out of range: %d", v)
		}
		if _, ok := seen[v]; ok {
			return nil, fmt.Errorf("duplicate draw number: %d", v)
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out, nil
}

// joinIntCSV 将整型数组按逗号拼接为字符串。
func joinIntCSV(nums []int) string {
	if len(nums) == 0 {
		return ""
	}
	parts := make([]string, 0, len(nums))
	for _, n := range nums {
		parts = append(parts, strconv.Itoa(n))
	}
	return strings.Join(parts, ",")
}

// containsInt 判断切片是否包含指定值。
func containsInt(nums []int, target int) bool {
	for _, n := range nums {
		if n == target {
			return true
		}
	}
	return false
}

// parseDateTimeOrDefault 解析时间字符串，失败时回退默认值。
func parseDateTimeOrDefault(raw string, fallback time.Time) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return fallback
}

// safeString 读取可空字符串指针并做空格裁剪。
func safeString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

// safeInt8 读取可空 int8 指针并提供默认值。
func safeInt8(v *int8, def int8) int8 {
	if v == nil {
		return def
	}
	return *v
}

// safeInt 读取可空 int 指针并提供默认值。
func safeInt(v *int, def int) int {
	if v == nil {
		return def
	}
	return *v
}

// safeInt64 读取可空 int64 指针并提供默认值。
func safeInt64(v *int64, def int64) int64 {
	if v == nil {
		return def
	}
	return *v
}

// valueOrCurrentUint 当请求值为空时回落到当前值。
func valueOrCurrentUint(v *uint, current uint) *uint {
	if v != nil {
		return v
	}
	return &current
}

// valueOrCurrentString 当请求值为空时回落到当前值。
func valueOrCurrentString(v *string, current string) *string {
	if v != nil {
		return v
	}
	return &current
}

// normalizeOptionNames 去重清洗动物竞猜选项，保持输入顺序。
func normalizeOptionNames(input []string) []string {
	out := make([]string, 0, len(input))
	seen := make(map[string]struct{}, len(input))
	for _, raw := range input {
		name := strings.TrimSpace(raw)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

// defaultAnimalOptionNames 返回默认 12 生肖竞猜选项。
func defaultAnimalOptionNames() []string {
	return []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
}
