# 短信模板生成系统 - 前端

## 技术栈
- Vue 3
- Vite
- Vue Router
- Element Plus
- Axios

## 快速开始

### 1. 安装依赖
```bash
cd frontend
npm install
```

### 2. 运行开发服务器
```bash
npm run dev
```

应用将在 `http://localhost:3000` 启动

### 3. 构建生产版本
```bash
npm run build
```

## 功能说明

### 1. 模板编辑页面
- 输入模板内容（支持括号位置和范围值）
- 选择字符编码（ASCII、Zawgyi、Unicode、其它）
- 选择生成方式（顺序/随机）
- 实时显示生成结果
- 支持编辑超出字符限制的内容

### 2. 位置配置页面
- 配置 a、b、c、d 四个位置的值
- 支持添加、编辑、删除位置值
- 支持批量设置位置值

## 项目结构
```
frontend/
├── index.html          # HTML入口
├── vite.config.js      # Vite配置
├── package.json        # 依赖管理
└── src/
    ├── main.js         # 应用入口
    ├── App.vue         # 根组件
    ├── api/            # API调用
    ├── components/     # 组件
    └── views/          # 页面视图
```

## API 配置

前端通过代理访问后端 API，代理配置在 `vite.config.js` 中：
- 开发环境：自动代理到 `http://localhost:8080`
- 生产环境：需要配置实际的后端地址

