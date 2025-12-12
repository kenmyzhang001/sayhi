# 短信模板生成系统 - 后端

## 技术栈
- Go 1.21+
- Gin Web框架
- JWT认证
- CORS支持

## 快速开始

### 1. 安装依赖
```bash
cd backend
go mod download
```

### 2. 运行服务
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

## API文档

### 认证相关

#### 1. 用户登录
**POST** `/api/auth/login`

请求体：
```json
{
  "username": "admin",
  "password": "admin123"
}
```

响应：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin",
  "message": "登录成功"
}
```

#### 2. 用户注册
**POST** `/api/auth/register`

请求体：
```json
{
  "username": "newuser",
  "password": "password123"
}
```

响应：
```json
{
  "message": "注册成功"
}
```

#### 3. 获取用户信息
**GET** `/api/auth/user`

需要认证：是（Bearer Token）

响应：
```json
{
  "id": 1,
  "username": "admin"
}
```

### 模板生成（需要认证）

#### 1. 生成短信内容
**POST** `/api/template/generate`

需要认证：是（Bearer Token）

请求头：
```
Authorization: Bearer <token>
```

请求体：
```json
{
  "template": "(1)(baidu.com)(2)(3-10)",
  "encoding": "Unicode",
  "generateMode": "sequential",
  "positions": {
    "a": ["1"],
    "b": ["baidu.com"],
    "c": ["2"],
    "d": ["3", "4", "5", "6", "7", "8", "9", "10"]
  }
}
```

响应：
```json
{
  "results": [
    {
      "content": "1 baidu.com 2 3",
      "charCount": 15,
      "isExceeded": false,
      "exceededChars": 0
    }
  ],
  "totalCount": 8,
  "exceededCount": 0
}
```

### 位置值管理（需要认证）

#### 1. 获取所有位置值
**GET** `/api/positions`

需要认证：是

#### 2. 获取指定位置值
**GET** `/api/positions/:position`

需要认证：是

#### 3. 添加位置值
**POST** `/api/positions`
需要认证：是

```json
{
  "position": "a",
  "value": "新值"
}
```

#### 4. 设置位置的所有值
**PUT** `/api/positions/:position`
需要认证：是

```json
{
  "values": ["值1", "值2", "值3"]
}
```

#### 5. 删除位置值
**DELETE** `/api/positions/:position?value=要删除的值`
需要认证：是

## 默认账号

系统初始化时会创建以下测试账号：

- **管理员**: `admin` / `admin123`
- **普通用户**: `user` / `user123`

## 认证说明

- 所有API接口（除登录和注册外）都需要在请求头中携带JWT token
- Token格式：`Authorization: Bearer <token>`
- Token有效期：24小时
- 密码使用MD5加密（生产环境建议使用bcrypt）

## 项目结构
```
backend/
├── main.go              # 主程序入口
├── models/              # 数据模型
│   ├── models.go        # 模板相关模型
│   └── user.go          # 用户模型
├── handlers/            # 请求处理器
│   ├── auth_handler.go  # 认证处理器
│   ├── template_handler.go
│   └── position_handler.go
├── services/            # 业务逻辑
│   ├── auth_service.go  # 认证服务
│   ├── generator.go
│   └── position_service.go
├── middleware/          # 中间件
│   └── auth.go          # 认证中间件
├── utils/               # 工具函数
│   ├── jwt.go           # JWT工具
│   ├── parser.go
│   └── encoder.go
└── go.mod              # 依赖管理
```
