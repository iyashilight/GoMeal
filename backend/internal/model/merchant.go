package model

import "time"

// Merchant 商家模型，映射 merchants 表
type Merchant struct {
	ID          uint       `json:"id" gorm:"primaryKey"`            // 商家ID，主键自增
	UserID      uint       `json:"user_id" gorm:"index;not null"`   // 关联的用户ID（店主）
	Name        string     `json:"name" gorm:"size:100;not null"`   // 商家名称
	Logo        string     `json:"logo" gorm:"size:255"`            // 商家 Logo URL
	Notice      string     `json:"notice" gorm:"size:500"`          // 商家公告
	Phone       string     `json:"phone" gorm:"size:20"`            // 商家联系电话
	Address     string     `json:"address" gorm:"size:255"`         // 商家地址
	MinPrice    float64    `json:"min_price" gorm:"default:0"`      // 起送价
	DeliveryFee float64    `json:"delivery_fee" gorm:"default:0"`   // 配送费
	Status      int        `json:"status" gorm:"default:1"`         // 营业状态：1-营业中，0-休息中
	Rating      float64    `json:"rating" gorm:"default:5.0"`       // 评分（满分 5.0）
	Sales       int        `json:"sales" gorm:"default:0"`          // 月销量
	Categories  []Category `json:"categories" gorm:"foreignKey:MerchantID"` // 关联的分类列表（HasMany）
	CreatedAt   time.Time  `json:"created_at"`                      // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`                      // 更新时间
}

// TableName 返回 merchants 表名
func (Merchant) TableName() string {
	return "merchants"
}

// Category 分类模型，映射 categories 表
type Category struct {
	ID         uint      `json:"id" gorm:"primaryKey"`                      // 分类ID
	MerchantID uint      `json:"merchant_id" gorm:"index;not null"`         // 所属商家ID
	Name       string    `json:"name" gorm:"size:50;not null"`              // 分类名称
	SortOrder  int       `json:"sort_order" gorm:"default:0"`               // 排序序号（升序）
	CreatedAt  time.Time `json:"created_at"`                                // 创建时间
	Foods      []Food    `json:"foods" gorm:"foreignKey:CategoryID"`        // 分类下的商品列表（HasMany）
}

// TableName 返回 categories 表名
func (Category) TableName() string {
	return "categories"
}

// Food 商品（食品）模型，映射 foods 表
type Food struct {
	ID          uint      `json:"id" gorm:"primaryKey"`               // 商品ID
	MerchantID  uint      `json:"merchant_id" gorm:"index;not null"`  // 所属商家ID
	CategoryID  uint      `json:"category_id" gorm:"index;not null"`  // 所属分类ID
	Name        string    `json:"name" gorm:"size:100;not null"`      // 商品名称
	Price       float64   `json:"price" gorm:"not null"`              // 当前售价
	OldPrice    float64   `json:"old_price"`                          // 原价（用于展示折扣）
	Stock       int       `json:"stock" gorm:"default:9999"`          // 库存数量
	Sales       int       `json:"sales" gorm:"default:0"`             // 累计销量
	Status      int       `json:"status" gorm:"default:1"`            // 状态：1-上架，0-下架
	IsRecommend int       `json:"is_recommend" gorm:"default:0"`      // 是否推荐：1-推荐，0-普通
	Description string    `json:"description" gorm:"size:500"`        // 商品描述
	Image       string    `json:"image" gorm:"size:255"`              // 商品图片 URL
	CreatedAt   time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`                         // 更新时间
}

// TableName 返回 foods 表名
func (Food) TableName() string {
	return "foods"
}
