-- TK 平台业务初始化数据（tk_ 前缀）
-- 说明：
-- 1) 本脚本仅包含业务数据，不包含 RBAC 系统数据；
-- 2) 可重复执行（幂等）；
-- 3) 包含可选的 w_* -> tk_* 数据迁移（自动判断源表是否存在）。

SET NAMES utf8mb4;

-- =========================================================
-- A. 基础初始化数据（幂等）
-- =========================================================

-- 1) 默认短信通道（开发/联调可用）。
INSERT INTO tk_sms_channel (
  id, provider, channel_name, access_key, access_secret, endpoint, sign_name,
  template_code_login, template_code_register,
  daily_limit, minute_limit, code_ttl_seconds, mock_mode, status, created_at, updated_at
) VALUES
  (1, 'mock', '默认模拟通道', '', '', '', 'TK平台', '', '', 20, 1, 300, 1, 1, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE
  provider = VALUES(provider),
  channel_name = VALUES(channel_name),
  sign_name = VALUES(sign_name),
  daily_limit = VALUES(daily_limit),
  minute_limit = VALUES(minute_limit),
  code_ttl_seconds = VALUES(code_ttl_seconds),
  mock_mode = VALUES(mock_mode),
  status = VALUES(status),
  updated_at = VALUES(updated_at);

-- 2) 默认彩种（澳彩/港彩）。
INSERT INTO tk_special_lottery (
  name, code, current_issue, next_draw_at, live_enabled, live_status, live_stream_url, status, sort, created_at, updated_at
) VALUES
  ('澳彩', 'macau', '', DATE_ADD(NOW(3), INTERVAL 1 DAY), 0, 'pending', '', 1, 1, NOW(3), NOW(3)),
  ('港彩', 'hk', '', DATE_ADD(NOW(3), INTERVAL 1 DAY), 0, 'pending', '', 1, 2, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  sort = VALUES(sort),
  status = VALUES(status),
  updated_at = VALUES(updated_at);

-- 3) 默认首页分类。
INSERT INTO tk_lottery_category (
  category_key, name, search_keywords, show_on_home, status, sort, created_at, updated_at
) VALUES
  ('jiuxiao', '九肖系列', '九肖 连肖 稳肖', 1, 1, 1, NOW(3), NOW(3)),
  ('neimu', '内幕系列', '内幕 爆料 速报', 1, 1, 2, NOW(3), NOW(3)),
  ('sibuxiang', '四不像系列', '四不像 图库', 1, 1, 3, NOW(3), NOW(3)),
  ('paolutu', '跑狗图系列', '跑狗图 图库 轨迹', 1, 1, 4, NOW(3), NOW(3)),
  ('guapai', '挂牌系列', '挂牌 波色', 1, 1, 5, NOW(3), NOW(3)),
  ('more', '更多', '更多 分类', 1, 1, 6, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  search_keywords = VALUES(search_keywords),
  show_on_home = VALUES(show_on_home),
  status = VALUES(status),
  sort = VALUES(sort),
  updated_at = VALUES(updated_at);

-- 4) 默认金刚导航（若不存在则插入）。
INSERT INTO tk_external_link (
  name, url, position, icon_url, group_key, status, sort, created_at, updated_at
)
SELECT
  '开奖现场', '/home/live-scene', 'home_kingkong', '', 'live_scene', 1, 1, NOW(3), NOW(3)
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1
  FROM tk_external_link
  WHERE position = 'home_kingkong' AND (name = '开奖现场' OR group_key = 'live_scene')
);

-- 5) 首页主题背景与左右浮动广告（若不存在则插入）。
INSERT INTO tk_external_link (
  name, url, position, icon_url, group_key, status, sort, created_at, updated_at
)
SELECT
  '首页主题背景', '/forum', 'home_theme_bg',
  'https://jmz.jlidesign.com:4949/unite49files/amyd/2026/02/15/20260215221804-115979298.jpg',
  'home_theme', 1, 1, NOW(3), NOW(3)
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1
  FROM tk_external_link
  WHERE position = 'home_theme_bg'
);

INSERT INTO tk_external_link (
  name, url, position, icon_url, group_key, status, sort, created_at, updated_at
)
SELECT
  '左侧活动广告', '/forum', 'home_float_left',
  'https://tk.jlidesign.com:4949/m/col/25/xbpgb.jpg',
  'home_float', 1, 1, NOW(3), NOW(3)
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1
  FROM tk_external_link
  WHERE position = 'home_float_left'
);

INSERT INTO tk_external_link (
  name, url, position, icon_url, group_key, status, sort, created_at, updated_at
)
SELECT
  '右侧活动广告', '/lottery/categories', 'home_float_right',
  'https://amo.jlidesign.com:4949/m/col/63/ampgt.jpg',
  'home_float', 1, 1, NOW(3), NOW(3)
FROM DUAL
WHERE NOT EXISTS (
  SELECT 1
  FROM tk_external_link
  WHERE position = 'home_float_right'
);

-- =========================================================
-- B. 可选迁移：w_* -> tk_*（自动检测源表）
-- =========================================================

DELIMITER $$

DROP PROCEDURE IF EXISTS migrate_w_to_tk_if_exists $$

