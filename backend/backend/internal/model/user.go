// Package model 定义数据库模型结构与 GORM 映射
package model

import "time"

const (
	RoleUser  = 0 // 普通用户
	RoleAdmin = 1 // 管理员
)

// User 用户模型，映射 users 表
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`               // 用户ID，主键自增
	Phone     string    `json:"phone" gorm:"uniqueIndex;not null;size:20"` // 手机号，唯一索引
	Password  string    `json:"-" gorm:"size:255;not null"`         // 密码（bcrypt 哈希），JSON 序列化时隐藏
	Nickname  string    `json:"nickname" gorm:"size:50"`            // 昵称
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
	Avatar    string    `json:"avatar" gorm:"size:255"`             // 头像 URL
	UserType  int       `json:"user_type" gorm:"default:0"`         // 用户类型：0-普通用户，1-管理员
}

// TableName 返回 users 表名
func (User) TableName() string {
	return "users"
}
