package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go-admin/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// drawResultStats 开奖主表需要同步回写的派生字段。
// 说明：
// 1. 这些字段仍保留在 tk_draw_record，兼容现有前端和业务接口；
// 2. 更完整的玩法结果会额外落到分表。
type drawResultStats struct {
	// SpecialSingleDouble 为特码单双。
	SpecialSingleDouble string
	// SpecialBigSmall 为特码大小。
	SpecialBigSmall string
	// SumSingleDouble 为总和单双。
	SumSingleDouble string
	// SumBigSmall 为总和大小。
	SumBigSmall string
	// SpecialCode 为特码数字。
	SpecialCode string
	// NormalCode 为正码数字串。
	NormalCode string
	// ZhengDescriptions 为正码1-6的人类可读描述。
	ZhengDescriptions [6]string
}

// drawResultBundle 为一次 6+1 开奖号码编译出的整套结果。
// 说明：
// 1. bundle 同时服务于主表字段回写和分表 upsert；
// 2. 这样可以确保一套开奖号码只经过一套规则计算，不会出现主表和分表不一致。
type drawResultBundle struct {
	// DrawLabels 为“生肖/五行”组合标签。
	DrawLabels string
	// ZodiacLabels 为生肖标签串。
	ZodiacLabels string
	// WuxingLabels 为五行标签串。
	WuxingLabels string
	// Stats 为主表派生字段。
	Stats drawResultStats
	// SpecialResult 为特码玩法结果。
	SpecialResult models.WDrawResultSpecial
	// RegularResult 为正码玩法结果。
	RegularResult models.WDrawResultRegular
	// CountResult 为统计玩法结果。
	CountResult models.WDrawResultCount
	// ZodiacTailResult 为生肖/尾数玩法结果。
	ZodiacTailResult models.WDrawResultZodiacTail
	// ComboResult 为组合玩法结果。
	ComboResult models.WDrawResultCombo
}

// compiledNumberDetail 为单个号码的完整规则属性。
type compiledNumberDetail struct {
	// Number 为当前号码。
	Number int `json:"number"`
	// Position 为开奖顺序位置（1-7）。
	Position int `json:"position"`
	// ColorWave 为波色。
	ColorWave string `json:"color_wave"`
	// BigSmall 为大小。
	BigSmall string `json:"big_small"`
	// SingleDouble 为单双。
	SingleDouble string `json:"single_double"`
	// SumSingleDouble 为合数单双。
	SumSingleDouble string `json:"sum_single_double"`
	// TailBigSmall 为尾数大小。
	TailBigSmall string `json:"tail_big_small"`
	// Zodiac 为生肖。
	Zodiac string `json:"zodiac"`
	// Wuxing 为五行。
	Wuxing string `json:"wuxing"`
	// Beast 为家畜/野兽。
	Beast string `json:"beast"`
	// TailLabel 为尾数标签（如 1尾）。
	TailLabel string `json:"tail_label"`
}

