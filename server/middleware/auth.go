package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员权限中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			c.Abort()
			return
		}

		if role.(string) != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OwnerAuthMiddleware 资源所有者权限中间件
func OwnerAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			c.Abort()
			return
		}

		resourceUserID := c.GetUint("resourceUserID") // 在路由处理函数中设置
		role := c.GetString("role")

		// 管理员或资源所有者可以访问
		if role == "admin" || currentUserID.(uint) == resourceUserID {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "权限不足",
		})
		c.Abort()
	}
}
