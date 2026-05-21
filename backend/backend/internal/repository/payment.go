package repository

import (
	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// PaymentRepository 支付仓储接口
type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByOrderID(orderID uint) (*model.Payment, error)
	UpdateStatus(id uint, status int) error
}

// paymentRepository 支付仓储的 GORM 实现
type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository 创建支付仓储实例
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByOrderID(orderID uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", id).Update("status", status).Error
}
