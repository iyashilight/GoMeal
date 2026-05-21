// Package utils 提供通用工具函数和类型：错误码、响应格式、JWT、密码加密、校验、分页
package utils

// AppError 应用层自定义错误，包含业务错误码和用户可读消息
// Service 层返回 AppError，Handler 层通过 ErrorFromErr 统一转换为 HTTP 响应
type AppError struct {
	Code    int    // 业务错误码（见 errcode.go）
	Message string // 错误描述
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError 创建自定义业务错误
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
