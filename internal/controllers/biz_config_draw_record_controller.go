package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go-admin-full/internal/models"
	"go-admin-full/internal/utils"
	"gorm.io/gorm"
)

// drawRecordUpsertRequest 开奖区记录新增/编辑请求结构。
type drawRecordUpsertRequest struct {
	SpecialLotteryID    *uint   `json:"special_lottery_id"`
	Issue               *string `json:"issue"`
	Year                *int    `json:"year"`
	DrawAt              *string `json:"draw_at"`
	NormalDrawResult    *string `json:"normal_draw_result"`
	SpecialDrawResult   *string `json:"special_draw_result"`
	DrawResult          *string `json:"draw_result"`
	DrawLabels          *string `json:"draw_labels"`
	PlaybackURL         *string `json:"playback_url"`
	SpecialSingleDouble *string `json:"special_single_double"`
	SpecialBigSmall     *string `json:"special_big_small"`
	SumSingleDouble     *string `json:"sum_single_double"`
	SumBigSmall         *string `json:"sum_big_small"`
	RecommendSix        *string `json:"recommend_six"`
	RecommendFour       *string `json:"recommend_four"`
	RecommendOne        *string `json:"recommend_one"`
	RecommendTen        *string `json:"recommend_ten"`
	SpecialCode         *string `json:"special_code"`
	NormalCode          *string `json:"normal_code"`
	Zheng1              *string `json:"zheng1"`
	Zheng2              *string `json:"zheng2"`
	Zheng3              *string `json:"zheng3"`
	Zheng4              *string `json:"zheng4"`
	Zheng5              *string `json:"zheng5"`
	Zheng6              *string `json:"zheng6"`
	IsCurrent           *int8   `json:"is_current"`
	Status              *int8   `json:"status"`
	Sort                *int    `json:"sort"`
}

// ListDrawRecords 查询开奖区开奖记录列表。
func (bc *BizConfigController) ListDrawRecords(c *gin.Context) {
	// 1) 读取筛选参数（彩种ID + 关键字）。
	specialLotteryID := strings.TrimSpace(c.Query("special_lottery_id"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	limit := parseIntWithDefault(c.Query("limit"), 300)
	if limit <= 0 || limit > 1000 {
		limit = 300
	}

	// 2) 组装查询条件，按开奖时间倒序输出。
	query := bc.db.Model(&models.WDrawRecord{}).Order("draw_at DESC, id DESC").Limit(limit)
	if sid, err := strconv.Atoi(specialLotteryID); err == nil && sid > 0 {
		query = query.Where("special_lottery_id = ?", sid)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("issue LIKE ?", like)
	}

	// 3) 执行查询并返回。
	var items []models.WDrawRecord
	if err := query.Find(&items).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"items": items})
}

