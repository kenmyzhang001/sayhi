-- PostgreSQL 版本数据库表结构

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

-- ============================================
-- 位置值配置表
-- ============================================
CREATE TABLE IF NOT EXISTS position_values (
  id BIGSERIAL PRIMARY KEY,
  position VARCHAR(10) NOT NULL,
  value VARCHAR(500) NOT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_position_values_position ON position_values(position);
CREATE INDEX idx_position_values_position_sort ON position_values(position, sort_order);

-- ============================================
-- 模板配置表
-- ============================================
CREATE TABLE IF NOT EXISTS templates (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  template TEXT NOT NULL,
  encoding VARCHAR(20) NOT NULL DEFAULT 'Unicode',
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_templates_user_id ON templates(user_id);

-- ============================================
-- 生成历史记录表
-- ============================================
CREATE TABLE IF NOT EXISTS generate_history (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  template TEXT NOT NULL,
  encoding VARCHAR(20) NOT NULL,
  generate_mode VARCHAR(20) NOT NULL,
  total_count INTEGER NOT NULL,
  exceeded_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_generate_history_user_id ON generate_history(user_id);
CREATE INDEX idx_generate_history_created_at ON generate_history(created_at);

-- ============================================
-- 触发器：自动更新 updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_position_values_updated_at BEFORE UPDATE ON position_values
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_templates_updated_at BEFORE UPDATE ON templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 初始化默认数据
-- ============================================
INSERT INTO users (username, password) VALUES 
('admin', '0192023a7bbd73250516f069df18b500'),
('user', '6ad14ba9986e3615423dfca256d04e3f')
ON CONFLICT (username) DO NOTHING;

INSERT INTO position_values (position, value, sort_order) VALUES
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
ON CONFLICT DO NOTHING;

