-- 2026-03-05 后台业务库语义同步脚本（tk_ 前缀）
-- 说明：
-- 1) tk_lottery_info 作为图库内容表，允许 special_lottery_id=0；
-- 2) 开奖区主数据使用 tk_draw_record。

ALTER TABLE tk_lottery_info
  MODIFY COLUMN special_lottery_id BIGINT UNSIGNED NOT NULL
  COMMENT '所属彩种ID（关联tk_special_lottery.id，0表示不绑定彩种）';

ALTER TABLE tk_lottery_info
  COMMENT = '图库图纸内容与竞猜配置表（不承载开奖区历史主数据）';

-- 3) 开奖记录拆分出独立属相/五行标签字段，支持历史页“只看属相”与“属相/五行”模式切换。
ALTER TABLE tk_draw_record
  ADD COLUMN zodiac_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码对应属相标签（与号码一一对应，示例：羊,蛇,马）' AFTER draw_labels,
  ADD COLUMN wuxing_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码对应五行标签（与号码一一对应，示例：土,金,火）' AFTER zodiac_labels;
