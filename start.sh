#!/bin/bash

# 启动脚本 - 同时启动前后端

echo "启动短信模板生成系统..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未安装 Go，请先安装 Go 1.21+"
    exit 1
fi

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo "错误: 未安装 Node.js，请先安装 Node.js"
    exit 1
fi

# 启动后端
echo "启动后端服务..."
cd backend
go mod download
go run main.go &
BACKEND_PID=$!
cd ..

# 等待后端启动
sleep 2

# 启动前端
echo "启动前端应用..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi
npm run dev &
FRONTEND_PID=$!
cd ..

echo "后端服务运行在: http://localhost:8080"
echo "前端应用运行在: http://localhost:3000"
echo ""
echo "按 Ctrl+C 停止服务"

# 等待用户中断
trap "kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT TERM
wait