// compileDrawResultBundle 根据当前开奖记录编译整套玩法结果。
func compileDrawResultBundle(item *models.WDrawRecord) (*drawResultBundle, error) {
	// 解析并校验 6+1 开奖号码。
	nums, err := parseDrawNumbersForResult(item.DrawResult)
	if err != nil {
		return nil, err
	}
	// 开奖区必须严格为 6+1。
	if len(nums) != 7 {
		return nil, fmt.Errorf("draw_result must contain 7 numbers")
	}

	// 逐个号码编译规则属性。
	details := make([]compiledNumberDetail, 0, len(nums))
	// 预分配开奖号码整数切片，后续构建 JSON 时可直接复用。
	allNumbers := make([]int, 0, len(nums))
	// 正码切片固定只取前 6 个。
	normalNumbers := make([]int, 0, 6)
	// 初始化总和。
	totalSum := 0
	// 初始化七码单双大小统计。
	oddCount, evenCount, bigCount, smallCount := 0, 0, 0, 0
	// 遍历全部开奖号码。
	for idx, num := range nums {
		// 编译单个号码的规则属性。
		detail := compileNumberDetail(num, idx+1)
		// 收集完整号码细节。
		details = append(details, detail)
		// 收集全部开奖号码。
		allNumbers = append(allNumbers, num)
		// 累加总和。
		totalSum += num
		// 七码玩法里 49 按单、大处理。
		if countAsOdd(num) {
			oddCount++
		} else {
			evenCount++
		}
		// 七码玩法里 49 按大处理。
		if countAsBig(num) {
			bigCount++
		} else {
			smallCount++
		}
		// 前 6 个位置属于正码。
		if idx < 6 {
			normalNumbers = append(normalNumbers, num)
		}
	}

	// 特码固定为第 7 个号码。
	special := details[6]
	// 输出“生肖/五行”组合标签。
	drawLabels := make([]string, 0, len(details))
	// 输出生肖标签。
	zodiacLabels := make([]string, 0, len(details))
	// 输出五行标签。
	wuxingLabels := make([]string, 0, len(details))
	// 顺序化去重集合。
	appearedZodiacSet := map[string]struct{}{}
	// 顺序化去重集合。
	appearedTailSet := map[string]struct{}{}
	// 顺序化去重集合。
	appearedWuxingSet := map[string]struct{}{}
	// 用于区分家畜和野兽命中集合。
	homeBeastHitSet := map[string]struct{}{}
	// 用于区分家畜和野兽命中集合。
	wildBeastHitSet := map[string]struct{}{}
	// 保存正码1-6的结构化结果 JSON。
	zhengJSON := [6]string{}
	// 保存正码1-6的人类可读描述。
	zhengDescriptions := [6]string{}
	// 遍历号码细节，构建标签和位置结果。
	for idx, detail := range details {
		// draw_labels 统一为“生肖/五行”。
		drawLabels = append(drawLabels, detail.Zodiac+"/"+detail.Wuxing)
		// 单独输出生肖标签。
		zodiacLabels = append(zodiacLabels, detail.Zodiac)
		// 单独输出五行标签。
		wuxingLabels = append(wuxingLabels, detail.Wuxing)
		// 记录已出现生肖。
		appearedZodiacSet[detail.Zodiac] = struct{}{}
		// 记录已出现尾数。
		appearedTailSet[detail.TailLabel] = struct{}{}
		// 记录已出现五行。
		appearedWuxingSet[detail.Wuxing] = struct{}{}
		// 按规则区分家畜和野兽集合。
		if detail.Beast == "家畜" {
			homeBeastHitSet[detail.Zodiac] = struct{}{}
		}
		// 按规则区分家畜和野兽集合。
		if detail.Beast == "野兽" {
			wildBeastHitSet[detail.Zodiac] = struct{}{}
		}
		// 只有前 6 个号码才写正码位置结果。
		if idx < 6 {
			// 保存结构化结果，便于后续做位置玩法结算。
			zhengJSON[idx] = jsonText(detail)
			// 保存传统字符串描述，兼容旧字段输出。
			zhengDescriptions[idx] = composeZhengDescription(detail)
		}
	}

	// 计算未命中的生肖集合。
	missedZodiacs := diffOrderedSet(models.DrawResultZodiacOrder, appearedZodiacSet)
	// 计算未命中的尾数集合。
	missedTails := diffOrderedSet(models.DrawResultTailOrder, appearedTailSet)
	// 计算已命中的生肖集合。
	appearedZodiacs := collectOrderedSet(models.DrawResultZodiacOrder, appearedZodiacSet)
	// 计算已命中的尾数集合。
	appearedTails := collectOrderedSet(models.DrawResultTailOrder, appearedTailSet)
	// 计算已命中的五行集合。
	appearedWuxings := collectOrderedSet(models.DrawResultWuxingOrder, appearedWuxingSet)
	// 计算已命中的家畜生肖集合。
	homeBeastZodiacs := collectOrderedSet(models.DrawResultZodiacOrder, homeBeastHitSet)
	// 计算已命中的野兽生肖集合。
	wildBeastZodiacs := collectOrderedSet(models.DrawResultZodiacOrder, wildBeastHitSet)

	// 总分大小按 >=175 为大，<=174 为小。
	totalBigSmall := "小"
	// 命中阈值时改为大。
	if totalSum >= 175 {
		totalBigSmall = "大"
	}
	// 总分单双按总和奇偶计算。
	totalSingleDouble := "双"
	// 奇数总和为单。
	if totalSum%2 == 1 {
		totalSingleDouble = "单"
	}

	// 特码半波（波色+大小）。
	halfWaveColorSize := buildHalfWaveColorSize(special)
	// 特码半波（波色+单双）。
	halfWaveColorParity := buildHalfWaveColorParity(special)

	// 先构建主表派生字段，后续直接回写到 tk_draw_record。
	stats := drawResultStats{
		// 处理当前语句逻辑。
		SpecialSingleDouble: special.SingleDouble,
		// 处理当前语句逻辑。
		SpecialBigSmall: special.BigSmall,
		// 处理当前语句逻辑。
		SumSingleDouble: totalSingleDouble,
		// 处理当前语句逻辑。
		SumBigSmall: totalBigSmall,
		// 处理当前语句逻辑。
		SpecialCode: strconv.Itoa(special.Number),
		// 处理当前语句逻辑。
		NormalCode: joinIntCSVForResult(normalNumbers),
		// 处理当前语句逻辑。
		ZhengDescriptions: zhengDescriptions,
	}

	// 构造特码玩法结果表。
	specialResult := models.WDrawResultSpecial{
		// 处理当前语句逻辑。
		SpecialLotteryID: item.SpecialLotteryID,
		// 处理当前语句逻辑。
		Issue: item.Issue,
		// 处理当前语句逻辑。
		Year: item.Year,
		// 处理当前语句逻辑。
		DrawAt: item.DrawAt,
		// 处理当前语句逻辑。
		SpecialNumber: special.Number,
		// 处理当前语句逻辑。
		SpecialColorWave: special.ColorWave,
		// 处理当前语句逻辑。
		SpecialBigSmall: special.BigSmall,
		// 处理当前语句逻辑。
		SpecialSingleDouble: special.SingleDouble,
		// 处理当前语句逻辑。
		SpecialSumSingleDouble: special.SumSingleDouble,
		// 处理当前语句逻辑。
		SpecialTailBigSmall: special.TailBigSmall,
		// 处理当前语句逻辑。
		SpecialZodiac: special.Zodiac,
		// 处理当前语句逻辑。
		SpecialWuxing: special.Wuxing,
		// 处理当前语句逻辑。
		SpecialHomeBeast: special.Beast,
		// 处理当前语句逻辑。
		HalfWaveColorSize: halfWaveColorSize,
		// 处理当前语句逻辑。
		HalfWaveColorParity: halfWaveColorParity,
		// 处理当前语句逻辑。
		PayloadJSON: jsonText(map[string]interface{}{
			"special":                special,
			"two_sides":              []string{special.BigSmall, special.SingleDouble},
			"half_wave_color_size":   halfWaveColorSize,
			"half_wave_color_parity": halfWaveColorParity,
		}),
	}

	// 构造正码玩法结果表。
	regularResult := models.WDrawResultRegular{
		// 处理当前语句逻辑。
		SpecialLotteryID: item.SpecialLotteryID,
		// 处理当前语句逻辑。
		Issue: item.Issue,
		// 处理当前语句逻辑。
		Year: item.Year,
		// 处理当前语句逻辑。
		DrawAt: item.DrawAt,
		// 处理当前语句逻辑。
		NormalNumbers: joinIntCSVForResult(normalNumbers),
		// 处理当前语句逻辑。
		TotalSum: totalSum,
		// 处理当前语句逻辑。
		TotalBigSmall: totalBigSmall,
		// 处理当前语句逻辑。
		TotalSingleDouble: totalSingleDouble,
		// 处理当前语句逻辑。
		Zheng1JSON: zhengJSON[0],
		// 处理当前语句逻辑。
		Zheng2JSON: zhengJSON[1],
		// 处理当前语句逻辑。
		Zheng3JSON: zhengJSON[2],
		// 处理当前语句逻辑。
		Zheng4JSON: zhengJSON[3],
		// 处理当前语句逻辑。
		Zheng5JSON: zhengJSON[4],
		// 处理当前语句逻辑。
		Zheng6JSON: zhengJSON[5],
		// 处理当前语句逻辑。
		PayloadJSON: jsonText(map[string]interface{}{
			"normal_numbers":      normalNumbers,
			"total_sum":           totalSum,
			"total_big_small":     totalBigSmall,
			"total_single_double": totalSingleDouble,
			"positions":           details[:6],
		}),
	}

	// 构造统计玩法结果表。
	countResult := models.WDrawResultCount{
		// 处理当前语句逻辑。
		SpecialLotteryID: item.SpecialLotteryID,
		// 处理当前语句逻辑。
		Issue: item.Issue,
		// 处理当前语句逻辑。
		Year: item.Year,
		// 处理当前语句逻辑。
		DrawAt: item.DrawAt,
		// 处理当前语句逻辑。
		TotalSum: totalSum,
		// 处理当前语句逻辑。
		OddCount: oddCount,
		// 处理当前语句逻辑。
		EvenCount: evenCount,
		// 处理当前语句逻辑。
		BigCount: bigCount,
		// 处理当前语句逻辑。
		SmallCount: smallCount,
		// 处理当前语句逻辑。
		DistinctZodiacCount: len(appearedZodiacs),
		// 处理当前语句逻辑。
		DistinctTailCount: len(appearedTails),
		// 处理当前语句逻辑。
		DistinctWuxingCount: len(appearedWuxings),
		// 处理当前语句逻辑。
		AppearedZodiacs: strings.Join(appearedZodiacs, ","),
		// 处理当前语句逻辑。
		MissedZodiacs: strings.Join(missedZodiacs, ","),
		// 处理当前语句逻辑。
		AppearedTails: strings.Join(appearedTails, ","),
		// 处理当前语句逻辑。
		MissedTails: strings.Join(missedTails, ","),
		// 处理当前语句逻辑。
		AppearedWuxings: strings.Join(appearedWuxings, ","),
		// 处理当前语句逻辑。
		PayloadJSON: jsonText(map[string]interface{}{
			"qi_ma": map[string]int{
				"单": oddCount,
				"双": evenCount,
				"大": bigCount,
				"小": smallCount,
			},
			"yixiao_liang":     len(appearedZodiacs),
			"weishu_liang":     len(appearedTails),
			"wuxing_liang":     len(appearedWuxings),
			"appeared_zodiacs": appearedZodiacs,
			"appeared_tails":   appearedTails,
			"appeared_wuxings": appearedWuxings,
		}),
	}

	// 构造生肖/尾数玩法结果表。
	zodiacTailResult := models.WDrawResultZodiacTail{
		// 处理当前语句逻辑。
		SpecialLotteryID: item.SpecialLotteryID,
		// 处理当前语句逻辑。
		Issue: item.Issue,
		// 处理当前语句逻辑。
		Year: item.Year,
		// 处理当前语句逻辑。
		DrawAt: item.DrawAt,
		// 处理当前语句逻辑。
		SpecialZodiac: special.Zodiac,
		// 处理当前语句逻辑。
		SpecialHomeBeast: special.Beast,
		// 处理当前语句逻辑。
		SpecialWuxing: special.Wuxing,
		// 处理当前语句逻辑。
		HitZodiacs: strings.Join(appearedZodiacs, ","),
		// 处理当前语句逻辑。
		MissZodiacs: strings.Join(missedZodiacs, ","),
		// 处理当前语句逻辑。
		HitTails: strings.Join(appearedTails, ","),
		// 处理当前语句逻辑。
		MissTails: strings.Join(missedTails, ","),
		// 处理当前语句逻辑。
		HomeBeastZodiacs: strings.Join(homeBeastZodiacs, ","),
		// 处理当前语句逻辑。
		WildBeastZodiacs: strings.Join(wildBeastZodiacs, ","),
		// 处理当前语句逻辑。
		PayloadJSON: jsonText(map[string]interface{}{
			"special_zodiac":     special.Zodiac,
			"special_beast":      special.Beast,
			"hit_zodiacs":        appearedZodiacs,
			"miss_zodiacs":       missedZodiacs,
			"hit_tails":          appearedTails,
			"miss_tails":         missedTails,
			"home_beast_zodiacs": homeBeastZodiacs,
			"wild_beast_zodiacs": wildBeastZodiacs,
		}),
	}

	// 构造组合玩法结果表。
	comboResult := models.WDrawResultCombo{
		// 处理当前语句逻辑。
		SpecialLotteryID: item.SpecialLotteryID,
		// 处理当前语句逻辑。
		Issue: item.Issue,
		// 处理当前语句逻辑。
		Year: item.Year,
		// 处理当前语句逻辑。
		DrawAt: item.DrawAt,
		// 处理当前语句逻辑。
		NormalNumbers: joinIntCSVForResult(normalNumbers),
		// 处理当前语句逻辑。
		AllNumbers: joinIntCSVForResult(allNumbers),
		// 处理当前语句逻辑。
		SpecialNumber: special.Number,
		// 处理当前语句逻辑。
		PayloadJSON: jsonText(map[string]interface{}{
			"lianma": map[string]interface{}{
				"normal_numbers": normalNumbers,
			},
			"guoguan": map[string]interface{}{
				"all_numbers":     allNumbers,
				"special_rule_49": "49 为和局时大小单双按 1 计，波色按绿波处理",
			},
			"erzhongte": map[string]interface{}{
				"normal_numbers": normalNumbers,
				"special_number": special.Number,
			},
			"techuan": map[string]interface{}{
				"normal_numbers": normalNumbers,
				"special_number": special.Number,
			},
			"buzhong": map[string]interface{}{
				"all_numbers": allNumbers,
			},
			"duoxuan_zhongyi": map[string]interface{}{
				"all_numbers": allNumbers,
			},
			"tepingzhong": map[string]interface{}{
				"all_numbers": allNumbers,
			},
		}),
	}

	// 返回完整 bundle。
	return &drawResultBundle{
		// 处理当前语句逻辑。
		DrawLabels: strings.Join(drawLabels, ","),
		// 处理当前语句逻辑。
		ZodiacLabels: strings.Join(zodiacLabels, ","),
		// 处理当前语句逻辑。
		WuxingLabels: strings.Join(wuxingLabels, ","),
		// 处理当前语句逻辑。
		Stats: stats,
		// 处理当前语句逻辑。
		SpecialResult: specialResult,
		// 处理当前语句逻辑。
		RegularResult: regularResult,
		// 处理当前语句逻辑。
		CountResult: countResult,
		// 处理当前语句逻辑。
		ZodiacTailResult: zodiacTailResult,
		// 处理当前语句逻辑。
		ComboResult: comboResult,
	}, nil
}

