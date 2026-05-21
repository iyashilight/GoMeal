// Package handler 实现 HTTP 接口层，接收请求并调用 Service 处理，返回统一响应格式
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// health 健康检查接口
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
