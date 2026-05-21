package model

import "time"

// Address 收货地址模型，映射 addresses 表
type Address struct {
	ID        uint      `json:"id" gorm:"primaryKey"`               // 地址ID，主键自增
	UserID    uint      `json:"user_id" gorm:"index;not null"`      // 所属用户ID，普通索引
	Name      string    `json:"name" gorm:"size:50;not null"`       // 收件人姓名
	Phone     string    `json:"phone" gorm:"size:20;not null"`      // 收件人手机号
	Address   string    `json:"address" gorm:"size:255;not null"`   // 详细地址
	Tag       string    `json:"tag" gorm:"size:20"`                 // 地址标签（如：家/公司）
	IsDefault bool      `json:"is_default" gorm:"default:false"`    // 是否为默认地址
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
}

// TableName 返回 addresses 表名
func (Address) TableName() string {
	return "addresses"
}
