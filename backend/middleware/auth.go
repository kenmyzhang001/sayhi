package middleware

import (
	"net/http"
	"sayhi/backend/services"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证token",
			})
			c.Abort()
			return
		}

		// 提取token (Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "认证格式错误，应为: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]
		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user", user)
		c.Set("username", user.Username)
		c.Next()
	}
}