// hydrateDrawRecordDerivedFields 将规则结果回写到开奖记录主表兼容字段。
func hydrateDrawRecordDerivedFields(item *models.WDrawRecord) (*drawResultBundle, error) {
	// 编译整套开奖结果。
	bundle, err := compileDrawResultBundle(item)
	if err != nil {
		return nil, err
	}
	// 将自动生成的官方标签回写到主表。
	item.DrawLabels = bundle.DrawLabels
	// 将生肖标签回写到主表。
	item.ZodiacLabels = bundle.ZodiacLabels
	// 将五行标签回写到主表。
	item.WuxingLabels = bundle.WuxingLabels
	// 将基础玩法结果回写到主表兼容字段。
	item.SpecialSingleDouble = bundle.Stats.SpecialSingleDouble
	// 将基础玩法结果回写到主表兼容字段。
	item.SpecialBigSmall = bundle.Stats.SpecialBigSmall
	// 将基础玩法结果回写到主表兼容字段。
	item.SumSingleDouble = bundle.Stats.SumSingleDouble
	// 将基础玩法结果回写到主表兼容字段。
	item.SumBigSmall = bundle.Stats.SumBigSmall
	// 将基础玩法结果回写到主表兼容字段。
	item.SpecialCode = bundle.Stats.SpecialCode
	// 将基础玩法结果回写到主表兼容字段。
	item.NormalCode = bundle.Stats.NormalCode
	// 位置玩法描述也跟随官方规则自动生成。
	item.Zheng1 = bundle.Stats.ZhengDescriptions[0]
	// 位置玩法描述也跟随官方规则自动生成。
	item.Zheng2 = bundle.Stats.ZhengDescriptions[1]
	// 位置玩法描述也跟随官方规则自动生成。
	item.Zheng3 = bundle.Stats.ZhengDescriptions[2]
	// 位置玩法描述也跟随官方规则自动生成。
	item.Zheng4 = bundle.Stats.ZhengDescriptions[3]
	// 位置玩法描述也跟随官方规则自动生成。
	item.Zheng5 = bundle.Stats.ZhengDescriptions[4]
	// 位置玩法描述也跟随官方规则自动生成。
	item.Zheng6 = bundle.Stats.ZhengDescriptions[5]
	// 返回 bundle 供事务内 upsert 分表复用。
	return bundle, nil
}

