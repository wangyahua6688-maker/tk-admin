USE nb_sys_001;

SET NAMES utf8mb4;

START TRANSACTION;

REPLACE INTO tk_special_lottery
  (id, name, code, current_issue, next_draw_at, live_enabled, live_status, live_stream_url, status, sort, created_at, updated_at)
VALUES
  (1, '澳彩', 'macau', '2026-068', '2026-03-28 19:46:00.000', 1, 'live', 'https://cph-p2p-msl.akamaized.net/hls/live/2000341/test/master.m3u8', 1, 1, NOW(3), NOW(3)),
  (2, '港彩', 'hongkong', '2026-188', '2026-03-28 18:19:00.000', 1, 'live', 'https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8', 1, 2, NOW(3), NOW(3));

REPLACE INTO tk_banner
  (id, title, image_url, link_url, type, position, positions, jump_type, jump_post_id, jump_url, content_html, status, sort, start_at, end_at, created_at, updated_at)
VALUES
  (1001, '首页今日策略', 'https://images.unsplash.com/photo-1514565131-fce0801e5785?auto=format&fit=crop&w=1200&q=80', '/lottery/6011', 'ad', 'home', 'home', 'external', 0, '/lottery/6011', NULL, 1, 1, '2026-01-01 00:00:00.000', '2026-12-31 23:59:59.000', NOW(3), NOW(3)),
  (1002, '港彩官方资讯', 'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1200&q=80', '/lottery/6013', 'official', 'home', 'home', 'external', 0, '/lottery/6013', NULL, 1, 2, '2026-01-01 00:00:00.000', '2026-12-31 23:59:59.000', NOW(3), NOW(3)),
  (1003, '图纸详情位推荐', 'https://images.unsplash.com/photo-1519389950473-47ba0277781c?auto=format&fit=crop&w=1200&q=80', '/lottery/6011', 'official', 'lottery_detail', 'lottery_detail', 'external', 0, '/lottery/6011', NULL, 1, 1, '2026-01-01 00:00:00.000', '2026-12-31 23:59:59.000', NOW(3), NOW(3));

REPLACE INTO tk_broadcast
  (id, title, content, status, sort, created_at, updated_at)
VALUES
  (1101, '系统维护通知', '系统将于凌晨 03:00 做常规维护，预计 10 分钟。', 1, 1, NOW(3), NOW(3)),
  (1102, '开奖提示', '澳彩与港彩开奖区已切换为真实测试数据，可直接联调。', 1, 2, NOW(3), NOW(3));

REPLACE INTO tk_home_popup
  (id, title, content, image_url, button_text, button_link, position, show_once, status, sort, start_at, end_at, created_at, updated_at)
VALUES
  (1201, '欢迎来到测试环境', '当前站点已写入首页、开奖区、论坛与高手榜测试数据。', '', '知道了', '/forum', 'home', 0, 1, 1, '2026-01-01 00:00:00.000', '2026-12-31 23:59:59.000', NOW(3), NOW(3));

REPLACE INTO tk_external_link
  (id, name, url, position, icon_url, group_key, status, sort, created_at, updated_at)
VALUES
  (1301, '开奖历史', '/history', 'home_external', '', 'default', 1, 1, NOW(3), NOW(3)),
  (1302, '高手论坛', '/forum', 'home_external', '', 'default', 1, 2, NOW(3), NOW(3)),
  (1303, '高手推荐', '/experts', 'home_nav', '', 'default', 1, 1, NOW(3), NOW(3)),
  (1304, '开奖详情推荐', '/lottery/6011', 'lottery_detail', '', 'default', 1, 1, NOW(3), NOW(3));

REPLACE INTO tk_lottery_category
  (id, category_key, name, search_keywords, show_on_home, status, sort, created_at, updated_at)
VALUES
  (1, 'jiuxiao', '九肖中特', '九肖 中特 推荐', 1, 1, 1, NOW(3), NOW(3)),
  (2, 'wuma', '五码中特', '五码 中特 稳胆', 1, 1, 2, NOW(3), NOW(3)),
  (3, 'pingteyi', '平特一肖', '平特 一肖 稳胆', 1, 1, 3, NOW(3), NOW(3)),
  (4, 'zhengma', '正码', '正码 两面 波色', 1, 1, 4, NOW(3), NOW(3));

REPLACE INTO tk_lottery_info
  (id, special_lottery_id, category_id, category_tag, issue, year, title, cover_image_url, detail_image_url, draw_code, normal_draw_result, special_draw_result, draw_result, draw_at, playback_url, is_current, status, sort, likes_count, comment_count, favorite_count, read_count, poll_enabled, poll_default_expand, recommend_info_ids, created_at, updated_at)
