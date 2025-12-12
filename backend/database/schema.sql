-- 短信模板生成系统数据库表结构
-- 数据库：MySQL 5.7+ / PostgreSQL 10+ / SQLite 3

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` VARCHAR(50) NOT NULL COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码（加密后）',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 位置值配置表
-- ============================================
CREATE TABLE IF NOT EXISTS `position_values` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `position` VARCHAR(10) NOT NULL COMMENT '位置标识（a, b, c, d）',
  `value` VARCHAR(500) NOT NULL COMMENT '位置值',
  `sort_order` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_position` (`position`),
  KEY `idx_position_sort` (`position`, `sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='位置值配置表';

-- ============================================
-- 模板配置表（可选，用于保存常用模板）
-- ============================================
CREATE TABLE IF NOT EXISTS `templates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '模板ID',
  `name` VARCHAR(100) NOT NULL COMMENT '模板名称',
  `template` TEXT NOT NULL COMMENT '模板内容',
  `encoding` VARCHAR(20) NOT NULL DEFAULT 'Unicode' COMMENT '字符编码',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '创建用户ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  CONSTRAINT `fk_templates_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模板配置表';

-- ============================================
-- 生成历史记录表（可选，用于记录生成历史）
-- ============================================
CREATE TABLE IF NOT EXISTS `generate_history` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `template` TEXT NOT NULL COMMENT '使用的模板',
  `encoding` VARCHAR(20) NOT NULL COMMENT '字符编码',
  `generate_mode` VARCHAR(20) NOT NULL COMMENT '生成方式（sequential/random）',
  `total_count` INT UNSIGNED NOT NULL COMMENT '生成总数',
  `exceeded_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '超出数量',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_history_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='生成历史记录表';

-- ============================================
-- 初始化默认数据
-- ============================================

-- 插入默认管理员账号（密码：admin123，MD5加密后）
INSERT INTO `users` (`username`, `password`) VALUES 
('admin', '0192023a7bbd73250516f069df18b500'),
('user', '6ad14ba9986e3615423dfca256d04e3f')
ON DUPLICATE KEY UPDATE `username`=`username`;

-- 插入示例位置值
INSERT INTO `position_values` (`position`, `value`, `sort_order`) VALUES
('a', '1', 1),
('b', 'baidu.com', 1),
('c', '2', 1),
('d', '3', 1),
('d', '4', 2),
('d', '5', 3),
('d', '6', 4),
('d', '7', 5),
('d', '8', 6),
('d', '9', 7),
('d', '10', 8)
ON DUPLICATE KEY UPDATE `value`=`value`;

