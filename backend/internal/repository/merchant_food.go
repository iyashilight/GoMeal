package repository

import (
	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// MerchantFoodRepository 商家端商品管理仓储接口
type MerchantFoodRepository interface {
	Create(food *model.Food) error
	Update(food *model.Food) error
	Delete(id uint, merchantID uint) error
	FindByID(id uint) (*model.Food, error)
	FindByMerchantID(merchantID uint) ([]model.Food, error)
	UpdateStatus(id uint, merchantID uint, status int) error
}

// merchantFoodRepository 商家端商品仓储的 GORM 实现
type merchantFoodRepository struct {
	db *gorm.DB
}

// NewMerchantFoodRepository 创建商家商品仓储实例
func NewMerchantFoodRepository(db *gorm.DB) MerchantFoodRepository {
	return &merchantFoodRepository{db: db}
}

func (r *merchantFoodRepository) Create(food *model.Food) error {
	return r.db.Create(food).Error
}

func (r *merchantFoodRepository) Update(food *model.Food) error {
	return r.db.Save(food).Error
}

func (r *merchantFoodRepository) Delete(id uint, merchantID uint) error {
	return r.db.Where("id = ? AND merchant_id = ?", id, merchantID).Delete(&model.Food{}).Error
}

func (r *merchantFoodRepository) FindByID(id uint) (*model.Food, error) {
	var food model.Food
	err := r.db.First(&food, id).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *merchantFoodRepository) FindByMerchantID(merchantID uint) ([]model.Food, error) {
	var foods []model.Food
	err := r.db.Where("merchant_id = ?", merchantID).Order("created_at desc").Find(&foods).Error
	if err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *merchantFoodRepository) UpdateStatus(id uint, merchantID uint, status int) error {
	return r.db.Model(&model.Food{}).
		Where("id = ? AND merchant_id = ?", id, merchantID).
		Update("status", status).Error
}