// CreateDrawRecord 新增开奖区开奖记录。
func (bc *BizConfigController) CreateDrawRecord(c *gin.Context) {
	var req drawRecordUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}
	// 开奖区记录必须具备彩种ID和期号。
	if req.SpecialLotteryID == nil || *req.SpecialLotteryID == 0 || strings.TrimSpace(safeString(req.Issue)) == "" {
		utils.JSONError(c, http.StatusBadRequest, "special_lottery_id/issue required")
		return
	}

	// 6+1 开奖号码为必填，并统一生成兼容字段 draw_result。
	normalRaw, specialRaw, mergedRaw, err := normalizeAndMergeDrawNumbers(
		safeString(req.NormalDrawResult),
		safeString(req.SpecialDrawResult),
		safeString(req.DrawResult),
	)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}
	drawAt := parseDateTimeOrDefault(safeString(req.DrawAt), time.Now())
	year := drawAt.Year()
	if req.Year != nil && *req.Year > 0 {
		year = *req.Year
	}

	// 标签为空时自动按号码生成默认“生肖/五行”标签，并同步拆分出独立属相/五行字段。
	labels, zodiacLabels, wuxingLabels := normalizeDrawLabels(safeString(req.DrawLabels), mergedRaw)

	// 自动计算开奖结果详情指标（当请求未填写时兜底）。
	stats := deriveDrawStats(mergedRaw)

	item := models.WDrawRecord{
		SpecialLotteryID:    *req.SpecialLotteryID,
		Issue:               strings.TrimSpace(safeString(req.Issue)),
		Year:                year,
		DrawAt:              drawAt,
		NormalDrawResult:    normalRaw,
		SpecialDrawResult:   specialRaw,
		DrawResult:          mergedRaw,
		DrawLabels:          labels,
		ZodiacLabels:        zodiacLabels,
		WuxingLabels:        wuxingLabels,
		PlaybackURL:         safeString(req.PlaybackURL),
		SpecialSingleDouble: valueOrAuto(req.SpecialSingleDouble, stats.SpecialSingleDouble),
		SpecialBigSmall:     valueOrAuto(req.SpecialBigSmall, stats.SpecialBigSmall),
		SumSingleDouble:     valueOrAuto(req.SumSingleDouble, stats.SumSingleDouble),
		SumBigSmall:         valueOrAuto(req.SumBigSmall, stats.SumBigSmall),
		RecommendSix:        safeString(req.RecommendSix),
		RecommendFour:       safeString(req.RecommendFour),
		RecommendOne:        safeString(req.RecommendOne),
		RecommendTen:        safeString(req.RecommendTen),
		SpecialCode:         valueOrAuto(req.SpecialCode, stats.SpecialCode),
		NormalCode:          valueOrAuto(req.NormalCode, stats.NormalCode),
		Zheng1:              safeString(req.Zheng1),
		Zheng2:              safeString(req.Zheng2),
		Zheng3:              safeString(req.Zheng3),
		Zheng4:              safeString(req.Zheng4),
		Zheng5:              safeString(req.Zheng5),
		Zheng6:              safeString(req.Zheng6),
		IsCurrent:           safeInt8(req.IsCurrent, 0),
		Status:              safeInt8(req.Status, 1),
		Sort:                safeInt(req.Sort, 0),
	}

	// 同一彩种只允许一条当前期记录。
	if err := bc.db.Transaction(func(tx *gorm.DB) error {
		if item.IsCurrent == 1 {
			if err := tx.Model(&models.WDrawRecord{}).
				Where("special_lottery_id = ?", item.SpecialLotteryID).
				Update("is_current", 0).Error; err != nil {
				return err
			}
		}
		return tx.Create(&item).Error
	}); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}

	utils.JSONOK(c, item)
}

