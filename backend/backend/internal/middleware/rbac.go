package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/utils"
)

// RequiredRole 基于角色的访问控制中间件
// 允许指定角色列表中的用户通过，其余用户返回 403
func RequiredRole(roles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("userType")
		if !exists {
			utils.Error(c, utils.ErrUnauthorized, "未登录")
			c.Abort()
			return
		}
		for _, role := range roles {
			if userType.(int) == role {
				c.Next()
				return
			}
		}
		utils.Error(c, utils.ErrAccessDenied, "无权限")
		c.Abort()
	}
}