VALUES
  (6011, 1, 1, 'jiuxiao', '2026-068', 2026, '澳门068期九肖稳定组合', 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?auto=format&fit=crop&w=900&q=80', 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?auto=format&fit=crop&w=1400&q=80', '13 49 39 36 05 42 32 34 26 33', '02,11,17,24,32,48', '09', '02 11 17 24 32 48 09', '2026-03-28 21:30:00.000', 'https://cdn.example.com/replay/macau/2026-068.m3u8', 1, 1, 1, 178, 7, 29, 3688, 1, 0, '6012,6013', NOW(3), NOW(3)),
  (6012, 1, 2, 'wuma', '2026-067', 2026, '澳门067期五码精选', 'https://images.unsplash.com/photo-1511497584788-876760111969?auto=format&fit=crop&w=900&q=80', 'https://images.unsplash.com/photo-1511497584788-876760111969?auto=format&fit=crop&w=1400&q=80', '01 03 07 14 29 33 36 41 44 48', '03,14,22,27,35,46', '18', '03 14 22 27 35 46 18', '2026-03-27 21:30:00.000', 'https://cdn.example.com/replay/macau/2026-067.m3u8', 0, 1, 2, 123, 4, 18, 2290, 1, 0, '6011,6013', NOW(3), NOW(3)),
  (6013, 2, 3, 'pingteyi', '2026-188', 2026, '港彩188期平特一肖', 'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=900&q=80', 'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=1400&q=80', '05 07 09 12 24 29 31 36 41 46', '01,07,15,29,38,46', '24', '01 07 15 29 38 46 24', '2026-03-28 18:19:00.000', 'https://cdn.example.com/replay/hongkong/2026-188.m3u8', 1, 1, 3, 86, 3, 11, 1870, 1, 0, '6011,6012', NOW(3), NOW(3)),
  (6014, 2, 4, 'zhengma', '2026-187', 2026, '港彩187期正码走势', 'https://images.unsplash.com/photo-1470770841072-f978cf4d019e?auto=format&fit=crop&w=900&q=80', 'https://images.unsplash.com/photo-1470770841072-f978cf4d019e?auto=format&fit=crop&w=1400&q=80', '02 08 11 18 23 32 35 39 44 47', '04,16,21,30,33,45', '11', '04 16 21 30 33 45 11', '2026-03-27 18:19:00.000', 'https://cdn.example.com/replay/hongkong/2026-187.m3u8', 0, 1, 4, 64, 2, 9, 1520, 1, 0, '6013,6011', NOW(3), NOW(3));

REPLACE INTO tk_draw_record
  (id, special_lottery_id, issue, year, draw_at, normal_draw_result, special_draw_result, draw_result, draw_labels, color_labels, zodiac_labels, wuxing_labels, playback_url, special_single_double, special_big_small, sum_single_double, sum_big_small, recommend_six, recommend_four, recommend_one, recommend_ten, special_code, normal_code, zheng1, zheng2, zheng3, zheng4, zheng5, zheng6, is_current, status, sort, created_at, updated_at)
