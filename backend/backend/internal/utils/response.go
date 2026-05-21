package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 HTTP 响应格式
type Response struct {
	Code    int         `json:"code"`    // 业务状态码，0 表示成功
	Message string      `json:"message"` // 提示消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Error 返回指定错误码的失败响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorFromErr 根据错误类型返回响应
// 如果 err 是 AppError，使用其中的 Code 和 Message
// 否则视为参数校验错误（1001）
func ErrorFromErr(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(http.StatusOK, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}
	Error(c, ErrInvalidParams, err.Error())
}