// upsertDrawResultTablesTx 在事务内写入全部玩法结果分表。
func upsertDrawResultTablesTx(tx *gorm.DB, drawRecordID uint, bundle *drawResultBundle) error {
	// 将主键关联补到特码玩法结果。
	bundle.SpecialResult.DrawRecordID = drawRecordID
	// 将主键关联补到正码玩法结果。
	bundle.RegularResult.DrawRecordID = drawRecordID
	// 将主键关联补到统计玩法结果。
	bundle.CountResult.DrawRecordID = drawRecordID
	// 将主键关联补到生肖/尾数玩法结果。
	bundle.ZodiacTailResult.DrawRecordID = drawRecordID
	// 将主键关联补到组合玩法结果。
	bundle.ComboResult.DrawRecordID = drawRecordID

	// 依次 upsert 特码玩法结果。
	if err := upsertByDrawRecordID(tx, &bundle.SpecialResult, []string{
		"special_lottery_id", "issue", "year", "draw_at", "special_number", "special_color_wave",
		"special_big_small", "special_single_double", "special_sum_single_double", "special_tail_big_small",
		"special_zodiac", "special_wuxing", "special_home_beast", "half_wave_color_size",
		"half_wave_color_parity", "payload_json", "updated_at",
	}); err != nil {
		return err
	}
	// 依次 upsert 正码玩法结果。
	if err := upsertByDrawRecordID(tx, &bundle.RegularResult, []string{
		"special_lottery_id", "issue", "year", "draw_at", "normal_numbers", "total_sum",
		"total_big_small", "total_single_double", "zheng1_json", "zheng2_json", "zheng3_json",
		"zheng4_json", "zheng5_json", "zheng6_json", "payload_json", "updated_at",
	}); err != nil {
		return err
	}
	// 依次 upsert 统计玩法结果。
	if err := upsertByDrawRecordID(tx, &bundle.CountResult, []string{
		"special_lottery_id", "issue", "year", "draw_at", "total_sum", "odd_count", "even_count",
		"big_count", "small_count", "distinct_zodiac_count", "distinct_tail_count", "distinct_wuxing_count",
		"appeared_zodiacs", "missed_zodiacs", "appeared_tails", "missed_tails", "appeared_wuxings",
		"payload_json", "updated_at",
	}); err != nil {
		return err
	}
	// 依次 upsert 生肖/尾数玩法结果。
	if err := upsertByDrawRecordID(tx, &bundle.ZodiacTailResult, []string{
		"special_lottery_id", "issue", "year", "draw_at", "special_zodiac", "special_home_beast",
		"special_wuxing", "hit_zodiacs", "miss_zodiacs", "hit_tails", "miss_tails",
		"home_beast_zodiacs", "wild_beast_zodiacs", "payload_json", "updated_at",
	}); err != nil {
		return err
	}
	// 最后 upsert 组合玩法结果。
	return upsertByDrawRecordID(tx, &bundle.ComboResult, []string{
		"special_lottery_id", "issue", "year", "draw_at", "normal_numbers", "all_numbers",
		"special_number", "payload_json", "updated_at",
	})
}