// UpdateDrawRecord 编辑开奖区开奖记录。
func (bc *BizConfigController) UpdateDrawRecord(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req drawRecordUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid request")
		return
	}

	var current models.WDrawRecord
	if err := bc.db.First(&current, id).Error; err != nil {
		utils.JSONError(c, http.StatusNotFound, "draw record not found")
		return
	}

	// 1) 基于当前值做增量合并。
	next := current
	if req.SpecialLotteryID != nil && *req.SpecialLotteryID > 0 {
		next.SpecialLotteryID = *req.SpecialLotteryID
	}
	if req.Issue != nil {
		next.Issue = strings.TrimSpace(*req.Issue)
	}
	if req.DrawAt != nil {
		next.DrawAt = parseDateTimeOrDefault(*req.DrawAt, next.DrawAt)
	}
	if req.Year != nil && *req.Year > 0 {
		next.Year = *req.Year
	}
	if req.PlaybackURL != nil {
		next.PlaybackURL = strings.TrimSpace(*req.PlaybackURL)
	}
	if req.RecommendSix != nil {
		next.RecommendSix = strings.TrimSpace(*req.RecommendSix)
	}
	if req.RecommendFour != nil {
		next.RecommendFour = strings.TrimSpace(*req.RecommendFour)
	}
	if req.RecommendOne != nil {
		next.RecommendOne = strings.TrimSpace(*req.RecommendOne)
	}
	if req.RecommendTen != nil {
		next.RecommendTen = strings.TrimSpace(*req.RecommendTen)
	}
	if req.Zheng1 != nil {
		next.Zheng1 = strings.TrimSpace(*req.Zheng1)
	}
	if req.Zheng2 != nil {
		next.Zheng2 = strings.TrimSpace(*req.Zheng2)
	}
	if req.Zheng3 != nil {
		next.Zheng3 = strings.TrimSpace(*req.Zheng3)
	}
	if req.Zheng4 != nil {
		next.Zheng4 = strings.TrimSpace(*req.Zheng4)
	}
	if req.Zheng5 != nil {
		next.Zheng5 = strings.TrimSpace(*req.Zheng5)
	}
	if req.Zheng6 != nil {
		next.Zheng6 = strings.TrimSpace(*req.Zheng6)
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

	// 2) 号码字段变更时，重新校验并重建 draw_result。
	if req.NormalDrawResult != nil || req.SpecialDrawResult != nil || req.DrawResult != nil {
		normalRaw, specialRaw, mergedRaw, drawErr := normalizeAndMergeDrawNumbers(
			safeString(valueOrCurrentString(req.NormalDrawResult, next.NormalDrawResult)),
			safeString(valueOrCurrentString(req.SpecialDrawResult, next.SpecialDrawResult)),
			safeString(valueOrCurrentString(req.DrawResult, next.DrawResult)),
		)
		if drawErr != nil {
			utils.JSONError(c, http.StatusBadRequest, drawErr.Error())
			return
		}
		next.NormalDrawResult = normalRaw
		next.SpecialDrawResult = specialRaw
		next.DrawResult = mergedRaw
	}

	// 3) 标签字段：传值按传值；未传值时基于最新号码自动生成，并同步独立属相/五行字段。
	if req.DrawLabels != nil {
		next.DrawLabels, next.ZodiacLabels, next.WuxingLabels = normalizeDrawLabels(strings.TrimSpace(*req.DrawLabels), next.DrawResult)
	} else {
		next.DrawLabels, next.ZodiacLabels, next.WuxingLabels = normalizeDrawLabels(next.DrawLabels, next.DrawResult)
	}

	// 4) 自动补齐开奖详情基础指标。
	stats := deriveDrawStats(next.DrawResult)
	next.SpecialSingleDouble = valueOrKeep(req.SpecialSingleDouble, next.SpecialSingleDouble, stats.SpecialSingleDouble)
	next.SpecialBigSmall = valueOrKeep(req.SpecialBigSmall, next.SpecialBigSmall, stats.SpecialBigSmall)
	next.SumSingleDouble = valueOrKeep(req.SumSingleDouble, next.SumSingleDouble, stats.SumSingleDouble)
	next.SumBigSmall = valueOrKeep(req.SumBigSmall, next.SumBigSmall, stats.SumBigSmall)
	next.SpecialCode = valueOrKeep(req.SpecialCode, next.SpecialCode, stats.SpecialCode)
	next.NormalCode = valueOrKeep(req.NormalCode, next.NormalCode, stats.NormalCode)

	// 5) 核心字段二次校验。
	if next.SpecialLotteryID == 0 || strings.TrimSpace(next.Issue) == "" {
		utils.JSONError(c, http.StatusBadRequest, "special_lottery_id/issue required")
		return
	}

	updates := map[string]interface{}{
		"special_lottery_id":    next.SpecialLotteryID,
		"issue":                 next.Issue,
		"year":                  next.Year,
		"draw_at":               next.DrawAt,
		"normal_draw_result":    next.NormalDrawResult,
		"special_draw_result":   next.SpecialDrawResult,
		"draw_result":           next.DrawResult,
		"draw_labels":           next.DrawLabels,
		"zodiac_labels":         next.ZodiacLabels,
		"wuxing_labels":         next.WuxingLabels,
		"playback_url":          next.PlaybackURL,
		"special_single_double": next.SpecialSingleDouble,
		"special_big_small":     next.SpecialBigSmall,
		"sum_single_double":     next.SumSingleDouble,
		"sum_big_small":         next.SumBigSmall,
		"recommend_six":         next.RecommendSix,
		"recommend_four":        next.RecommendFour,
		"recommend_one":         next.RecommendOne,
		"recommend_ten":         next.RecommendTen,
		"special_code":          next.SpecialCode,
		"normal_code":           next.NormalCode,
		"zheng1":                next.Zheng1,
		"zheng2":                next.Zheng2,
		"zheng3":                next.Zheng3,
		"zheng4":                next.Zheng4,
		"zheng5":                next.Zheng5,
		"zheng6":                next.Zheng6,
		"is_current":            next.IsCurrent,
		"status":                next.Status,
		"sort":                  next.Sort,
	}

	// 6) 写库并维护“同彩种唯一当前期”约束。
	if err := bc.db.Transaction(func(tx *gorm.DB) error {
		if next.IsCurrent == 1 {
			if err := tx.Model(&models.WDrawRecord{}).
				Where("special_lottery_id = ? AND id <> ?", next.SpecialLotteryID, id).
				Update("is_current", 0).Error; err != nil {
				return err
			}
		}
		return tx.Model(&models.WDrawRecord{}).Where("id = ?", id).Updates(updates).Error
	}); err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// DeleteDrawRecord 删除开奖区开奖记录。
func (bc *BizConfigController) DeleteDrawRecord(c *gin.Context) {
	id, err := parseUintID(c)
	if err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := bc.db.Delete(&models.WDrawRecord{}, id).Error; err != nil {
		utils.JSONError(c, 500, err.Error())
		return
	}
	utils.JSONOK(c, gin.H{"id": id})
}

// drawStats 开奖详情自动计算结果。
type drawStats struct {
	SpecialSingleDouble string
	SpecialBigSmall     string
	SumSingleDouble     string
	SumBigSmall         string
	SpecialCode         string
	NormalCode          string
}

// deriveDrawStats 根据 6+1 开奖号码计算详情页基础字段。
func deriveDrawStats(drawResult string) drawStats {
	nums, err := parseDrawNumbers(drawResult)
	if err != nil || len(nums) < 7 {
		return drawStats{}
	}
	special := nums[6]
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return drawStats{
		SpecialSingleDouble: oddEvenCN(special),
		SpecialBigSmall:     bigSmallCN(special, 24),
		SumSingleDouble:     oddEvenCN(sum),
		SumBigSmall:         bigSmallCN(sum, 175),
		SpecialCode:         strconv.Itoa(special),
		NormalCode:          joinIntCSV(nums[:6]),
	}
}

// oddEvenCN 返回中文单双。
func oddEvenCN(v int) string {
	if v%2 == 0 {
		return "双"
	}
	return "单"
}

// bigSmallCN 按阈值返回中文大小（> threshold 为大）。
func bigSmallCN(v, threshold int) string {
	if v > threshold {
		return "大"
	}
	return "小"
}

// normalizeDrawLabels 清洗开奖标签；标签缺失时自动生成 7 个默认标签，并输出独立属相/五行串。
func normalizeDrawLabels(raw, drawResult string) (string, string, string) {
	labels := parseLabelCSV(raw)
	if len(labels) == 7 {
		return buildSplitLabelResult(labels, drawResult)
	}
	nums, err := parseDrawNumbers(drawResult)
	if err != nil || len(nums) != 7 {
		joined := strings.Join(labels, ",")
		return joined, joined, ""
	}
	return buildSplitLabelResult(buildDefaultDrawLabels(nums), drawResult)
}

// buildSplitLabelResult 将组合标签拆分为“draw_labels/zodiac_labels/wuxing_labels”三种存储格式。
func buildSplitLabelResult(labels []string, drawResult string) (string, string, string) {
	// 1) 预先准备号码默认标签，防止人工只填属相不填五行导致数据不完整。
	defaultPairs := map[int]string{}
	nums := numsOrEmpty(drawResult)
	if len(nums) == 7 {
		pairs := buildDefaultDrawLabels(nums)
		for idx, n := range nums {
			if idx < len(pairs) {
				defaultPairs[n] = pairs[idx]
			}
		}
	}

	// 2) 按顺序拆分每个标签的属相/五行。
	drawOut := make([]string, 0, len(labels))
	zodiacOut := make([]string, 0, len(labels))
	wuxingOut := make([]string, 0, len(labels))
	for idx, item := range labels {
		raw := strings.TrimSpace(item)
		if raw == "" {
			raw = fallbackPairByIndex(defaultPairs, nums, idx)
		}
		zodiac, wuxing := splitPair(raw)
		if wuxing == "" {
			// 3) 缺五行时回退默认映射，保证开奖记录字段完整。
			fallback := fallbackPairByIndex(defaultPairs, nums, idx)
			fz, fw := splitPair(fallback)
			if zodiac == "" {
				zodiac = fz
			}
			wuxing = fw
		}
		if zodiac == "" {
			zodiac = raw
		}
		drawLabel := zodiac
		if wuxing != "" {
			drawLabel = strings.TrimSpace(zodiac + "/" + wuxing)
		}
		drawOut = append(drawOut, drawLabel)
		zodiacOut = append(zodiacOut, zodiac)
		wuxingOut = append(wuxingOut, wuxing)
	}
	return strings.Join(drawOut, ","), strings.Join(zodiacOut, ","), strings.Join(wuxingOut, ",")
}

// splitPair 将“属相/五行”标签拆成两个片段。
func splitPair(raw string) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(raw), "/", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(raw), ""
}

