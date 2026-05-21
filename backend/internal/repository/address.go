package repository

import (
	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// AddressRepository 地址仓储接口
type AddressRepository interface {
	Create(address *model.Address) error          // 创建地址
	FindByUser(userID uint) ([]model.Address, error) // 查询用户所有地址
	FindByID(id uint) (*model.Address, error)     // 按 ID 查询地址
	Update(address *model.Address) error          // 更新地址
	Delete(userID, addressID uint) error          // 删除地址（校验用户归属）
	ClearDefault(userID uint) error               // 清除用户所有地址的默认标记
}

// addressRepository 地址仓储的 GORM 实现
type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository 创建地址仓储实例
func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

// Create 插入新地址记录
func (r *addressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

// FindByUser 查询指定用户的所有地址
func (r *addressRepository) FindByUser(userID uint) ([]model.Address, error) {
	var address []model.Address
	err := r.db.Where("user_id = ?", userID).Find(&address).Error
	if err != nil {
		return nil, err
	}
	return address, nil
}

// FindByID 按主键 ID 查询地址
func (r *addressRepository) FindByID(id uint) (*model.Address, error) {
	var address model.Address
	err := r.db.First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Update 全字段更新地址记录
func (r *addressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

// Delete 按用户ID + 地址ID 删除（防止越权删除他人地址）
func (r *addressRepository) Delete(userID, addressID uint) error {
	return r.db.Where("user_id = ? AND id = ?", userID, addressID).
		Delete(&model.Address{}).Error
}

// ClearDefault 将用户所有地址的 is_default 置为 false
func (r *addressRepository) ClearDefault(userID uint) error {
	return r.db.Model(&model.Address{}).Where("user_id = ?", userID).
		Update("is_default", false).Error
}
