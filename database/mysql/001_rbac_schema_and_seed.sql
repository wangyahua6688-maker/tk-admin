-- =========================================================
-- 后台管理系统 RBAC 初始化脚本（MySQL 8.x）
-- 说明：
-- 1. 本脚本包含表结构 + 初始化数据（角色/权限/菜单/关联关系）。
-- 2. 表结构与当前 Go 代码（GORM 模型）保持一致。
-- 3. 执行前请确认数据库字符集为 utf8mb4。
-- =========================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- -------------------------
-- 1) 用户表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(100) NOT NULL,
  `password_hash` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) DEFAULT '',
  `avatar` VARCHAR(255) DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1,
  `refresh_token` VARCHAR(512) DEFAULT NULL,
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统用户表';


-- -------------------------
-- 2) 角色表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_roles` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) DEFAULT NULL,
  `name` VARCHAR(100) NOT NULL,
  `code` VARCHAR(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_roles_name` (`name`),
  UNIQUE KEY `uk_roles_code` (`code`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- -------------------------
-- 3) 权限表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_permissions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) DEFAULT NULL,
  `name` VARCHAR(100) NOT NULL,
  `code` VARCHAR(100) NOT NULL,
  `type` VARCHAR(50) DEFAULT 'api',
  `method` VARCHAR(50) DEFAULT '',
  `path` VARCHAR(200) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_permissions_name` (`name`),
  UNIQUE KEY `uk_permissions_code` (`code`),
  KEY `idx_permissions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- -------------------------
