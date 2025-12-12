-- 迁移脚本 003: 添加生成历史记录表
-- 执行时间: 2024-01-03
-- 说明: 添加生成历史记录功能

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

