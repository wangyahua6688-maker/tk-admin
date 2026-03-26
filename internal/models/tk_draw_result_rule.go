package models

import common_model "github.com/wangyahua6688-maker/tk-common/models"

// 统一复用 tk-common 中维护的开奖结果规则常量，避免各业务服务各自维护一套映射。
var (
	// DrawResultZodiacOrder 固定生肖顺序，用于稳定输出。
	DrawResultZodiacOrder = common_model.DrawResultZodiacOrder
	// DrawResultWuxingOrder 固定五行顺序，用于稳定输出。
	DrawResultWuxingOrder = common_model.DrawResultWuxingOrder
	// DrawResultTailOrder 固定尾数顺序，用于稳定输出。
	DrawResultTailOrder = common_model.DrawResultTailOrder
	// DrawResultColorWaveMap 号码到波色的映射。
	DrawResultColorWaveMap = common_model.DrawResultColorWaveMap
	// DrawResultZodiacMap 号码到生肖的映射。
	DrawResultZodiacMap = common_model.DrawResultZodiacMap
	// DrawResultWuxingMap 号码到五行的映射。
	DrawResultWuxingMap = common_model.DrawResultWuxingMap
	// DrawResultBeastMap 号码到家畜/野兽分类的映射。
	DrawResultBeastMap = common_model.DrawResultBeastMap
)
