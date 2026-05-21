package model

import "time"

// Cart 购物车模型，映射 carts 表
type Cart struct {
	ID         uint      `json:"id" gorm:"primaryKey"`                       // 购物车条目ID
	UserID     uint      `json:"user_id" gorm:"index:idx_user_merchant;not null"` // 用户ID，复合索引
	MerchantID uint      `json:"merchant_id" gorm:"index:idx_user_merchant;not null"` // 商家ID，复合索引
	FoodID     uint      `json:"food_id" gorm:"not null"`                    // 商品ID
	Quantity   int       `json:"quantity" gorm:"default:1"`                  // 数量
	Food       Food      `json:"food" gorm:"foreignKey:FoodID"`              // 关联的商品信息（Preload 使用）
	Merchant   Merchant  `json:"merchant" gorm:"foreignKey:MerchantID"`      // 关联的商家信息（Preload 使用）
	CreatedAt  time.Time `json:"created_at"`                                 // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`                                 // 更新时间
}

// TableName 返回 carts 表名
func (Cart) TableName() string {
	return "carts"
}
