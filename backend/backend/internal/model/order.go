package model

import "time"

// 订单状态常量
const (
	OrderStatusPending   = 0 // 待支付
	OrderStatusPaid      = 1 // 已支付
	OrderStatusAccepted  = 2 // 已接单
	OrderStatusDelivery  = 3 // 配送中
	OrderStatusCompleted = 4 // 已完成
	OrderStatusCancelled = 5 // 已取消
	OrderStatusRefunded  = 6 // 已退款
)

// Order 订单模型，映射 orders 表
type Order struct {
	ID          uint   `json:"id" gorm:"primaryKey"`                 // 订单ID，主键自增
	OrderNo     string `json:"order_no" gorm:"uniqueIndex;size:32;not null"` // 订单号，唯一索引（格式：yyyymmddhhmmss + 6位随机数）
	UserID      uint        `json:"user_id" gorm:"index;not null"`          // 下单用户ID
	MerchantID  uint        `json:"merchant_id" gorm:"index;not null"`      // 商家ID
	AddressID   uint        `json:"address_id" gorm:"not null"`             // 收货地址ID
	TotalAmount float64     `json:"total_amount" gorm:"not null"`           // 订单总金额（含配送费）
	DeliveryFee float64     `json:"delivery_fee"`                           // 配送费
	Remark      string      `json:"remark" gorm:"size:500"`                 // 订单备注
	Status      int         `json:"status" gorm:"default:0"`                // 订单状态，默认待支付
	PaidAt      *time.Time  `json:"paid_at"`                                // 支付时间（指针类型，允许为空）
	DeliveredAt *time.Time  `json:"delivered_at"`                           // 送达时间
	CompletedAt *time.Time  `json:"completed_at"`                           // 完成时间
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`        // 订单商品列表（HasMany）
	Merchant    Merchant    `json:"merchant" gorm:"foreignKey:MerchantID"`  // 商家信息（BelongsTo）
	Address     Address     `json:"address" gorm:"foreignKey:AddressID"`    // 收货地址（BelongsTo）
	CreatedAt   time.Time   `json:"created_at"`                             // 创建时间
	UpdatedAt   time.Time   `json:"updated_at"`                             // 更新时间
}

// TableName 返回 orders 表名
func (Order) TableName() string {
	return "orders"
}

// OrderItem 订单商品明细模型，映射 order_items 表
type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`              // 明细ID
	OrderID   uint      `json:"order_id" gorm:"index;not null"`    // 所属订单ID
	FoodID    uint      `json:"food_id" gorm:"not null"`           // 商品ID
	FoodName  string    `json:"food_name" gorm:"size:100;not null"` // 商品名称（下单时快照）
	FoodImage string    `json:"food_image" gorm:"size:255"`        // 商品图片（下单时快照）
	Price     float64   `json:"price" gorm:"not null"`             // 下单时单价
	Quantity  int       `json:"quantity" gorm:"not null"`          // 购买数量
	CreatedAt time.Time `json:"created_at"`                        // 创建时间
}

// TableName 返回 order_items 表名
func (OrderItem) TableName() string {
	return "order_items"
}
