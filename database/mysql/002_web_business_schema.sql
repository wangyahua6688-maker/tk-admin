-- TK 平台业务表结构（tk_ 前缀）
-- 说明：
-- 1. 本脚本仅包含业务域表，不包含 RBAC 系统表；
-- 2. 执行前请先选择目标数据库：USE your_db;

CREATE TABLE IF NOT EXISTS tk_users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  username VARCHAR(64) NOT NULL COMMENT '用户名（唯一）',
  phone VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号（用于验证码登录，唯一）',
  nickname VARCHAR(64) NOT NULL DEFAULT '' COMMENT '昵称',
  avatar VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像地址',
  password_hash VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码哈希（bcrypt）',
  register_source VARCHAR(20) NOT NULL DEFAULT 'password' COMMENT '注册来源：password/sms/admin/import',
  last_login_at DATETIME(3) NULL COMMENT '最近登录时间',
  user_type VARCHAR(20) NOT NULL DEFAULT 'natural' COMMENT '用户类型：natural自然用户；official官方账号；robot机器人账号',
  fans_count BIGINT NOT NULL DEFAULT 0 COMMENT '粉丝数',
  following_count BIGINT NOT NULL DEFAULT 0 COMMENT '关注数',
  growth_value BIGINT NOT NULL DEFAULT 0 COMMENT '成长值',
  read_post_count BIGINT NOT NULL DEFAULT 0 COMMENT '阅读帖子数',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_tk_users_username (username),
  UNIQUE KEY uk_tk_users_phone (phone),
  KEY idx_tk_users_user_type (user_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS tk_banner (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  title VARCHAR(120) NOT NULL COMMENT 'Banner标题',
  image_url VARCHAR(255) NOT NULL COMMENT 'Banner图片地址',
  link_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '兼容字段：历史跳转地址',
  type VARCHAR(32) NOT NULL COMMENT 'Banner类型：ad广告；official官方通知',
  position VARCHAR(32) NOT NULL COMMENT '兼容字段：主展示位置（取positions第一个）',
  positions VARCHAR(255) NOT NULL DEFAULT '' COMMENT '展示位置，多选逗号分隔：home,lottery_detail,post_detail',
  jump_type VARCHAR(20) NOT NULL DEFAULT 'none' COMMENT '跳转类型：none不跳转；post关联帖子；external外链；custom自定义内容',
  jump_post_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联帖子ID（jump_type=post）',
  jump_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '跳转外链地址（jump_type=external）',
  content_html LONGTEXT NULL COMMENT '自定义富文本内容（jump_type=custom）',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  start_at DATETIME(3) NULL COMMENT '生效开始时间',
  end_at DATETIME(3) NULL COMMENT '生效结束时间',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_banner_type (type),
  KEY idx_tk_banner_position (position),
  KEY idx_tk_banner_jump_type (jump_type),
  KEY idx_tk_banner_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Banner配置表';

CREATE TABLE IF NOT EXISTS tk_broadcast (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  title VARCHAR(120) NOT NULL COMMENT '广播标题',
  content VARCHAR(500) NOT NULL COMMENT '广播内容',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_broadcast_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统广播表';

CREATE TABLE IF NOT EXISTS tk_special_lottery (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  name VARCHAR(64) NOT NULL COMMENT '彩种名称',
  code VARCHAR(32) NOT NULL COMMENT '彩种编码（唯一）',
  current_issue VARCHAR(32) NOT NULL DEFAULT '' COMMENT '当前期号',
  next_draw_at DATETIME(3) NOT NULL COMMENT '下期开奖时间',
  live_enabled TINYINT NOT NULL DEFAULT 0 COMMENT '直播开关：1开启；0关闭',
  live_status VARCHAR(16) NOT NULL DEFAULT 'pending' COMMENT '直播状态：pending未开始；live直播中；ended已结束',
  live_stream_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '外部直播流地址',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_tk_special_lottery_code (code),
  KEY idx_tk_special_lottery_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='彩种配置表';

CREATE TABLE IF NOT EXISTS tk_lottery_category (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  category_key VARCHAR(32) NOT NULL COMMENT '分类键（唯一）',
  name VARCHAR(32) NOT NULL COMMENT '分类名称',
  search_keywords VARCHAR(255) NOT NULL DEFAULT '' COMMENT '搜索关键字（空格/逗号分隔）',
  show_on_home TINYINT NOT NULL DEFAULT 1 COMMENT '是否首页展示：1展示；0仅在更多分类中展示',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_tk_lottery_category_key (category_key),
  KEY idx_tk_lottery_category_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='图库分类配置表';

CREATE TABLE IF NOT EXISTS tk_lottery_info (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  special_lottery_id BIGINT UNSIGNED NOT NULL COMMENT '所属彩种ID（关联tk_special_lottery.id，0表示不绑定彩种）',
  category_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '图库分类ID（关联tk_lottery_category.id）',
  category_tag VARCHAR(32) NOT NULL DEFAULT '' COMMENT '分类标识兼容字段（通常等于category_key）',
  issue VARCHAR(32) NOT NULL COMMENT '期号',
  year INT NOT NULL COMMENT '年份（如2026）',
  title VARCHAR(120) NOT NULL COMMENT '标题',
  cover_image_url VARCHAR(255) NOT NULL COMMENT '列表封面图地址',
  detail_image_url VARCHAR(255) NOT NULL COMMENT '详情图地址',
  draw_code VARCHAR(120) NOT NULL DEFAULT '' COMMENT '暗码',
  normal_draw_result VARCHAR(64) NOT NULL DEFAULT '' COMMENT '普通号码（6个，逗号分隔）',
  special_draw_result VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特别号码（1个）',
  draw_result VARCHAR(120) NOT NULL DEFAULT '' COMMENT '兼容字段：完整开奖号码（普通6个+特别号）',
  draw_at DATETIME(3) NOT NULL COMMENT '开奖时间',
  playback_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '直播回放地址（直播结束后录入）',
  is_current TINYINT NOT NULL DEFAULT 0 COMMENT '是否当前期：1是；0否',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  likes_count BIGINT NOT NULL DEFAULT 0 COMMENT '点赞数（详情页展示）',
  comment_count BIGINT NOT NULL DEFAULT 0 COMMENT '评论数（详情页展示）',
  favorite_count BIGINT NOT NULL DEFAULT 0 COMMENT '收藏数（详情页展示）',
  read_count BIGINT NOT NULL DEFAULT 0 COMMENT '阅读数（详情页展示）',
  poll_enabled TINYINT NOT NULL DEFAULT 1 COMMENT '投票开关：1显示；0隐藏',
  poll_default_expand TINYINT NOT NULL DEFAULT 0 COMMENT '投票默认展开：1展开；0收起',
  recommend_info_ids VARCHAR(255) NOT NULL DEFAULT '' COMMENT '推荐图纸ID列表，逗号分隔',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_lottery_info_special (special_lottery_id),
  KEY idx_tk_lottery_info_category_id (category_id),
  KEY idx_tk_lottery_info_category_tag (category_tag),
  KEY idx_tk_lottery_info_issue (issue),
  KEY idx_tk_lottery_info_year (year),
  KEY idx_tk_lottery_info_draw_at (draw_at),
  KEY idx_tk_lottery_info_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='图库图纸内容与竞猜配置表（不承载开奖区历史主数据）';

CREATE TABLE IF NOT EXISTS tk_draw_record (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  special_lottery_id BIGINT UNSIGNED NOT NULL COMMENT '所属彩种ID（关联tk_special_lottery.id）',
  issue VARCHAR(32) NOT NULL COMMENT '开奖期号（如2026-063）',
  year INT NOT NULL COMMENT '年份（如2026）',
  draw_at DATETIME(3) NOT NULL COMMENT '开奖时间',
  normal_draw_result VARCHAR(64) NOT NULL DEFAULT '' COMMENT '普通号码（6个，逗号分隔）',
  special_draw_result VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特别号码（1个）',
  draw_result VARCHAR(120) NOT NULL DEFAULT '' COMMENT '兼容字段：完整开奖号码（普通6个+特别号）',
  draw_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码标签（与号码一一对应，格式示例：羊/土,蛇/金）',
  zodiac_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码对应属相标签（与号码一一对应，示例：羊,蛇,马）',
  wuxing_labels VARCHAR(255) NOT NULL DEFAULT '' COMMENT '号码对应五行标签（与号码一一对应，示例：土,金,火）',
  playback_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '开奖回放地址（直播结束后录入）',
  special_single_double VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特码单双（如：双）',
  special_big_small VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特码大小（如：大）',
  sum_single_double VARCHAR(16) NOT NULL DEFAULT '' COMMENT '总和单双（如：双）',
  sum_big_small VARCHAR(16) NOT NULL DEFAULT '' COMMENT '总和大小（如：大）',
  recommend_six VARCHAR(120) NOT NULL DEFAULT '' COMMENT '六肖推荐（空格分隔）',
  recommend_four VARCHAR(120) NOT NULL DEFAULT '' COMMENT '四肖推荐（空格分隔）',
  recommend_one VARCHAR(32) NOT NULL DEFAULT '' COMMENT '一肖推荐',
  recommend_ten VARCHAR(255) NOT NULL DEFAULT '' COMMENT '十码推荐（空格分隔）',
  special_code VARCHAR(16) NOT NULL DEFAULT '' COMMENT '特码数字',
  normal_code VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正码（逗号分隔）',
  zheng1 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正1特描述（如：大双,合双,尾大,蓝波）',
  zheng2 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正2特描述',
  zheng3 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正3特描述',
  zheng4 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正4特描述',
  zheng5 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正5特描述',
  zheng6 VARCHAR(120) NOT NULL DEFAULT '' COMMENT '正6特描述',
  is_current TINYINT NOT NULL DEFAULT 0 COMMENT '是否当前期：1是；0否',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_tk_draw_record_issue (special_lottery_id, issue),
  KEY idx_tk_draw_record_special (special_lottery_id),
  KEY idx_tk_draw_record_year (year),
  KEY idx_tk_draw_record_draw_at (draw_at),
  KEY idx_tk_draw_record_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='开奖区开奖记录表（首页开奖区/历史开奖/开奖详情）';

CREATE TABLE IF NOT EXISTS tk_lottery_option (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  lottery_info_id BIGINT UNSIGNED NOT NULL COMMENT '开奖内容ID（关联tk_lottery_info.id）',
  option_name VARCHAR(32) NOT NULL COMMENT '投票选项名称（生肖）',
  votes BIGINT NOT NULL DEFAULT 0 COMMENT '票数',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_lottery_option_info (lottery_info_id),
  KEY idx_tk_lottery_option_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='投票选项表';

CREATE TABLE IF NOT EXISTS tk_lottery_vote_record (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  lottery_info_id BIGINT UNSIGNED NOT NULL COMMENT '开奖内容ID',
  option_id BIGINT UNSIGNED NOT NULL COMMENT '投票选项ID',
  voter_hash VARCHAR(64) NOT NULL COMMENT '投票指纹哈希（防刷）',
  device_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '设备ID（前端传入）',
  client_ip VARCHAR(64) NOT NULL DEFAULT '' COMMENT '客户端IP',
  user_agent VARCHAR(255) NOT NULL DEFAULT '' COMMENT '客户端UA',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_tk_lottery_vote_record_unique (lottery_info_id, voter_hash),
  KEY idx_tk_lottery_vote_record_option (option_id),
  KEY idx_tk_lottery_vote_record_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='投票记录表';

CREATE TABLE IF NOT EXISTS tk_post_article (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  lottery_info_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '关联开奖内容ID（关联tk_lottery_info.id，0表示不关联）',
  user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '发帖用户ID（关联tk_users.id）',
  title VARCHAR(160) NOT NULL COMMENT '帖子标题',
  cover_image VARCHAR(255) NOT NULL DEFAULT '' COMMENT '封面图地址',
  content TEXT NULL COMMENT '帖子富文本内容',
  is_official TINYINT NOT NULL DEFAULT 0 COMMENT '帖子类型：1官方发帖；0网友发帖',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_post_article_lottery_info_id (lottery_info_id),
  KEY idx_tk_post_article_user_id (user_id),
  KEY idx_tk_post_article_official (is_official),
  KEY idx_tk_post_article_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='帖子表';

CREATE TABLE IF NOT EXISTS tk_comment (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  post_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '帖子ID（关联tk_post_article.id，0表示非帖子评论）',
  lottery_info_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '开奖内容ID（关联tk_lottery_info.id，0表示非彩种详情评论）',
  user_id BIGINT UNSIGNED NOT NULL COMMENT '评论用户ID（关联tk_users.id）',
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父评论ID，0表示一级评论',
  content VARCHAR(1000) NOT NULL COMMENT '评论内容',
  likes BIGINT NOT NULL DEFAULT 0 COMMENT '点赞数',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_comment_post_id (post_id),
  KEY idx_tk_comment_lottery_info_id (lottery_info_id),
  KEY idx_tk_comment_user_id (user_id),
  KEY idx_tk_comment_parent_id (parent_id),
  KEY idx_tk_comment_status_likes (status, likes)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='评论表';

CREATE TABLE IF NOT EXISTS tk_external_link (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  name VARCHAR(80) NOT NULL COMMENT '外链名称',
  url VARCHAR(255) NOT NULL COMMENT '外链地址',
  position VARCHAR(32) NOT NULL COMMENT '展示位置',
  icon_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '图标地址（用于金刚导航）',
  group_key VARCHAR(32) NOT NULL DEFAULT '' COMMENT '分组键（如：aocai/hkcai/default）',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_external_link_position (position),
  KEY idx_tk_external_link_group_key (group_key),
  KEY idx_tk_external_link_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='外链配置表';

CREATE TABLE IF NOT EXISTS tk_home_popup (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  title VARCHAR(120) NOT NULL COMMENT '弹窗标题',
  content TEXT NULL COMMENT '弹窗内容（支持富文本）',
  image_url VARCHAR(255) NOT NULL DEFAULT '' COMMENT '弹窗图片地址',
  button_text VARCHAR(40) NOT NULL DEFAULT '' COMMENT '按钮文案',
  button_link VARCHAR(255) NOT NULL DEFAULT '' COMMENT '按钮跳转地址',
  position VARCHAR(32) NOT NULL DEFAULT 'home' COMMENT '展示位置：home-首页',
  show_once TINYINT NOT NULL DEFAULT 1 COMMENT '是否单设备仅展示一次：1是；0否',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  sort INT NOT NULL DEFAULT 0 COMMENT '排序值，越小越靠前',
  start_at DATETIME(3) NULL COMMENT '生效开始时间',
  end_at DATETIME(3) NULL COMMENT '生效结束时间',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_home_popup_position (position),
  KEY idx_tk_home_popup_status_sort (status, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='首页首屏弹窗配置表';

CREATE TABLE IF NOT EXISTS tk_sms_channel (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  provider VARCHAR(32) NOT NULL DEFAULT 'custom' COMMENT '短信服务商：aliyun/tencent/twilio/custom/mock',
  channel_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '通道名称',
  access_key VARCHAR(128) NOT NULL DEFAULT '' COMMENT '服务凭证 AccessKey',
  access_secret VARCHAR(255) NOT NULL DEFAULT '' COMMENT '服务凭证 AccessSecret',
  endpoint VARCHAR(255) NOT NULL DEFAULT '' COMMENT '网关地址',
  sign_name VARCHAR(64) NOT NULL DEFAULT '' COMMENT '短信签名',
  template_code_login VARCHAR(64) NOT NULL DEFAULT '' COMMENT '登录验证码模板编码',
  template_code_register VARCHAR(64) NOT NULL DEFAULT '' COMMENT '注册验证码模板编码',
  daily_limit INT NOT NULL DEFAULT 20 COMMENT '单手机号日发送上限',
  minute_limit INT NOT NULL DEFAULT 1 COMMENT '单手机号分钟发送上限',
  code_ttl_seconds INT NOT NULL DEFAULT 300 COMMENT '验证码有效时长（秒）',
  mock_mode TINYINT NOT NULL DEFAULT 1 COMMENT '模拟发送开关：1模拟；0真实发送',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用；0停用',
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_tk_sms_channel_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='短信通道配置表';
