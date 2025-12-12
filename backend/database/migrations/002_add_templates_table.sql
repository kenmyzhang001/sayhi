-- 迁移脚本 002: 添加模板配置表
-- 执行时间: 2024-01-02
-- 说明: 添加模板保存功能

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