-- 4) 菜单表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_menus` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(200) NOT NULL,
  `path` VARCHAR(200) NOT NULL,
  `icon` VARCHAR(100) DEFAULT '',
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `component` VARCHAR(200) DEFAULT '',
  `order_num` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `deleted_at` DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_menus_parent_id` (`parent_id`),
  KEY `idx_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- -------------------------
-- 5) 用户-角色关联表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_user_roles` (
  `user_id` BIGINT UNSIGNED NOT NULL,
  `role_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`),
  KEY `idx_user_roles_role_id` (`role_id`),
  CONSTRAINT `fk_sys_user_roles_user_id` FOREIGN KEY (`user_id`) REFERENCES `sys_users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sys_user_roles_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- -------------------------
-- 6) 角色-权限关联表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_role_permissions` (
  `role_id` BIGINT UNSIGNED NOT NULL,
  `permission_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`role_id`, `permission_id`),
  KEY `idx_role_permissions_permission_id` (`permission_id`),
  CONSTRAINT `fk_sys_role_permissions_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sys_role_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `sys_permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- -------------------------
-- 7) 菜单-权限关联表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_menu_permissions` (
  `menu_id` BIGINT UNSIGNED NOT NULL,
  `permission_id` BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (`menu_id`, `permission_id`),
  KEY `idx_menu_permissions_permission_id` (`permission_id`),
  CONSTRAINT `fk_sys_menu_permissions_menu_id` FOREIGN KEY (`menu_id`) REFERENCES `sys_menus` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sys_menu_permissions_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `sys_permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单权限关联表';

-- -------------------------
-- 8) 登录日志表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_login_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `username` VARCHAR(100) DEFAULT '',
  `ip` VARCHAR(64) DEFAULT '',
  `device` VARCHAR(255) DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_login_logs_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志';

-- -------------------------
-- 9) 刷新令牌记录表（审计用途）
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_refresh_token_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `device_id` VARCHAR(255) DEFAULT '',
  `token` VARCHAR(1024) NOT NULL,
  `expires_at` DATETIME(3) NOT NULL,
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_refresh_token_records_user_id` (`user_id`),
  KEY `idx_refresh_token_records_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='刷新令牌审计表';

-- -------------------------
-- 10) 系统消息表
-- -------------------------
CREATE TABLE IF NOT EXISTS `sys_system_messages` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `title` VARCHAR(200) NOT NULL,
  `content` TEXT NOT NULL,
  `level` VARCHAR(20) NOT NULL DEFAULT 'info',
  `is_read` TINYINT(1) NOT NULL DEFAULT 0,
  `operator_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `biz_type` VARCHAR(50) DEFAULT '',
  `biz_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_sys_msg_user_id` (`user_id`),
  KEY `idx_sys_msg_operator_id` (`operator_id`),
  KEY `idx_sys_msg_biz` (`biz_type`, `biz_id`),
  KEY `idx_sys_msg_user_read` (`user_id`, `is_read`),
  CONSTRAINT `fk_sys_system_messages_user_id` FOREIGN KEY (`user_id`) REFERENCES `sys_users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统消息表';

-- -------------------------
-- 11) 初始化角色
-- -------------------------
INSERT INTO `sys_roles` (`id`, `name`, `code`, `created_at`, `updated_at`)
VALUES
  (1, '超级管理员', 'admin', NOW(3), NOW(3)),
  (2, '安全审计员', 'security_auditor', NOW(3), NOW(3)),
  (3, '运营人员', 'operator', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE `updated_at` = VALUES(`updated_at`);

-- -------------------------
-- 12) 初始化权限（与当前后端路由保持一致）
-- -------------------------
INSERT INTO `sys_permissions` (`id`, `name`, `code`, `type`, `method`, `path`, `created_at`, `updated_at`)
VALUES
  (1,  '角色列表',                 'sys:role:list',              'api', 'GET',    '/api/roles/',                     NOW(3), NOW(3)),
  (2,  '创建角色',                 'sys:role:create',            'api', 'POST',   '/api/roles/',                     NOW(3), NOW(3)),
  (3,  '角色详情',                 'sys:role:view',              'api', 'GET',    '/api/roles/:id',                  NOW(3), NOW(3)),
  (4,  '更新角色',                 'sys:role:update',            'api', 'PUT',    '/api/roles/:id',                  NOW(3), NOW(3)),
  (5,  '删除角色',                 'sys:role:delete',            'api', 'DELETE', '/api/roles/:id',                  NOW(3), NOW(3)),
  (6,  '角色权限查看',             'sys:role_permission:view',   'api', 'GET',    '/api/roles/:id/permissions',      NOW(3), NOW(3)),
  (7,  '角色权限绑定',             'sys:role_permission:bind',   'api', 'PUT',    '/api/roles/:id/permissions',      NOW(3), NOW(3)),
  (8,  '权限列表',                 'sys:permission:list',        'api', 'GET',    '/api/permissions/',               NOW(3), NOW(3)),
  (9,  '创建权限',                 'sys:permission:create',      'api', 'POST',   '/api/permissions/',               NOW(3), NOW(3)),
  (10, '权限详情',                 'sys:permission:view',        'api', 'GET',    '/api/permissions/:id',            NOW(3), NOW(3)),
  (11, '更新权限',                 'sys:permission:update',      'api', 'PUT',    '/api/permissions/:id',            NOW(3), NOW(3)),
  (12, '删除权限',                 'sys:permission:delete',      'api', 'DELETE', '/api/permissions/:id',            NOW(3), NOW(3)),
  (13, '菜单列表',                 'sys:menu:list',              'api', 'GET',    '/api/menus/',                     NOW(3), NOW(3)),
  (14, '创建菜单',                 'sys:menu:create',            'api', 'POST',   '/api/menus/',                     NOW(3), NOW(3)),
  (15, '菜单详情',                 'sys:menu:view',              'api', 'GET',    '/api/menus/:id',                  NOW(3), NOW(3)),
  (16, '更新菜单',                 'sys:menu:update',            'api', 'PUT',    '/api/menus/:id',                  NOW(3), NOW(3)),
  (17, '删除菜单',                 'sys:menu:delete',            'api', 'DELETE', '/api/menus/:id',                  NOW(3), NOW(3)),
  (18, '前端菜单树',               'sys:menu:frontend',          'api', 'GET',    '/api/menus/frontend/tree',        NOW(3), NOW(3)),
  (19, '菜单权限查看',             'sys:menu_permission:view',   'api', 'GET',    '/api/menus/:id/permissions',      NOW(3), NOW(3)),
  (20, '菜单权限绑定',             'sys:menu_permission:bind',   'api', 'PUT',    '/api/menus/:id/permissions',      NOW(3), NOW(3)),
  (21, '用户列表',                 'sys:user:list',              'api', 'GET',    '/api/users/',                     NOW(3), NOW(3)),
  (22, '创建用户',                 'sys:user:create',            'api', 'POST',   '/api/users/',                     NOW(3), NOW(3)),
  (23, '用户详情',                 'sys:user:view',              'api', 'GET',    '/api/users/:id',                  NOW(3), NOW(3)),
  (24, '更新用户',                 'sys:user:update',            'api', 'PUT',    '/api/users/:id',                  NOW(3), NOW(3)),
  (25, '删除用户',                 'sys:user:delete',            'api', 'DELETE', '/api/users/:id',                  NOW(3), NOW(3)),
  (26, '用户资料查看',             'sys:user:profile',           'api', 'GET',    '/api/users/profile',              NOW(3), NOW(3)),
  (27, '用户角色查看',             'sys:user_role:view',         'api', 'GET',    '/api/users/role/:id',             NOW(3), NOW(3)),
  (28, '用户角色绑定',             'sys:user_role:bind',         'api', 'POST',   '/api/users/role/bind',            NOW(3), NOW(3)),
  (29, '用户角色新增',             'sys:user_role:add',          'api', 'POST',   '/api/users/role/add',             NOW(3), NOW(3)),
  (30, '用户角色移除',             'sys:user_role:remove',       'api', 'POST',   '/api/users/role/remove',          NOW(3), NOW(3)),
  (31, '登录日志查看',             'sys:audit:login_log:list',   'api', 'GET',    '/api/audit/login-logs',           NOW(3), NOW(3)),
  (32, '系统消息列表',             'sys:message:list',           'api', 'GET',    '/api/messages/',                  NOW(3), NOW(3)),
  (33, '系统消息已读',             'sys:message:read',           'api', 'POST',   '/api/messages/:id/read',          NOW(3), NOW(3)),
  (34, '系统消息全部已读',         'sys:message:read_all',       'api', 'POST',   '/api/messages/read-all',          NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE `updated_at` = VALUES(`updated_at`);

-- -------------------------
-- 13) 初始化菜单
-- -------------------------
INSERT INTO `sys_menus` (`id`, `title`, `path`, `icon`, `parent_id`, `component`, `order_num`, `created_at`, `updated_at`)
VALUES
  (1, '系统管理', '/system',               'setting', 0, 'Layout',                        10, NOW(3), NOW(3)),
  (8, '用户管理', '/system/users',         'user',    1, 'system/user/index',             5,  NOW(3), NOW(3)),
  (2, '角色管理', '/system/roles',         'team',    1, 'system/role/index',             10, NOW(3), NOW(3)),
  (3, '权限管理', '/system/permissions',   'lock',    1, 'system/permission/index',       20, NOW(3), NOW(3)),
  (4, '菜单管理', '/system/menus',         'menu',    1, 'system/menu/index',             30, NOW(3), NOW(3)),
  (5, '用户角色', '/system/user-roles',    'user',    1, 'system/user-role/index',        40, NOW(3), NOW(3)),
  (6, '审计中心', '/audit',                'audit',   0, 'Layout',                        20, NOW(3), NOW(3)),
  (7, '登录日志', '/audit/login-logs',     'history', 6, 'audit/login-log/index',         10, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE `updated_at` = VALUES(`updated_at`);

-- -------------------------
-- 14) 角色-权限初始化
-- admin：所有权限
-- security_auditor：只读类权限 + 菜单树
-- operator：基础运营权限（菜单树 + 用户角色查看）
-- -------------------------
DELETE FROM `sys_role_permissions` WHERE `role_id` IN (1, 2, 3);

-- admin -> 全部权限
INSERT INTO `sys_role_permissions` (`role_id`, `permission_id`)
SELECT 1 AS role_id, p.id AS permission_id FROM `sys_permissions` p;

-- security_auditor -> 只读权限集合
INSERT INTO `sys_role_permissions` (`role_id`, `permission_id`) VALUES
  (2, 1), (2, 3), (2, 6), (2, 8), (2, 10), (2, 13), (2, 15), (2, 18), (2, 19), (2, 21), (2, 23), (2, 26), (2, 27), (2, 31), (2, 32), (2, 33), (2, 34);

-- operator -> 基础权限集合
INSERT INTO `sys_role_permissions` (`role_id`, `permission_id`) VALUES
  (3, 13), (3, 15), (3, 18), (3, 21), (3, 23), (3, 26), (3, 27), (3, 32), (3, 33), (3, 34);

-- -------------------------
-- 15) 菜单-权限初始化
-- 说明：用于“按权限返回前端菜单树”
-- -------------------------
DELETE FROM `sys_menu_permissions` WHERE `menu_id` IN (1, 2, 3, 4, 5, 6, 7, 8);

INSERT INTO `sys_menu_permissions` (`menu_id`, `permission_id`) VALUES
  -- 根菜单：系统管理（任意管理类权限即可看到）
  (1, 13), (1, 21), (1, 27), (1, 18),
  -- 用户管理
  (8, 21), (8, 22), (8, 23), (8, 24), (8, 25),
  -- 角色管理
  (2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6), (2, 7),
  -- 权限管理
  (3, 8), (3, 9), (3, 10), (3, 11), (3, 12),
  -- 菜单管理
  (4, 13), (4, 14), (4, 15), (4, 16), (4, 17), (4, 19), (4, 20),
  -- 用户角色
  (5, 27), (5, 28), (5, 29), (5, 30),
  -- 审计中心
  (6, 31),
  -- 登录日志
  (7, 31);

-- -------------------------
-- 16) 初始化管理员账号
-- 默认账号：admin
-- 默认密码：Admin@123456
-- -------------------------
INSERT INTO `sys_users` (`id`, `username`, `password_hash`, `email`, `avatar`, `status`, `created_at`, `updated_at`)
VALUES
  (1, 'admin', '$2a$10$AwkXKvSji6Xb.rZdjQGkR.60r7eZZaoGYJo50gUoIFLiOpH/W1pA.', 'admin@example.com', '', 1, NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE
  `password_hash` = VALUES(`password_hash`),
  `email` = VALUES(`email`),
  `avatar` = VALUES(`avatar`),
  `status` = VALUES(`status`),
  `updated_at` = VALUES(`updated_at`);

-- admin 用户绑定 admin 角色
INSERT INTO `sys_user_roles` (`user_id`, `role_id`)
VALUES (1, 1)
ON DUPLICATE KEY UPDATE `role_id` = VALUES(`role_id`);

SET FOREIGN_KEY_CHECKS = 1;

-- =========================================================
-- 使用说明：
-- 1. 执行：mysql -u用户名 -p 数据库名 < 001_rbac_schema_and_seed.sql
-- 2. 首次登录：
--    用户名：admin
--    密码：Admin@123456
-- 3. 生产环境请立即修改默认管理员密码。
-- =========================================================
