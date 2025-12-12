package main

import (
	"fmt"
	"sayhi/backend/config"
	"sayhi/backend/handlers"
	"sayhi/backend/middleware"
	"sayhi/backend/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 初始化服务
	authService := services.NewAuthService()
	positionService := services.NewPositionService()
	speechService := services.NewSpeechService()

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(authService)
	templateHandler := handlers.NewTemplateHandler(speechService)
	positionHandler := handlers.NewPositionHandler(positionService)
	speechHandler := handlers.NewSpeechHandler(speechService)

	// 认证中间件
	authMiddleware := middleware.AuthMiddleware(authService)

	// 公开路由（不需要认证）
	public := r.Group("/api")
	{
		// 认证相关
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
	}

	// 需要认证的路由
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		// 用户信息
		api.GET("/auth/user", authHandler.GetUserInfo)

		// 模板生成
		api.POST("/template/generate", templateHandler.Generate)

		// 位置值管理
		api.GET("/positions", positionHandler.GetAllPositions)
		api.GET("/positions/:position", positionHandler.GetPositionValues)
		api.POST("/positions", positionHandler.AddPositionValue)
		api.PUT("/positions/:position", positionHandler.SetPositionValues)
		api.DELETE("/positions/:position", positionHandler.DeletePositionValue)

		// 话术组管理
		api.GET("/speech-groups", speechHandler.GetAllGroups)
		api.GET("/speech-groups/:id", speechHandler.GetGroup)
		api.POST("/speech-groups", speechHandler.CreateGroup)
		api.PUT("/speech-groups/:id", speechHandler.UpdateGroup)
		api.DELETE("/speech-groups/:id", speechHandler.DeleteGroup)
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("服务器启动在: http://%s\n", addr)
	r.Run(addr)
}