// deleteDrawResultTablesTx 删除一条开奖记录关联的全部玩法结果。
func deleteDrawResultTablesTx(tx *gorm.DB, drawRecordID uint) error {
	// 逐表删除，保证分表不会残留孤儿数据。
	if err := tx.Where("draw_record_id = ?", drawRecordID).Delete(&models.WDrawResultSpecial{}).Error; err != nil {
		return err
	}
	// 删除正码玩法结果。
	if err := tx.Where("draw_record_id = ?", drawRecordID).Delete(&models.WDrawResultRegular{}).Error; err != nil {
		return err
	}
	// 删除统计玩法结果。
	if err := tx.Where("draw_record_id = ?", drawRecordID).Delete(&models.WDrawResultCount{}).Error; err != nil {
		return err
	}
	// 删除生肖/尾数玩法结果。
	if err := tx.Where("draw_record_id = ?", drawRecordID).Delete(&models.WDrawResultZodiacTail{}).Error; err != nil {
		return err
	}
	// 删除组合玩法结果。
	return tx.Where("draw_record_id = ?", drawRecordID).Delete(&models.WDrawResultCombo{}).Error
}

// upsertByDrawRecordID 以 draw_record_id 为唯一键做幂等写入。
func upsertByDrawRecordID(tx *gorm.DB, model interface{}, updateColumns []string) error {
	// 统一用 draw_record_id 做冲突更新，保证“一期开奖结果一条分表记录”。
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "draw_record_id"}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(model).Error
}