CREATE PROCEDURE migrate_w_to_tk_if_exists()
BEGIN
  -- ---------- 外链 ----------
  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = DATABASE() AND table_name = 'w_external_link'
  ) THEN
    UPDATE tk_external_link t
    JOIN w_external_link w
      ON BINARY t.name = BINARY w.name
     AND BINARY t.url = BINARY w.url
     AND BINARY t.position = BINARY w.position
    SET
      t.icon_url = IFNULL(w.icon_url, ''),
      t.group_key = IFNULL(w.group_key, ''),
      t.status = w.status,
      t.sort = w.sort,
      t.updated_at = IFNULL(w.updated_at, t.updated_at);

    INSERT INTO tk_external_link (
      name, url, position, icon_url, group_key, status, sort, created_at, updated_at
    )
    SELECT
      w.name, w.url, w.position, IFNULL(w.icon_url, ''), IFNULL(w.group_key, ''), w.status, w.sort,
      IFNULL(w.created_at, NOW(3)), IFNULL(w.updated_at, NOW(3))
    FROM w_external_link w
    LEFT JOIN tk_external_link t
      ON BINARY t.name = BINARY w.name
     AND BINARY t.url = BINARY w.url
     AND BINARY t.position = BINARY w.position
    WHERE t.id IS NULL;
  END IF;

  -- ---------- 分类 ----------
  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = DATABASE() AND table_name = 'w_lottery_category'
  ) THEN
    UPDATE tk_lottery_category t
    JOIN w_lottery_category w
      ON BINARY t.category_key = BINARY w.category_key
    SET
      t.name = w.name,
      t.search_keywords = IFNULL(w.search_keywords, ''),
      t.show_on_home = w.show_on_home,
      t.status = w.status,
      t.sort = w.sort,
      t.updated_at = IFNULL(w.updated_at, t.updated_at);

    INSERT INTO tk_lottery_category (
      category_key, name, search_keywords, show_on_home, status, sort, created_at, updated_at
    )
    SELECT
      w.category_key, w.name, IFNULL(w.search_keywords, ''), w.show_on_home, w.status, w.sort,
      IFNULL(w.created_at, NOW(3)), IFNULL(w.updated_at, NOW(3))
    FROM w_lottery_category w
    LEFT JOIN tk_lottery_category t
      ON BINARY t.category_key = BINARY w.category_key
    WHERE t.id IS NULL;
  END IF;

  -- ---------- 图纸 ----------
  IF EXISTS (
    SELECT 1 FROM information_schema.tables
    WHERE table_schema = DATABASE() AND table_name = 'w_lottery_info'
  ) THEN
    UPDATE tk_lottery_info t
    JOIN w_lottery_info w
      ON t.special_lottery_id = w.special_lottery_id
     AND BINARY t.issue = BINARY w.issue
     AND BINARY t.title = BINARY w.title
    SET
      t.category_id = COALESCE((
        SELECT c.id FROM tk_lottery_category c
        WHERE c.category_key COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
           OR c.name COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
        ORDER BY c.id ASC LIMIT 1
      ), t.category_id),
      t.category_tag = IFNULL(w.category_tag, ''),
      t.year = w.year,
      t.cover_image_url = w.cover_image_url,
      t.detail_image_url = w.detail_image_url,
      t.draw_code = w.draw_code,
      t.normal_draw_result = TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', 6)),
      t.special_draw_result = TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', -1)),
      t.draw_result = w.draw_result,
      t.draw_at = w.draw_at,
      t.playback_url = IFNULL(t.playback_url, ''),
      t.is_current = w.is_current,
      t.status = w.status,
      t.sort = w.sort,
      t.likes_count = w.likes_count,
      t.comment_count = w.comment_count,
      t.favorite_count = w.favorite_count,
      t.read_count = w.read_count,
      t.poll_enabled = w.poll_enabled,
      t.poll_default_expand = w.poll_default_expand,
      t.recommend_info_ids = IFNULL(w.recommend_info_ids, ''),
      t.updated_at = IFNULL(w.updated_at, t.updated_at);

    INSERT INTO tk_lottery_info (
      special_lottery_id, category_id, category_tag, issue, year, title,
      cover_image_url, detail_image_url, draw_code, normal_draw_result, special_draw_result, draw_result, draw_at, playback_url,
      is_current, status, sort, likes_count, comment_count, favorite_count, read_count,
      poll_enabled, poll_default_expand, recommend_info_ids, created_at, updated_at
    )
    SELECT
      w.special_lottery_id,
      COALESCE((
        SELECT c.id FROM tk_lottery_category c
        WHERE c.category_key COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
           OR c.name COLLATE utf8mb4_general_ci = IFNULL(w.category_tag, '') COLLATE utf8mb4_general_ci
        ORDER BY c.id ASC LIMIT 1
      ), 0),
      IFNULL(w.category_tag, ''), w.issue, w.year, w.title,
      w.cover_image_url, w.detail_image_url, w.draw_code,
      TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', 6)),
      TRIM(BOTH ',' FROM SUBSTRING_INDEX(REPLACE(IFNULL(w.draw_result, ''), ' ', ''), ',', -1)),
      w.draw_result, w.draw_at, '',
      w.is_current, w.status, w.sort, w.likes_count, w.comment_count, w.favorite_count, w.read_count,
      w.poll_enabled, w.poll_default_expand, IFNULL(w.recommend_info_ids, ''),
      IFNULL(w.created_at, NOW(3)), IFNULL(w.updated_at, NOW(3))
    FROM w_lottery_info w
    LEFT JOIN tk_lottery_info t
      ON t.special_lottery_id = w.special_lottery_id
     AND BINARY t.issue = BINARY w.issue
     AND BINARY t.title = BINARY w.title
    WHERE t.id IS NULL;
  END IF;
END $$

CALL migrate_w_to_tk_if_exists() $$
DROP PROCEDURE migrate_w_to_tk_if_exists $$

DELIMITER ;

-- =========================================================
-- C. 清理历史遗留字段（幂等）
-- =========================================================
-- 说明：
-- 1) 倒计时配置统一放在 tk_special_lottery.next_draw_at；
-- 2) tk_draw_record 仅保存“已开奖结果查询”，不再承载倒计时配置；
-- 3) 字段存在才删除，兼容旧库。

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
