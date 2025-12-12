# 数据库脚本说明

本目录包含数据库建表和迁移脚本。

## 文件说明

### 主文件
- `schema.sql` - MySQL 完整建表脚本（包含所有表）
- `postgresql_schema.sql` - PostgreSQL 版本建表脚本
- `sqlite_schema.sql` - SQLite 版本建表脚本
- `init_data.sql` - 初始化默认数据

### 迁移脚本
- `migrations/001_initial_schema.sql` - 初始表结构（用户表、位置值表）
- `migrations/002_add_templates_table.sql` - 添加模板表
- `migrations/003_add_generate_history_table.sql` - 添加历史记录表

## 使用方法

### MySQL

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE sayhi CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行建表脚本
mysql -u root -p sayhi < schema.sql

# 执行初始化数据
mysql -u root -p sayhi < init_data.sql
```

### PostgreSQL

```bash
# 创建数据库
createdb sayhi

# 执行建表脚本
psql -U postgres -d sayhi -f postgresql_schema.sql
```

### SQLite

```bash
# 执行建表脚本
sqlite3 sayhi.db < sqlite_schema.sql
```

## 表结构说明

### users - 用户表
- `id` - 用户ID（主键）
- `username` - 用户名（唯一）
- `password` - 密码（MD5加密）
- `created_at` - 创建时间
- `updated_at` - 更新时间

### position_values - 位置值配置表
- `id` - ID（主键）
- `position` - 位置标识（a, b, c, d）
- `value` - 位置值
- `sort_order` - 排序顺序
- `created_at` - 创建时间
- `updated_at` - 更新时间

### templates - 模板配置表（可选）
- `id` - 模板ID（主键）
- `name` - 模板名称
- `template` - 模板内容
- `encoding` - 字符编码
- `user_id` - 创建用户ID（外键）
- `created_at` - 创建时间
- `updated_at` - 更新时间

### generate_history - 生成历史记录表（可选）
- `id` - 记录ID（主键）
- `user_id` - 用户ID（外键）
- `template` - 使用的模板
- `encoding` - 字符编码
- `generate_mode` - 生成方式
- `total_count` - 生成总数
- `exceeded_count` - 超出数量
- `created_at` - 创建时间

## 默认账号

系统初始化后会创建以下默认账号：

- **管理员**: `admin` / `admin123`
- **普通用户**: `user` / `user123`

密码已使用 MD5 加密存储。

## 迁移说明

如果需要按顺序执行迁移脚本：

```bash
# MySQL
mysql -u root -p sayhi < migrations/001_initial_schema.sql
mysql -u root -p sayhi < migrations/002_add_templates_table.sql
mysql -u root -p sayhi < migrations/003_add_generate_history_table.sql
mysql -u root -p sayhi < init_data.sql
```

## 注意事项

1. **字符集**: MySQL 使用 `utf8mb4` 以支持完整的 Unicode 字符
2. **外键约束**: 确保外键约束正确设置，删除用户时会级联删除相关数据
3. **索引**: 已为常用查询字段创建索引，提升查询性能
4. **密码加密**: 当前使用 MD5，生产环境建议使用 bcrypt 或 argon2

## 后续集成

要将代码从内存存储迁移到数据库，需要：

1. 安装数据库驱动（如 `gorm` 或 `database/sql`）
2. 创建数据库连接
3. 修改 `auth_service.go` 和 `position_service.go` 使用数据库操作
4. 更新 `go.mod` 添加数据库驱动依赖

示例依赖：
```go
// MySQL
github.com/go-sql-driver/mysql

// PostgreSQL
github.com/lib/pq

// SQLite
github.com/mattn/go-sqlite3

// ORM (可选)
gorm.io/gorm
gorm.io/driver/mysql
```