// parseDrawNumbersForResult 解析 6+1 开奖号码。
func parseDrawNumbersForResult(raw string) ([]int, error) {
	// 统一清洗各种分隔符。
	normalized := strings.NewReplacer(
		"，", ",",
		"、", ",",
		";", ",",
		"|", ",",
		" ", ",",
		"\n", ",",
		"\r", ",",
		"\t", ",",
	).Replace(strings.TrimSpace(raw))
	// 拆分出号码片段。
	parts := strings.Split(normalized, ",")
	// 为号码结果预留 7 个位置。
	out := make([]int, 0, 7)
	// 使用 seen 防止重复号码。
	seen := map[int]struct{}{}
	// 逐个解析号码。
	for _, part := range parts {
		// 去掉空白片段。
		token := strings.TrimSpace(part)
		// 空片段直接跳过。
		if token == "" {
			continue
		}
		// 将文本号码转成整数。
		num, err := strconv.Atoi(token)
		// 非数字直接报错。
		if err != nil {
			return nil, fmt.Errorf("invalid draw number: %s", token)
		}
		// 六合彩号码范围固定为 1-49。
		if num < 1 || num > 49 {
			return nil, fmt.Errorf("draw number out of range: %d", num)
		}
		// 开奖号不可重复。
		if _, exists := seen[num]; exists {
			return nil, fmt.Errorf("duplicate draw number: %d", num)
		}
		// 记录已出现号码。
		seen[num] = struct{}{}
		// 收集解析后的号码。
		out = append(out, num)
	}
	// 返回解析结果。
	return out, nil
}

// compileNumberDetail 根据号码生成规则属性。
func compileNumberDetail(num, position int) compiledNumberDetail {
	// 构建当前号码的完整属性。
	return compiledNumberDetail{
		// 处理当前语句逻辑。
		Number: num,
		// 处理当前语句逻辑。
		Position: position,
		// 处理当前语句逻辑。
		ColorWave: models.DrawResultColorWaveMap[num],
		// 处理当前语句逻辑。
		BigSmall: specialBigSmall(num),
		// 处理当前语句逻辑。
		SingleDouble: specialSingleDouble(num),
		// 处理当前语句逻辑。
		SumSingleDouble: specialSumSingleDouble(num),
		// 处理当前语句逻辑。
		TailBigSmall: specialTailBigSmall(num),
		// 处理当前语句逻辑。
		Zodiac: models.DrawResultZodiacMap[num],
		// 处理当前语句逻辑。
		Wuxing: models.DrawResultWuxingMap[num],
		// 处理当前语句逻辑。
		Beast: models.DrawResultBeastMap[models.DrawResultZodiacMap[num]],
		// 处理当前语句逻辑。
		TailLabel: fmt.Sprintf("%d尾", num%10),
	}
}

