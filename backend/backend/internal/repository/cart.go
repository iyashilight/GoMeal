package repository

import (
	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// CartRepository 购物车仓储接口
type CartRepository interface {
	GetCartByUser(userID uint) ([]model.Cart, error)    // 查询用户购物车（含商品和商家信息）
	GetCartItem(userID, foodID uint) (*model.Cart, error) // 查询用户特定商品的购物车记录
	Create(cart *model.Cart) error                       // 添加购物车条目
	Update(cart *model.Cart) error                       // 更新购物车条目
	Delete(userID, cartID uint) error                    // 删除购物车条目（校验用户归属）
	ClearByUser(userID uint) error                       // 清空用户购物车
}

// cartRepository 购物车仓储的 GORM 实现
type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository 创建购物车仓储实例
func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

// GetCartByUser 查询用户的购物车，预加载 Food 和 Merchant 关联数据
func (r *cartRepository) GetCartByUser(userID uint) ([]model.Cart, error) {
	var carts []model.Cart
	err := r.db.Where("user_id = ?", userID).Preload("Food").Preload("Merchant").
		Find(&carts).Error
	if err != nil {
		return nil, err
	}
	return carts, nil
}

// GetCartItem 查询用户购物车中是否已包含某商品
func (r *cartRepository) GetCartItem(userID, foodID uint) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("user_id =? AND food_id = ?", userID, foodID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// Create 插入购物车记录
func (r *cartRepository) Create(cart *model.Cart) error {
	return r.db.Create(cart).Error
}

// Update 全字段更新购物车记录
func (r *cartRepository) Update(cart *model.Cart) error {
	return r.db.Save(cart).Error
}

// Delete 按用户和购物车ID删除记录
func (r *cartRepository) Delete(userID, cartID uint) error {
	return r.db.Where("user_id = ? AND id = ?", userID, cartID).
		Delete(&model.Cart{}).Error
}

// ClearByUser 删除用户的所有购物车记录
func (r *cartRepository) ClearByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).
		Delete(&model.Cart{}).Error
}
