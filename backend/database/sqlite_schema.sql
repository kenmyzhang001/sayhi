-- SQLite 版本数据库表结构

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

-- ============================================
-- 位置值配置表
-- ============================================
CREATE TABLE IF NOT EXISTS position_values (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  position TEXT NOT NULL,
  value TEXT NOT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_position_values_position ON position_values(position);
CREATE INDEX idx_position_values_position_sort ON position_values(position, sort_order);

-- ============================================
-- 模板配置表
-- ============================================
CREATE TABLE IF NOT EXISTS templates (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  template TEXT NOT NULL,
  encoding TEXT NOT NULL DEFAULT 'Unicode',
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_user_id ON templates(user_id);

-- ============================================
-- 生成历史记录表
-- ============================================
CREATE TABLE IF NOT EXISTS generate_history (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  template TEXT NOT NULL,
  encoding TEXT NOT NULL,
  generate_mode TEXT NOT NULL,
  total_count INTEGER NOT NULL,
  exceeded_count INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_generate_history_user_id ON generate_history(user_id);
CREATE INDEX idx_generate_history_created_at ON generate_history(created_at);

-- ============================================
-- 触发器：自动更新 updated_at
-- ============================================
CREATE TRIGGER IF NOT EXISTS update_users_updated_at 
  AFTER UPDATE ON users
  FOR EACH ROW
BEGIN
  UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_position_values_updated_at 
  AFTER UPDATE ON position_values
  FOR EACH ROW
BEGIN
  UPDATE position_values SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_templates_updated_at 
  AFTER UPDATE ON templates
  FOR EACH ROW
BEGIN
  UPDATE templates SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- ============================================
-- 初始化默认数据
-- ============================================
INSERT OR IGNORE INTO users (username, password) VALUES 
('admin', '0192023a7bbd73250516f069df18b500'),
('user', '6ad14ba9986e3615423dfca256d04e3f');

INSERT OR IGNORE INTO position_values (position, value, sort_order) VALUES
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
('d', '10', 8);