// specialBigSmall 返回特码/正码位用的大小结果。
func specialBigSmall(num int) string {
	// 49 为和局。
	if num == 49 {
		return "和"
	}
	// 25-48 为大。
	if num >= 25 {
		return "大"
	}
	// 1-24 为小。
	return "小"
}

// specialSingleDouble 返回特码/正码位用的单双结果。
func specialSingleDouble(num int) string {
	// 49 为和局。
	if num == 49 {
		return "和"
	}
	// 偶数为双。
	if num%2 == 0 {
		return "双"
	}
	// 奇数为单。
	return "单"
}

// specialSumSingleDouble 返回特码/正码位用的合数单双结果。
func specialSumSingleDouble(num int) string {
	// 49 为和局。
	if num == 49 {
		return "和"
	}
	// 计算十位和个位之和。
	sum := num/10 + num%10
	// 和数为偶数则合双。
	if sum%2 == 0 {
		return "合双"
	}
	// 和数为奇数则合单。
	return "合单"
}

// specialTailBigSmall 返回特码尾数大小结果。
func specialTailBigSmall(num int) string {
	// 49 为和局。
	if num == 49 {
		return "和"
	}
	// 尾数 0-4 为尾小。
	if num%10 <= 4 {
		return "尾小"
	}
	// 尾数 5-9 为尾大。
	return "尾大"
}

// countAsOdd 七码统计里判断是否记为单。
func countAsOdd(num int) bool {
	// 49 按单计算。
	if num == 49 {
		return true
	}
	// 其它号码按正常奇偶判断。
	return num%2 == 1
}

// countAsBig 七码统计里判断是否记为大。
func countAsBig(num int) bool {
	// 49 按大计算。
	if num == 49 {
		return true
	}
	// 25 及以上按大计算。
	return num >= 25
}

// buildHalfWaveColorSize 组装半波（波色+大小）结果。
func buildHalfWaveColorSize(detail compiledNumberDetail) string {
	// 49 对半波统一视为和局。
	if detail.Number == 49 {
		return "和局"
	}
	// 去掉“波”后拼接大小单双。
	return strings.TrimSuffix(detail.ColorWave, "波") + detail.BigSmall
}

// buildHalfWaveColorParity 组装半波（波色+单双）结果。
func buildHalfWaveColorParity(detail compiledNumberDetail) string {
	// 49 对半波统一视为和局。
	if detail.Number == 49 {
		return "和局"
	}
	// 去掉“波”后拼接大小单双。
	return strings.TrimSuffix(detail.ColorWave, "波") + detail.SingleDouble
}

// composeZhengDescription 组装正码位置的描述串。
func composeZhengDescription(detail compiledNumberDetail) string {
	// 按固定顺序输出，便于前端直接展示。
	parts := []string{
		detail.BigSmall,
		detail.SingleDouble,
		detail.ColorWave,
		detail.SumSingleDouble,
		detail.TailBigSmall,
		detail.Zodiac,
		detail.Wuxing,
	}
	// 返回逗号分隔描述。
	return strings.Join(parts, ",")
}

// collectOrderedSet 按给定顺序收集命中集合。
func collectOrderedSet(order []string, set map[string]struct{}) []string {
	// 预分配结果切片。
	out := make([]string, 0, len(set))
	// 按固定顺序输出。
	for _, label := range order {
		// 命中则输出。
		if _, ok := set[label]; ok {
			out = append(out, label)
		}
	}
	// 返回结果集合。
	return out
}

// diffOrderedSet 按给定顺序收集未命中集合。
func diffOrderedSet(order []string, set map[string]struct{}) []string {
	// 预分配结果切片。
	out := make([]string, 0, len(order))
	// 按固定顺序扫描全集。
	for _, label := range order {
		// 未命中时输出。
		if _, ok := set[label]; !ok {
			out = append(out, label)
		}
	}
	// 返回结果集合。
	return out
}

// jsonText 将任意结构稳定编码为 JSON 字符串。
func jsonText(v interface{}) string {
	// 执行 JSON 编码。
	payload, err := json.Marshal(v)
	// 编码失败时回退空对象，避免写库报错。
	if err != nil {
		return "{}"
	}
	// 返回 JSON 文本。
	return string(payload)
}

