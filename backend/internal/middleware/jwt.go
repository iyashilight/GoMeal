// Package middleware 提供 HTTP 中间件：JWT 鉴权、CORS、访问日志、限流、RBAC
package middleware

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/utils"
)

// JWTAuth JWT 鉴权中间件
// 从 Authorization Header 解析 Bearer Token，验证后将用户信息注入 Gin Context
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"code": 401, "message": "未登录"})
			slog.Warn("missing authorization header", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"code": 401, "message": "认证格式错误"})
			slog.Warn("invalid authorization format", "path", c.Request.URL.Path, "ip", c.ClientIP())
			c.Abort()
			return
		}
		claims, err := utils.ParseToken(parts[1], secret)
		if err != nil {
			c.JSON(401, gin.H{"code": 401, "message": "token无效或已过期"})
			slog.Warn("invalid token", "path", c.Request.URL.Path, "ip", c.ClientIP(), "error", err)
			c.Abort()
			return
		}
		// 将用户信息注入上下文，后续 Handler/Service 通过 GetUserID 等函数获取
		c.Set("userID", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("userType", claims.UserType)
		c.Next()
	}
}

// GetUserID 从 Gin Context 中提取当前登录用户 ID
func GetUserID(c *gin.Context) uint {
	id, _ := c.Get("userID")
	return id.(uint)
}
