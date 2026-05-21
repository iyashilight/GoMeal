package repository

import (
	"time"

	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// OrderRepository 订单仓储接口
type OrderRepository interface {
	Create(order *model.Order) error                                  // 创建订单（含订单商品）
	FindByUser(userID uint, status int, page, size int) ([]model.Order, error) // 按用户和状态分页查询
	FindByID(orderID uint) (*model.Order, error)                      // 按 ID 查询订单
	UpdateStatus(orderID uint, status int) error                      // 更新订单状态
	UpdatePaid(orderID uint) error                                    // 支付成功：更新状态+支付时间
	CountByUser(userID uint, status int) (int64, error)               // 统计用户订单数
}

// orderRepository 订单仓储的 GORM 实现
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储实例
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create 创建订单，包含关联的 OrderItem（通过 GORM 的关联自动插入）
func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

// FindByUser 分页查询用户订单，支持按状态筛选（status = -1 时查询所有）
// 按创建时间倒序排列，预加载商家和订单商品信息
func (r *orderRepository) FindByUser(userID uint, status int, page, size int) (
	[]model.Order, error) {
	var orders []model.Order
	query := r.db.Where("user_id = ?", userID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at desc").Limit(size).Offset((page - 1) * size).
		Preload("Merchant").Preload("Items").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// FindByID 按主键查询订单详情，预加载商家、地址、订单商品
func (r *orderRepository) FindByID(orderID uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Merchant").Preload("Address").Preload("Items").
		First(&order, orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatus 更新订单状态字段
func (r *orderRepository) UpdateStatus(orderID uint, status int) error {
	return r.db.Model(&model.Order{}).Where("id = ?", orderID).
		Update("status", status).Error
}

// UpdatePaid 支付成功，更新订单状态和支付时间
func (r *orderRepository) UpdatePaid(orderID uint) error {
	now := time.Now()
	return r.db.Model(&model.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":  model.OrderStatusPaid,
			"paid_at": &now,
		}).Error
}

// CountByUser 统计用户订单总数，支持按状态筛选
func (r *orderRepository) CountByUser(userID uint, status int) (int64, error) {
	var total int64
	query := r.db.Model(&model.Order{}).Where("user_id = ?", userID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	err := query.Count(&total).Error
	return total, err
}
