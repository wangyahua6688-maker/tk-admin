-- 清理开奖结果表中历史遗留的“倒计时配置字段”。
-- 说明：
-- 1) 倒计时配置统一放在 tk_special_lottery.next_draw_at；
-- 2) tk_draw_record 只保存“已开奖结果”，不再承载倒计时配置；
-- 3) 本脚本按“列存在才删除”执行，兼容不同环境版本。

DROP PROCEDURE IF EXISTS tk_drop_col_if_exists;
DELIMITER $$
CREATE PROCEDURE tk_drop_col_if_exists(IN p_table VARCHAR(64), IN p_column VARCHAR(64))
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = p_table
      AND column_name = p_column
  ) THEN
    SET @ddl = CONCAT('ALTER TABLE `', p_table, '` DROP COLUMN `', p_column, '`');
    PREPARE stmt FROM @ddl;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
  END IF;
END $$
DELIMITER ;

CALL tk_drop_col_if_exists('tk_draw_record', 'next_draw_at');
CALL tk_drop_col_if_exists('tk_draw_record', 'countdown_sec');
CALL tk_drop_col_if_exists('tk_draw_record', 'countdown_seconds');
CALL tk_drop_col_if_exists('tk_draw_record', 'draw_countdown_sec');
CALL tk_drop_col_if_exists('tk_draw_record', 'draw_countdown_seconds');

DROP PROCEDURE IF EXISTS tk_drop_col_if_exists;
