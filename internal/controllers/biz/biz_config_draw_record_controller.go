package biz

import (
	"errors"
	"strconv"
	"strings"
	"time"

	commonresp "github.com/wangyahua6688-maker/tk-common/utils/httpresp"
	"go-admin/internal/constants"
	admindto "go-admin/internal/dto/admin"
	"go-admin/internal/models"
	bizsvc "go-admin/internal/services/biz"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListDrawRecords 查询开奖区开奖记录列表。
func (bc *LotteryController) ListDrawRecords(c *gin.Context) {
	// 1) 读取筛选参数（彩种ID + 关键字）。
	specialLotteryID := strings.TrimSpace(c.Query("special_lottery_id"))
	// 定义并初始化当前变量。
	keyword := strings.TrimSpace(c.Query("keyword"))
	// 定义并初始化当前变量。
	limit := parseIntWithDefault(c.Query("limit"), 300)
	// 判断条件并进入对应分支逻辑。
	if limit <= 0 || limit > 1000 {
		// 更新当前变量或字段值。
		limit = 300
	}

	// 2) 组装查询条件并调用服务层。
	filter := bizsvc.DrawRecordFilter{
		// 处理当前语句逻辑。
		Limit: limit,
		// 处理当前语句逻辑。
		Keyword: keyword,
	}
	// 判断条件并进入对应分支逻辑。
	if sid, err := strconv.Atoi(specialLotteryID); err == nil && sid > 0 {
		// 更新当前变量或字段值。
		filter.SpecialLotteryID = uint(sid)
	}

	// 3) 执行查询并返回。
	items, err := bc.drawRecordSvc.ListDrawRecords(c.Request.Context(), filter)
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

// CreateDrawRecord 新增开奖区开奖记录。
func (bc *LotteryController) CreateDrawRecord(c *gin.Context) {
	// 声明当前变量。
	var req admindto.DrawRecordUpsertRequest
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}
	// 开奖区记录必须具备彩种ID和期号。
	if req.SpecialLotteryID == nil || *req.SpecialLotteryID == 0 || strings.TrimSpace(safeString(req.Issue)) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "special_lottery_id/issue required")
		// 返回当前处理结果。
		return
	}

	// 6+1 开奖号码为必填，并统一生成兼容字段 draw_result。
	normalRaw, specialRaw, mergedRaw, err := normalizeAndMergeDrawNumbers(
		// 调用safeString完成当前处理。
		safeString(req.NormalDrawResult),
		// 调用safeString完成当前处理。
		safeString(req.SpecialDrawResult),
		// 调用safeString完成当前处理。
		safeString(req.DrawResult),
	)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, err.Error())
		// 返回当前处理结果。
		return
	}
	// 定义并初始化当前变量。
	drawAt := parseDateTimeOrDefault(safeString(req.DrawAt), time.Now())
	// 定义并初始化当前变量。
	year := drawAt.Year()
	// 判断条件并进入对应分支逻辑。
	if req.Year != nil && *req.Year > 0 {
		// 更新当前变量或字段值。
		year = *req.Year
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
	item := models.WDrawRecord{
		// 处理当前语句逻辑。
		SpecialLotteryID: *req.SpecialLotteryID,
		// 调用strings.TrimSpace完成当前处理。
		Issue: strings.TrimSpace(safeString(req.Issue)),
		// 处理当前语句逻辑。
		Year: year,
		// 处理当前语句逻辑。
		DrawAt: drawAt,
		// 处理当前语句逻辑。
		NormalDrawResult: normalRaw,
		// 处理当前语句逻辑。
		SpecialDrawResult: specialRaw,
		// 处理当前语句逻辑。
		DrawResult: mergedRaw,
		// 处理当前语句逻辑。
		PlaybackURL: playbackURL,
		// 调用safeString完成当前处理。
		RecommendSix: safeString(req.RecommendSix),
		// 调用safeString完成当前处理。
		RecommendFour: safeString(req.RecommendFour),
		// 调用safeString完成当前处理。
		RecommendOne: safeString(req.RecommendOne),
		// 调用safeString完成当前处理。
		RecommendTen: safeString(req.RecommendTen),
		// 调用safeInt8完成当前处理。
		IsCurrent: safeInt8(req.IsCurrent, 0),
		// 调用safeInt8完成当前处理。
		Status: safeInt8(req.Status, 1),
		// 调用safeInt完成当前处理。
		Sort: safeInt(req.Sort, 0),
	}

	// 同一彩种只允许一条当前期记录。
	if err := bc.drawRecordSvc.CreateDrawRecord(c.Request.Context(), &item); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 开奖记录已变更，需要同步清公开首页/开奖区/开奖现场缓存。
	_ = invalidatePublicLotteryCaches(c.Request.Context(), item.SpecialLotteryID)

	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, item)
}

