# 配置文件说明

## 配置文件位置

配置文件位于 `backend/config/` 目录下：

- `config.go` - Go 配置结构体和加载逻辑
- `config.yaml` - YAML 格式配置文件（可选）
- `config.example.yaml` - 配置文件示例
- `.env.example` - 环境变量配置示例

## 配置方式

系统支持两种配置方式（按优先级排序）：

1. **环境变量**（最高优先级）
2. **代码默认值**（最低优先级）

### 使用环境变量

创建 `.env` 文件（在 `backend/` 目录下）：

```bash
cp .env.example .env
# 编辑 .env 文件
```

或在启动时设置环境变量：

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=sayhi
go run main.go
```

### 使用代码配置

直接修改 `config/config.go` 中的默认值。

## 配置项说明

### 服务器配置

| 环境变量 | 配置项 | 默认值 | 说明 |
|---------|--------|--------|------|
| `SERVER_HOST` | `server.host` | `0.0.0.0` | 服务器监听地址 |
| `SERVER_PORT` | `server.port` | `8080` | 服务器端口 |
| `GIN_MODE` | `server.mode` | `debug` | 运行模式 |

### 数据库配置

| 环境变量 | 配置项 | 默认值 | 说明 |
|---------|--------|--------|------|
| `DB_TYPE` | `database.type` | `mysql` | 数据库类型 |
| `DB_HOST` | `database.host` | `localhost` | 数据库主机 |
| `DB_PORT` | `database.port` | `3306` | 数据库端口 |
| `DB_USER` | `database.user` | `root` | 数据库用户名 |
| `DB_PASSWORD` | `database.password` | `` | 数据库密码 |
| `DB_NAME` | `database.database` | `sayhi` | 数据库名称 |
| `DB_CHARSET` | `database.charset` | `utf8mb4` | 字符集 |
| `DB_MAX_IDLE` | `database.max_idle` | `10` | 最大空闲连接数 |
| `DB_MAX_OPEN` | `database.max_open` | `100` | 最大打开连接数 |
| `DB_DSN` | `database.dsn` | `` | 完整连接字符串（可选） |

### JWT配置

| 环境变量 | 配置项 | 默认值 | 说明 |
|---------|--------|--------|------|
| `JWT_SECRET` | `jwt.secret` | `sayhi-secret-key...` | JWT密钥 |
| `JWT_EXPIRE_TIME` | `jwt.expire_time` | `24` | Token过期时间（小时） |

## 使用示例

### 在代码中使用配置

```go
import "sayhi/backend/config"

func main() {
    // 加载配置
    cfg := config.LoadConfig()
    
    // 使用配置
    fmt.Println("Server Port:", cfg.Server.Port)
    fmt.Println("Database:", cfg.Database.Database)
    fmt.Println("JWT Secret:", cfg.JWT.Secret)
    
    // 获取数据库连接字符串
    dsn := cfg.Database.GetDSN()
}
```

### 不同数据库的配置示例

#### MySQL
```bash
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password123
DB_NAME=sayhi
```

#### PostgreSQL
```bash
DB_TYPE=postgresql
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password123
DB_NAME=sayhi
```

#### SQLite
```bash
DB_TYPE=sqlite
DB_NAME=sayhi.db
```

## 安全建议

1. **生产环境**：
   - 修改 `JWT_SECRET` 为强随机字符串
   - 使用环境变量而不是硬编码
   - 不要将 `.env` 文件提交到版本控制

2. **数据库密码**：
   - 使用强密码
   - 定期更换密码
   - 限制数据库访问权限

3. **配置文件**：
   - 将 `.env.example` 提交到版本控制
   - 将 `.env` 添加到 `.gitignore`
   - 使用密钥管理服务（如 AWS Secrets Manager）

## 加载配置

在 `main.go` 中加载配置：

```go
import "sayhi/backend/config"

func main() {
    // 加载配置
    cfg := config.LoadConfig()
    
    // 设置 Gin 模式
    gin.SetMode(cfg.Server.Mode)
    
    // 使用配置启动服务器
    r.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}
```

