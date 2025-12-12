-- 初始化数据脚本
-- 说明: 插入系统默认数据

-- 插入默认用户账号
-- 密码说明：
-- admin123 -> MD5: 0192023a7bbd73250516f069df18b500
-- user123  -> MD5: 6ad14ba9986e3615423dfca256d04e3f

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