// numsOrEmpty 解析 drawResult，失败则返回空切片。
func numsOrEmpty(drawResult string) []int {
	nums, err := parseDrawNumbers(drawResult)
	if err != nil {
		return []int{}
	}
	return nums
}

// fallbackPairByIndex 根据号码顺序取默认“属相/五行”标签。
func fallbackPairByIndex(pairMap map[int]string, nums []int, idx int) string {
	if idx < 0 || idx >= len(nums) {
		return ""
	}
	return pairMap[nums[idx]]
}

// parseLabelCSV 解析标签字符串（逗号/空格/换行分隔）。
func parseLabelCSV(raw string) []string {
	tokens := strings.FieldsFunc(strings.TrimSpace(raw), func(r rune) bool {
		return r == ',' || r == '|' || r == ';' || r == '\n' || r == '\r' || r == '\t'
	})
	out := make([]string, 0, len(tokens))
	for _, t := range tokens {
		v := strings.TrimSpace(t)
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}

// buildDefaultDrawLabels 基于号码生成默认“生肖/五行”标签。
func buildDefaultDrawLabels(nums []int) []string {
	zodiacs := []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	elements := []string{"金", "木", "水", "火", "土"}
	out := make([]string, 0, len(nums))
	for _, n := range nums {
		zodiac := zodiacs[(n-1)%len(zodiacs)]
		element := elements[(n-1)%len(elements)]
		out = append(out, zodiac+"/"+element)
	}
	return out
}

// valueOrAuto 请求传值优先，否则回退自动计算值。
func valueOrAuto(v *string, auto string) string {
	manual := safeString(v)
	if manual != "" {
		return manual
	}
	return auto
}

// valueOrKeep 请求传值优先；否则保留当前值；当前值为空时回退自动值。
func valueOrKeep(v *string, current, auto string) string {
	if v != nil {
		manual := strings.TrimSpace(*v)
		if manual != "" {
			return manual
		}
	}
	if strings.TrimSpace(current) != "" {
		return current
	}
	return auto
}

// parseIntWithDefault 将字符串解析为整数，失败时回退默认值。
func parseIntWithDefault(raw string, def int) int {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return def
	}
	return v
}
