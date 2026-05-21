package repository

import (
	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	FindByMerchantID(merchantID uint) ([]model.Category, error) // 查询商家的所有分类
	Create(category *model.Category) error                     // 创建分类
}

// categoryRepository 分类仓储的 GORM 实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储实例
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// FindByMerchantID 查询指定商家下的所有分类，按 sort_order 升序排列
func (r *categoryRepository) FindByMerchantID(merchantID uint) ([]model.Category,
	error) {
	var categories []model.Category
	err := r.db.Where("Merchant_id = ?", merchantID).Order("sort_order asc").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Create 插入新分类
func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}
