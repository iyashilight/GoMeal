package model

import "time"

// Payment 支付记录模型，映射 payments 表
// 每条支付记录对应一个订单，通过 order_id 唯一索引保证幂等
type Payment struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	OrderID   uint       `json:"order_id" gorm:"uniqueIndex;not null"`     // 订单ID（唯一约束，防重复支付）
	UserID    uint       `json:"user_id" gorm:"index;not null"`            // 支付用户ID
	Amount    float64    `json:"amount" gorm:"not null"`                   // 支付金额
	Method    string     `json:"method" gorm:"size:20"`                    // 支付方式：wechat/alipay/mock
	TradeNo   string     `json:"trade_no" gorm:"uniqueIndex;size:64"`      // 支付网关流水号
	Status    int        `json:"status" gorm:"default:0"`                  // 0=待支付, 1=已支付, 2=已退款
	PaidAt    *time.Time `json:"paid_at"`                                  // 支付时间
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}
