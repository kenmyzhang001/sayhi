# 快速开始指南

## 环境要求

- Go 1.21 或更高版本
- Node.js 16+ 和 npm
- 现代浏览器（Chrome、Firefox、Safari、Edge）

## 安装步骤

### 方式一：使用启动脚本（推荐）

```bash
./start.sh
```

脚本会自动：
1. 检查环境
2. 安装依赖
3. 启动后端服务（端口 8080）
4. 启动前端应用（端口 3000）

### 方式二：手动启动

#### 1. 启动后端

```bash
cd backend
go mod download
go run main.go
```

后端将在 `http://localhost:8080` 启动

#### 2. 启动前端（新终端）

```bash
cd frontend
npm install
npm run dev
```

前端将在 `http://localhost:3000` 启动

## 使用说明

### 1. 登录系统

1. 访问 `http://localhost:3000`
2. 使用默认账号登录：
   - 管理员：`admin` / `admin123`
   - 普通用户：`user` / `user123`
3. 或点击"还没有账号？立即注册"创建新账号

### 2. 配置位置值

1. 登录后，点击顶部菜单的"位置配置"
2. 为 a、b、c、d 四个位置配置值
3. 可以添加、编辑、删除或批量设置值

### 3. 生成短信内容

1. 点击"模板生成"菜单
2. 输入模板，例如：`(1)(baidu.com)(2)(3-10)`
3. 选择字符编码（默认 Unicode）
4. 选择生成方式（顺序/随机）
5. 点击"生成内容"按钮
6. 查看生成结果，可以编辑超出字符限制的内容

### 4. 模板语法

- **固定值**：`(1)`、`(baidu.com)`、`(2)`
- **范围值**：`(3-10)` 表示从 3 到 10，共 8 种话术
- **位置对应**：第一个括号是位置 a，第二个是位置 b，以此类推

### 5. 示例

**模板**：`(1)(baidu.com)(2)(3-10)`

**顺序生成结果**：
```
1 baidu.com 2 3
1 baidu.com 2 4
1 baidu.com 2 5
...
1 baidu.com 2 10
```

**随机生成结果**：
```
2 baidu.com 1 3
1 4 baidu.com 2
...
```

## 认证说明

- 所有功能都需要先登录
- Token会自动保存在浏览器本地存储
- Token过期后需要重新登录
- 可以点击右上角用户名下拉菜单退出登录

## 常见问题

### 登录失败

- 检查用户名和密码是否正确
- 默认账号：`admin` / `admin123` 或 `user` / `user123`
- 可以注册新账号

### 后端启动失败

- 检查端口 8080 是否被占用
- 确保已安装 Go 1.21+
- 运行 `go mod download` 安装依赖

### 前端启动失败

- 检查端口 3000 是否被占用
- 确保已安装 Node.js 16+
- 运行 `npm install` 安装依赖
- 清除缓存：`rm -rf node_modules package-lock.json && npm install`

### API 请求失败

- 确保后端服务已启动
- 检查浏览器控制台的错误信息
- 确认 CORS 配置正确

## 开发说明

### 后端开发

- 代码结构：`models/` → `services/` → `handlers/` → `main.go`
- 添加新功能：在对应的目录下添加文件
- 测试：可以使用 Postman 或 curl 测试 API

### 前端开发

- 代码结构：`src/views/` 页面，`src/components/` 组件，`src/api/` API 调用
- 修改样式：编辑对应组件的 `<style>` 部分
- 添加功能：在对应组件中添加逻辑

## 下一步

- 查看 `README.md` 了解项目详情
- 查看 `需求文档-短信模板生成系统.md` 了解完整需求
- 查看 `backend/README.md` 和 `frontend/README.md` 了解技术细节

