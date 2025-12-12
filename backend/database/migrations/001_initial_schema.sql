-- 迁移脚本 001: 初始数据库结构
-- 执行时间: 2024-01-01
-- 说明: 创建初始数据库表结构

-- 用户表
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

-- 位置值配置表
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