VALUES
  (7011, 1, '2026-067', 2026, '2026-03-27 21:30:00.000', '03,14,22,27,35,46', '18', '03 14 22 27 35 46 18', '龙/金,蛇/水,鸡/水,龙/土,猴/土,鸡/木,牛/金', '蓝波,蓝波,绿波,绿波,红波,红波,红波', '龙,蛇,鸡,龙,猴,鸡,牛', '金,水,水,土,土,木,金', 'https://cdn.example.com/replay/macau/2026-067.m3u8', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', 0, 1, 2, NOW(3), NOW(3)),
  (7012, 1, '2026-068', 2026, '2026-03-28 21:30:00.000', '02,11,17,24,32,48', '09', '02 11 17 24 32 48 09', '蛇/火,猴/金,虎/火,羊/木,猪/火,羊/火,狗/火', '红波,绿波,绿波,红波,绿波,蓝波,蓝波', '蛇,猴,虎,羊,猪,羊,狗', '火,金,火,木,火,火,火', 'https://cdn.example.com/replay/macau/2026-068.m3u8', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', 1, 1, 1, NOW(3), NOW(3)),
  (7021, 2, '2026-187', 2026, '2026-03-27 18:19:00.000', '04,16,21,30,33,45', '11', '04 16 21 30 33 45 11', '兔/金,兔/木,狗/水,牛/水,鼠/土,狗/木,猴/金', '蓝波,绿波,绿波,红波,绿波,红波,绿波', '兔,兔,狗,牛,鼠,狗,猴', '金,木,水,水,土,木,金', 'https://cdn.example.com/replay/hongkong/2026-187.m3u8', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', 0, 1, 2, NOW(3), NOW(3)),
  (7022, 2, '2026-188', 2026, '2026-03-28 18:19:00.000', '01,07,15,29,38,46', '24', '01 07 15 29 38 46 24', '马/火,鼠/木,龙/木,虎/水,蛇/木,鸡/木,羊/木', '红波,红波,蓝波,红波,绿波,红波,红波', '马,鼠,龙,虎,蛇,鸡,羊', '火,木,木,水,木,木,木', 'https://cdn.example.com/replay/hongkong/2026-188.m3u8', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', 1, 1, 1, NOW(3), NOW(3));

REPLACE INTO tk_lottery_option
  (id, lottery_info_id, option_name, votes, sort, created_at, updated_at)
VALUES
  (9101, 6011, '鼠', 6, 1, NOW(3), NOW(3)),
  (9102, 6011, '牛', 8, 2, NOW(3), NOW(3)),
  (9103, 6011, '虎', 12, 3, NOW(3), NOW(3)),
  (9104, 6011, '兔', 9, 4, NOW(3), NOW(3)),
  (9105, 6011, '龙', 10, 5, NOW(3), NOW(3)),
  (9106, 6011, '蛇', 13, 6, NOW(3), NOW(3)),
  (9107, 6011, '马', 7, 7, NOW(3), NOW(3)),
  (9108, 6011, '羊', 15, 8, NOW(3), NOW(3)),
  (9109, 6011, '猴', 11, 9, NOW(3), NOW(3)),
  (9110, 6011, '鸡', 5, 10, NOW(3), NOW(3)),
  (9111, 6011, '狗', 9, 11, NOW(3), NOW(3)),
  (9112, 6011, '猪', 6, 12, NOW(3), NOW(3)),
  (9113, 6013, '鼠', 5, 1, NOW(3), NOW(3)),
  (9114, 6013, '牛', 9, 2, NOW(3), NOW(3)),
  (9115, 6013, '虎', 16, 3, NOW(3), NOW(3)),
  (9116, 6013, '兔', 8, 4, NOW(3), NOW(3)),
  (9117, 6013, '龙', 12, 5, NOW(3), NOW(3)),
  (9118, 6013, '蛇', 7, 6, NOW(3), NOW(3)),
  (9119, 6013, '马', 14, 7, NOW(3), NOW(3)),
  (9120, 6013, '羊', 10, 8, NOW(3), NOW(3)),
  (9121, 6013, '猴', 6, 9, NOW(3), NOW(3)),
  (9122, 6013, '鸡', 4, 10, NOW(3), NOW(3)),
  (9123, 6013, '狗', 11, 11, NOW(3), NOW(3)),
  (9124, 6013, '猪', 3, 12, NOW(3), NOW(3));

REPLACE INTO tk_users
  (id, username, phone, nickname, avatar, password_hash, register_source, last_login_at, user_type, fans_count, following_count, growth_value, read_post_count, status, created_at, updated_at)
VALUES
  (3001, 'tk_official_helper', '13800138001', 'TK官方助手', 'https://i.pravatar.cc/150?img=12', '', 'admin', NOW(3), 'official', 980, 12, 2860, 1360, 1, NOW(3), NOW(3)),
  (3002, 'tk_reporter', '13800138002', '公告播报员', 'https://i.pravatar.cc/150?img=13', '', 'admin', NOW(3), 'official', 860, 8, 2310, 1190, 1, NOW(3), NOW(3)),
  (3003, 'caiyou_ajie', '13800138003', '彩友阿杰', 'https://i.pravatar.cc/150?img=31', '', 'sms', NOW(3), 'natural', 1320, 42, 3560, 2140, 1, NOW(3), NOW(3)),
  (3004, 'steady_xiaoyu', '13800138004', '稳胆小雨', 'https://i.pravatar.cc/150?img=41', '', 'sms', NOW(3), 'natural', 740, 21, 2050, 980, 1, NOW(3), NOW(3)),
  (3005, 'smart_assistant', '13800138005', '智投助手', 'https://i.pravatar.cc/150?img=52', '', 'admin', NOW(3), 'robot', 1180, 6, 2980, 1730, 1, NOW(3), NOW(3)),
  (3006, 'laochen_live', '13800138006', '老陈实战', 'https://i.pravatar.cc/150?img=56', '', 'sms', NOW(3), 'natural', 905, 18, 2670, 1490, 1, NOW(3), NOW(3));

REPLACE INTO tk_post_article
  (id, lottery_info_id, user_id, title, cover_image, content, is_official, status, created_at, updated_at)
VALUES
  (8101, 6011, 3001, '第068期官方解析：六码组合与稳胆', 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?auto=format&fit=crop&w=900&q=80', '本期关注前区连号与两组同尾号，建议先看开奖详情再做讨论。', 1, 1, '2026-03-28 08:20:00.000', NOW(3)),
  (8102, 6013, 3002, '港彩188期官方提醒：注意尾数分布', 'https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=900&q=80', '港彩临近开奖，尾数与波色分布值得重点关注。', 1, 1, '2026-03-28 10:30:00.000', NOW(3)),
  (8103, 6011, 3003, '澳门068期：我看好的九肖组合', 'https://images.unsplash.com/photo-1465146344425-f00d5f5c8f07?auto=format&fit=crop&w=900&q=80', '这期我更看好蛇、猴、虎、羊、猪，防一手狗。', 0, 1, '2026-03-28 12:10:00.000', NOW(3)),
  (8104, 6012, 3003, '2025年末复盘：节奏比号码更重要', 'https://images.unsplash.com/photo-1500534623283-312aade485b7?auto=format&fit=crop&w=900&q=80', '回看历史走势，执行纪律往往比单期预测更重要。', 0, 1, '2026-03-27 22:10:00.000', NOW(3)),
  (8105, 6013, 3004, '港彩068期：我看好的四肖', 'https://images.unsplash.com/photo-1473448912268-2022ce9509d8?auto=format&fit=crop&w=900&q=80', '这期我会重点关注虎、龙、马、羊的组合。', 0, 1, '2026-03-27 20:10:00.000', NOW(3)),
  (8106, 6014, 3005, '机器人复盘：五码中特回报率统计', 'https://images.unsplash.com/photo-1511497584788-876760111969?auto=format&fit=crop&w=900&q=80', '结合最近20期结果，五码中特的风险回报比依然靠前。', 0, 1, '2026-03-27 18:20:00.000', NOW(3)),
  (8107, 6011, 3006, '老陈实战：正码组合拆解', 'https://images.unsplash.com/photo-1522202176988-66273c2fd55f?auto=format&fit=crop&w=900&q=80', '胆码、拖码和波色搭配需要一起看，单看一个维度容易偏。', 0, 1, '2026-03-28 14:45:00.000', NOW(3));

REPLACE INTO tk_comment
  (id, post_id, lottery_info_id, user_id, parent_id, content, likes, status, created_at, updated_at)
VALUES
  (9001, 8101, 0, 3003, 0, '这一期官方给出的六码思路和我自己的判断接近。', 12, 1, '2026-03-28 08:40:00.000', NOW(3)),
  (9002, 8101, 0, 3005, 0, '前区连号值得防一手，后区更看好奇数组合。', 16, 1, '2026-03-28 09:10:00.000', NOW(3)),
  (9003, 8103, 0, 3004, 0, '九肖里我也会保留羊和狗，整体方向一致。', 8, 1, '2026-03-28 12:50:00.000', NOW(3)),
  (9004, 8103, 0, 3006, 0, '这期波色分布也可以一起看，别只盯生肖。', 10, 1, '2026-03-28 13:12:00.000', NOW(3)),
  (9005, 8106, 0, 3005, 0, '五码中特这类玩法更适合看回报率，不是只看命中率。', 18, 1, '2026-03-27 18:45:00.000', NOW(3)),
  (9006, 8107, 0, 3006, 0, '正码拆解里两面和波色最好一起看。', 21, 1, '2026-03-28 15:05:00.000', NOW(3)),
  (9007, 0, 6011, 3001, 0, '系统评论：澳门068期开奖后将同步更新玩法结果。', 5, 1, '2026-03-28 21:35:00.000', NOW(3)),
  (9008, 0, 6011, 3003, 0, '网友评论：这期蛇、猴、虎整体命中表现不错。', 9, 1, '2026-03-28 21:50:00.000', NOW(3)),
  (9009, 0, 6011, 3005, 0, '热门评论：本期五行分布偏火木，下期可能会回补金水。', 13, 1, '2026-03-28 22:05:00.000', NOW(3)),
  (9010, 0, 6013, 3002, 0, '系统评论：港彩188期的平特与波色结果已经同步。', 4, 1, '2026-03-28 18:25:00.000', NOW(3)),
  (9011, 0, 6013, 3004, 0, '网友评论：港彩这期红波占比偏高，下一期可以留意蓝绿。', 7, 1, '2026-03-28 18:40:00.000', NOW(3));

REPLACE INTO tk_sms_channel
  (id, provider, channel_name, access_key, access_secret, endpoint, sign_name, template_code_login, template_code_register, daily_limit, minute_limit, code_ttl_seconds, mock_mode, status, created_at, updated_at)
VALUES
  (1, 'mock', '默认模拟通道', '', '', '', 'TK', 'LOGIN', 'REGISTER', 20, 1, 300, 1, 1, NOW(3), NOW(3));

COMMIT;