// joinIntCSVForResult 将号码列表转成逗号分隔字符串。
func joinIntCSVForResult(nums []int) string {
	// 预分配字符串切片。
	parts := make([]string, 0, len(nums))
	// 逐个转成字符串。
	for _, num := range nums {
		parts = append(parts, strconv.Itoa(num))
	}
	// 返回 CSV 文本。
	return strings.Join(parts, ",")
}

// buildDrawRecordUpdateMap 将完整开奖记录对象转成可写入数据库的 map。
func buildDrawRecordUpdateMap(item *models.WDrawRecord) map[string]interface{} {
	// 统一输出控制字段，避免 update 路径遗漏新衍生字段。
	return map[string]interface{}{
		"special_lottery_id":    item.SpecialLotteryID,
		"issue":                 item.Issue,
		"year":                  item.Year,
		"draw_at":               item.DrawAt,
		"normal_draw_result":    item.NormalDrawResult,
		"special_draw_result":   item.SpecialDrawResult,
		"draw_result":           item.DrawResult,
		"draw_labels":           item.DrawLabels,
		"zodiac_labels":         item.ZodiacLabels,
		"wuxing_labels":         item.WuxingLabels,
		"playback_url":          item.PlaybackURL,
		"special_single_double": item.SpecialSingleDouble,
		"special_big_small":     item.SpecialBigSmall,
		"sum_single_double":     item.SumSingleDouble,
		"sum_big_small":         item.SumBigSmall,
		"recommend_six":         item.RecommendSix,
		"recommend_four":        item.RecommendFour,
		"recommend_one":         item.RecommendOne,
		"recommend_ten":         item.RecommendTen,
		"special_code":          item.SpecialCode,
		"normal_code":           item.NormalCode,
		"zheng1":                item.Zheng1,
		"zheng2":                item.Zheng2,
		"zheng3":                item.Zheng3,
		"zheng4":                item.Zheng4,
		"zheng5":                item.Zheng5,
		"zheng6":                item.Zheng6,
		"is_current":            item.IsCurrent,
		"status":                item.Status,
		"sort":                  item.Sort,
	}
}

// CreateDrawRecord 创建开奖记录并同步玩法结果分表。
func (s *BizConfigService) CreateDrawRecord(ctx context.Context, item *models.WDrawRecord) error {
	// 先统一回写主表派生字段，确保主表和分表来自同一套规则。
	bundle, err := hydrateDrawRecordDerivedFields(item)
	if err != nil {
		return err
	}
	// 通过事务保证主表与分表原子一致。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 同彩种只允许一条当前期记录。
		if item.IsCurrent == 1 {
			// 先清除同彩种其它记录的 current 标识。
			if err := s.dao.ResetCurrentDrawRecordTx(tx, item.SpecialLotteryID, 0); err != nil {
				return err
			}
		}
		// 写入开奖记录主表。
		if err := s.dao.CreateDrawRecordTx(tx, item); err != nil {
			return err
		}
		// 主表写入成功后，同步写入全部玩法结果分表。
		return upsertDrawResultTablesTx(tx, item.ID, bundle)
	})
}

// UpdateDrawRecord 更新开奖记录并同步玩法结果分表。
func (s *BizConfigService) UpdateDrawRecord(ctx context.Context, id uint, item *models.WDrawRecord) error {
	// 先统一回写主表派生字段，避免更新时只改号码不改玩法结果。
	bundle, err := hydrateDrawRecordDerivedFields(item)
	if err != nil {
		return err
	}
	// 通过事务保证主表和分表原子一致。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// current 置为 1 时，先清理同彩种其它记录。
		if item.IsCurrent == 1 && item.SpecialLotteryID > 0 {
			if err := s.dao.ResetCurrentDrawRecordTx(tx, item.SpecialLotteryID, id); err != nil {
				return err
			}
		}
		// 更新开奖记录主表。
		if err := s.dao.UpdateDrawRecordTx(tx, id, buildDrawRecordUpdateMap(item)); err != nil {
			return err
		}
		// 主表更新成功后，同步更新全部玩法结果分表。
		return upsertDrawResultTablesTx(tx, id, bundle)
	})
}

// DeleteDrawRecord 删除开奖记录及其玩法结果分表。
func (s *BizConfigService) DeleteDrawRecord(ctx context.Context, id uint) error {
	// 通过事务保证主表和分表一起删除。
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删玩法结果分表，避免残留孤儿数据。
		if err := deleteDrawResultTablesTx(tx, id); err != nil {
			return err
		}
		// 最后删除开奖记录主表。
		return tx.Delete(&models.WDrawRecord{}, id).Error
	})
}