// UpdateDrawRecord 编辑开奖区开奖记录。
func (bc *LotteryController) UpdateDrawRecord(c *gin.Context) {
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
	var req admindto.DrawRecordUpsertRequest
	// 判断条件并进入对应分支逻辑。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid request")
		// 返回当前处理结果。
		return
	}

	// 定义并初始化当前变量。
	current, err := bc.drawRecordSvc.GetDrawRecordByID(c.Request.Context(), id)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 判断条件并进入对应分支逻辑。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizResourceNotFound, "draw record not found")
			// 返回当前处理结果。
			return
		}
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if current == nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizResourceNotFound, "draw record not found")
		// 返回当前处理结果。
		return
	}

	// 1) 基于当前值做增量合并。
	next := *current
	// 判断条件并进入对应分支逻辑。
	if req.SpecialLotteryID != nil && *req.SpecialLotteryID > 0 {
		// 更新当前变量或字段值。
		next.SpecialLotteryID = *req.SpecialLotteryID
	}
	// 判断条件并进入对应分支逻辑。
	if req.Issue != nil {
		// 更新当前变量或字段值。
		next.Issue = strings.TrimSpace(*req.Issue)
	}
	// 判断条件并进入对应分支逻辑。
	if req.DrawAt != nil {
		// 更新当前变量或字段值。
		next.DrawAt = parseDateTimeOrDefault(*req.DrawAt, next.DrawAt)
		// 未显式传年份时，默认跟随开奖时间年份变化。
		next.Year = next.DrawAt.Year()
	}
	// 判断条件并进入对应分支逻辑。
	if req.Year != nil && *req.Year > 0 {
		// 更新当前变量或字段值。
		next.Year = *req.Year
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
	if req.RecommendSix != nil {
		// 更新当前变量或字段值。
		next.RecommendSix = strings.TrimSpace(*req.RecommendSix)
	}
	// 判断条件并进入对应分支逻辑。
	if req.RecommendFour != nil {
		// 更新当前变量或字段值。
		next.RecommendFour = strings.TrimSpace(*req.RecommendFour)
	}
	// 判断条件并进入对应分支逻辑。
	if req.RecommendOne != nil {
		// 更新当前变量或字段值。
		next.RecommendOne = strings.TrimSpace(*req.RecommendOne)
	}
	// 判断条件并进入对应分支逻辑。
	if req.RecommendTen != nil {
		// 更新当前变量或字段值。
		next.RecommendTen = strings.TrimSpace(*req.RecommendTen)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Zheng1 != nil {
		// 更新当前变量或字段值。
		next.Zheng1 = strings.TrimSpace(*req.Zheng1)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Zheng2 != nil {
		// 更新当前变量或字段值。
		next.Zheng2 = strings.TrimSpace(*req.Zheng2)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Zheng3 != nil {
		// 更新当前变量或字段值。
		next.Zheng3 = strings.TrimSpace(*req.Zheng3)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Zheng4 != nil {
		// 更新当前变量或字段值。
		next.Zheng4 = strings.TrimSpace(*req.Zheng4)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Zheng5 != nil {
		// 更新当前变量或字段值。
		next.Zheng5 = strings.TrimSpace(*req.Zheng5)
	}
	// 判断条件并进入对应分支逻辑。
	if req.Zheng6 != nil {
		// 更新当前变量或字段值。
		next.Zheng6 = strings.TrimSpace(*req.Zheng6)
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

	// 2) 号码字段变更时，重新校验并重建 draw_result。
	if req.NormalDrawResult != nil || req.SpecialDrawResult != nil || req.DrawResult != nil {
		// 定义并初始化当前变量。
		normalRaw, specialRaw, mergedRaw, drawErr := normalizeAndMergeDrawNumbers(
			// 调用safeString完成当前处理。
			safeString(valueOrCurrentString(req.NormalDrawResult, next.NormalDrawResult)),
			// 调用safeString完成当前处理。
			safeString(valueOrCurrentString(req.SpecialDrawResult, next.SpecialDrawResult)),
			// 调用safeString完成当前处理。
			safeString(valueOrCurrentString(req.DrawResult, next.DrawResult)),
		)
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

	// 3) 核心字段二次校验。
	if next.SpecialLotteryID == 0 || strings.TrimSpace(next.Issue) == "" {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "special_lottery_id/issue required")
		// 返回当前处理结果。
		return
	}

	// 4) 写库并维护“同彩种唯一当前期”约束。
	if err := bc.drawRecordSvc.UpdateDrawRecord(c.Request.Context(), id, &next); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 彩种可能被修改，需要同时失效旧彩种和新彩种缓存。
	_ = invalidatePublicLotteryCaches(c.Request.Context(), current.SpecialLotteryID)
	if next.SpecialLotteryID != current.SpecialLotteryID {
		_ = invalidatePublicLotteryCaches(c.Request.Context(), next.SpecialLotteryID)
	}
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// DeleteDrawRecord 删除开奖区开奖记录。
func (bc *LotteryController) DeleteDrawRecord(c *gin.Context) {
	// 定义并初始化当前变量。
	id, err := parseUintID(c)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminBizInvalidRequest, "invalid id")
		// 返回当前处理结果。
		return
	}
	// 删除前先取当前记录，用于删除后清理对应彩种缓存。
	current, err := bc.drawRecordSvc.GetDrawRecordByID(c.Request.Context(), id)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 判断条件并进入对应分支逻辑。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 调用utils.JSONError完成当前处理。
			commonresp.GinError(c, constants.AdminBizResourceNotFound, "draw record not found")
			// 返回当前处理结果。
			return
		}
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 判断条件并进入对应分支逻辑。
	if err := bc.drawRecordSvc.DeleteDrawRecord(c.Request.Context(), id); err != nil {
		// 调用utils.JSONError完成当前处理。
		commonresp.GinError(c, constants.AdminSysInternalError, err.Error())
		// 返回当前处理结果。
		return
	}
	// 删除成功后，清理首页/开奖区/开奖现场缓存。
	_ = invalidatePublicLotteryCaches(c.Request.Context(), current.SpecialLotteryID)
	// 调用utils.JSONOK完成当前处理。
	commonresp.GinOK(c, gin.H{"id": id})
}

// drawStats 开奖详情自动计算结果。
type drawStats struct {
	// 处理当前语句逻辑。
	SpecialSingleDouble string
	// 处理当前语句逻辑。
	SpecialBigSmall string
	// 处理当前语句逻辑。
	SumSingleDouble string
	// 处理当前语句逻辑。
	SumBigSmall string
	// 处理当前语句逻辑。
	SpecialCode string
	// 处理当前语句逻辑。
	NormalCode string
}

// deriveDrawStats 根据 6+1 开奖号码计算详情页基础字段。
func deriveDrawStats(drawResult string) drawStats {
	// 定义并初始化当前变量。
	nums, err := parseDrawNumbers(drawResult)
	// 判断条件并进入对应分支逻辑。
	if err != nil || len(nums) < 7 {
		// 返回当前处理结果。
		return drawStats{}
	}
	// 定义并初始化当前变量。
	special := nums[6]
	// 定义并初始化当前变量。
	sum := 0
	// 循环处理当前数据集合。
	for _, n := range nums {
		// 更新当前变量或字段值。
		sum += n
	}
	// 返回当前处理结果。
	return drawStats{
		// 调用oddEvenCN完成当前处理。
		SpecialSingleDouble: oddEvenCN(special),
		// 调用bigSmallCN完成当前处理。
		SpecialBigSmall: bigSmallCN(special, 24),
		// 调用oddEvenCN完成当前处理。
		SumSingleDouble: oddEvenCN(sum),
		// 调用bigSmallCN完成当前处理。
		SumBigSmall: bigSmallCN(sum, 175),
		// 调用strconv.Itoa完成当前处理。
		SpecialCode: strconv.Itoa(special),
		// 调用joinIntCSV完成当前处理。
		NormalCode: joinIntCSV(nums[:6]),
	}
}

// oddEvenCN 返回中文单双。
func oddEvenCN(v int) string {
	// 判断条件并进入对应分支逻辑。
	if v%2 == 0 {
		return "双"
	}
	return "单"
}

// bigSmallCN 按阈值返回中文大小（> threshold 为大）。
func bigSmallCN(v, threshold int) string {
	// 判断条件并进入对应分支逻辑。
	if v > threshold {
		return "大"
	}
	return "小"
}

// normalizeDrawLabels 清洗开奖标签；标签缺失时自动生成 7 个默认标签，并输出独立属相/五行串。
func normalizeDrawLabels(raw, drawResult string) (string, string, string) {
	// 定义并初始化当前变量。
	labels := parseLabelCSV(raw)
	// 判断条件并进入对应分支逻辑。
	if len(labels) == 7 {
		// 返回当前处理结果。
		return buildSplitLabelResult(labels, drawResult)
	}
	// 定义并初始化当前变量。
	nums, err := parseDrawNumbers(drawResult)
	// 判断条件并进入对应分支逻辑。
	if err != nil || len(nums) != 7 {
		// 定义并初始化当前变量。
		joined := strings.Join(labels, ",")
		// 返回当前处理结果。
		return joined, joined, ""
	}
	// 返回当前处理结果。
	return buildSplitLabelResult(buildDefaultDrawLabels(nums), drawResult)
}

// buildSplitLabelResult 将组合标签拆分为“draw_labels/zodiac_labels/wuxing_labels”三种存储格式。
func buildSplitLabelResult(labels []string, drawResult string) (string, string, string) {
	// 1) 预先准备号码默认标签，防止人工只填属相不填五行导致数据不完整。
	defaultPairs := map[int]string{}
	// 定义并初始化当前变量。
	nums := numsOrEmpty(drawResult)
	// 判断条件并进入对应分支逻辑。
	if len(nums) == 7 {
		// 定义并初始化当前变量。
		pairs := buildDefaultDrawLabels(nums)
		// 循环处理当前数据集合。
		for idx, n := range nums {
			// 判断条件并进入对应分支逻辑。
			if idx < len(pairs) {
				// 更新当前变量或字段值。
				defaultPairs[n] = pairs[idx]
			}
		}
	}

	// 2) 按顺序拆分每个标签的属相/五行。
	drawOut := make([]string, 0, len(labels))
	// 定义并初始化当前变量。
	zodiacOut := make([]string, 0, len(labels))
	// 定义并初始化当前变量。
	wuxingOut := make([]string, 0, len(labels))
	// 循环处理当前数据集合。
	for idx, item := range labels {
		// 定义并初始化当前变量。
		raw := strings.TrimSpace(item)
		// 判断条件并进入对应分支逻辑。
		if raw == "" {
			// 更新当前变量或字段值。
			raw = fallbackPairByIndex(defaultPairs, nums, idx)
		}
		// 定义并初始化当前变量。
		zodiac, wuxing := splitPair(raw)
		// 判断条件并进入对应分支逻辑。
		if wuxing == "" {
			// 3) 缺五行时回退默认映射，保证开奖记录字段完整。
			fallback := fallbackPairByIndex(defaultPairs, nums, idx)
			// 定义并初始化当前变量。
			fz, fw := splitPair(fallback)
			// 判断条件并进入对应分支逻辑。
			if zodiac == "" {
				// 更新当前变量或字段值。
				zodiac = fz
			}
			// 更新当前变量或字段值。
			wuxing = fw
		}
		// 判断条件并进入对应分支逻辑。
		if zodiac == "" {
			// 更新当前变量或字段值。
			zodiac = raw
		}
		// 定义并初始化当前变量。
		drawLabel := zodiac
		// 判断条件并进入对应分支逻辑。
		if wuxing != "" {
			// 更新当前变量或字段值。
			drawLabel = strings.TrimSpace(zodiac + "/" + wuxing)
		}
		// 更新当前变量或字段值。
		drawOut = append(drawOut, drawLabel)
		// 更新当前变量或字段值。
		zodiacOut = append(zodiacOut, zodiac)
		// 更新当前变量或字段值。
		wuxingOut = append(wuxingOut, wuxing)
	}
	// 返回当前处理结果。
	return strings.Join(drawOut, ","), strings.Join(zodiacOut, ","), strings.Join(wuxingOut, ",")
}

// splitPair 将“属相/五行”标签拆成两个片段。
func splitPair(raw string) (string, string) {
	// 定义并初始化当前变量。
	parts := strings.SplitN(strings.TrimSpace(raw), "/", 2)
	// 判断条件并进入对应分支逻辑。
	if len(parts) == 2 {
		// 返回当前处理结果。
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	// 返回当前处理结果。
	return strings.TrimSpace(raw), ""
}

// numsOrEmpty 解析 drawResult，失败则返回空切片。
func numsOrEmpty(drawResult string) []int {
	// 定义并初始化当前变量。
	nums, err := parseDrawNumbers(drawResult)
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return []int{}
	}
	// 返回当前处理结果。
	return nums
}

// fallbackPairByIndex 根据号码顺序取默认“属相/五行”标签。
func fallbackPairByIndex(pairMap map[int]string, nums []int, idx int) string {
	// 判断条件并进入对应分支逻辑。
	if idx < 0 || idx >= len(nums) {
		// 返回当前处理结果。
		return ""
	}
	// 返回当前处理结果。
	return pairMap[nums[idx]]
}

// parseLabelCSV 解析标签字符串（逗号/空格/换行分隔）。
func parseLabelCSV(raw string) []string {
	// 定义并初始化当前变量。
	tokens := strings.FieldsFunc(strings.TrimSpace(raw), func(r rune) bool {
		// 返回当前处理结果。
		return r == ',' || r == '|' || r == ';' || r == '\n' || r == '\r' || r == '\t'
	})
	// 定义并初始化当前变量。
	out := make([]string, 0, len(tokens))
	// 循环处理当前数据集合。
	for _, t := range tokens {
		// 定义并初始化当前变量。
		v := strings.TrimSpace(t)
		// 判断条件并进入对应分支逻辑。
		if v == "" {
			// 处理当前语句逻辑。
			continue
		}
		// 更新当前变量或字段值。
		out = append(out, v)
	}
	// 返回当前处理结果。
	return out
}

// buildDefaultDrawLabels 基于号码生成默认“生肖/五行”标签。
func buildDefaultDrawLabels(nums []int) []string {
	// 定义并初始化当前变量。
	zodiacs := []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	// 定义并初始化当前变量。
	elements := []string{"金", "木", "水", "火", "土"}
	// 定义并初始化当前变量。
	out := make([]string, 0, len(nums))
	// 循环处理当前数据集合。
	for _, n := range nums {
		// 定义并初始化当前变量。
		zodiac := zodiacs[(n-1)%len(zodiacs)]
		// 定义并初始化当前变量。
		element := elements[(n-1)%len(elements)]
		// 更新当前变量或字段值。
		out = append(out, zodiac+"/"+element)
	}
	// 返回当前处理结果。
	return out
}

// valueOrAuto 请求传值优先，否则回退自动计算值。
func valueOrAuto(v *string, auto string) string {
	// 定义并初始化当前变量。
	manual := safeString(v)
	// 判断条件并进入对应分支逻辑。
	if manual != "" {
		// 返回当前处理结果。
		return manual
	}
	// 返回当前处理结果。
	return auto
}

// valueOrKeep 请求传值优先；否则保留当前值；当前值为空时回退自动值。
func valueOrKeep(v *string, current, auto string) string {
	// 判断条件并进入对应分支逻辑。
	if v != nil {
		// 定义并初始化当前变量。
		manual := strings.TrimSpace(*v)
		// 判断条件并进入对应分支逻辑。
		if manual != "" {
			// 返回当前处理结果。
			return manual
		}
	}
	// 判断条件并进入对应分支逻辑。
	if strings.TrimSpace(current) != "" {
		// 返回当前处理结果。
		return current
	}
	// 返回当前处理结果。
	return auto
}

// parseIntWithDefault 将字符串解析为整数，失败时回退默认值。
func parseIntWithDefault(raw string, def int) int {
	// 定义并初始化当前变量。
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	// 判断条件并进入对应分支逻辑。
	if err != nil {
		// 返回当前处理结果。
		return def
	}
	// 返回当前处理结果。
	return v
}
