package repository

import (
	"fmt"

	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// FoodRepository 商品仓储接口
type FoodRepository interface {
	FindFoodsByCategory(categoryID uint) ([]model.Food, error) // 按分类查询商品
	FindFoodByID(id uint) (*model.Food, error)                 // 按 ID 查询商品
	DecreaseStock(foodID uint, quantity int) error             // 扣减库存（带库存不足校验）
	IncreaseStock(foodID uint, quantity int) error             // 恢复库存
}

// foodRepository 商品仓储的 GORM 实现
type foodRepository struct {
	db *gorm.DB
}

// NewFoodRepository 创建商品仓储实例
func NewFoodRepository(db *gorm.DB) FoodRepository {
	return &foodRepository{db: db}
}

// FindFoodsByCategory 查询分类下所有上架商品（status = 1）
func (r *foodRepository) FindFoodsByCategory(categoryID uint) ([]model.Food, error) {
	var foods []model.Food
	err := r.db.Where("category_id = ? AND status = ?", categoryID, 1).Find(&foods).Error
	if err != nil {
		return nil, err
	}
	return foods, nil
}

// FindFoodByID 按主键查询商品
func (r *foodRepository) FindFoodByID(id uint) (*model.Food, error) {
	var food model.Food
	err := r.db.First(&food, id).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

// DecreaseStock 扣减库存，使用 stock >= quantity 条件保证原子性，防止超卖
func (r *foodRepository) DecreaseStock(foodID uint, quantity int) error {
	result := r.db.Model(&model.Food{}).Where("id = ? AND stock >= ?", foodID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("insufficient stock for food %d", foodID)
	}
	return nil
}

// IncreaseStock 增加库存，用于订单取消时恢复 MySQL 库存
func (r *foodRepository) IncreaseStock(foodID uint, quantity int) error {
	return r.db.Model(&model.Food{}).Where("id = ?", foodID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
