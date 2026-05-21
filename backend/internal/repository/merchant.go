package repository

import (
	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// MerchantRepository 商家仓储接口
type MerchantRepository interface {
	FindAll() ([]model.Merchant, error)                // 查询所有营业中的商家
	FindByID(id uint) (*model.Merchant, error)         // 按 ID 查询商家基本信息
	FindByIDWithDetail(id uint) (*model.Merchant, error) // 按 ID 查询商家及其分类、商品（预加载）
	FindByUserID(userID uint) (*model.Merchant, error) // 按店主用户ID 查询商家
	Create(merchant *model.Merchant) error             // 创建商家
	Update(merchant *model.Merchant) error             // 更新商家信息
}

// merchantRepository 商家仓储的 GORM 实现
type merchantRepository struct {
	db *gorm.DB
}

// NewMerchantRepository 创建商家仓储实例
func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

// FindAll 查询所有 status = 1（营业中）的商家
func (r *merchantRepository) FindAll() ([]model.Merchant, error) {
	var merchants []model.Merchant
	err := r.db.Where("status = ?", 1).Find(&merchants).Error
	if err != nil {
		return nil, err
	}
	return merchants, nil
}

// FindByID 按主键查询商家基本信息
func (r *merchantRepository) FindByID(id uint) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.First(&merchant, id).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// FindByIDWithDetail 查询商家详情，预加载其分类和商品（三层嵌套）
func (r *merchantRepository) FindByIDWithDetail(id uint) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.Preload("Categories.Foods").First(&merchant, id).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// FindByUserID 根据店主用户ID 查询商家
func (r *merchantRepository) FindByUserID(userID uint) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.db.Where("user_id = ?", userID).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// Create 插入商家记录
func (r *merchantRepository) Create(merchant *model.Merchant) error {
	return r.db.Create(merchant).Error
}

// Update 全字段更新商家信息
func (r *merchantRepository) Update(merchant *model.Merchant) error {
	return r.db.Save(merchant).Error
}
