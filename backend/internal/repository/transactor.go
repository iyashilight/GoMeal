package repository

import "gorm.io/gorm"

type Transactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) *Transactor {
	return &Transactor{db: db}
}

type TxFactory struct {
	OrderRepo   *orderRepository
	CartRepo    *cartRepository
	FoodRepo    *foodRepository
	PaymentRepo *paymentRepository
}

func (t *Transactor) ExecTx(fn func(tx *TxFactory) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		factory := &TxFactory{
			OrderRepo:   &orderRepository{db: tx},
			CartRepo:    &cartRepository{db: tx},
			FoodRepo:    &foodRepository{db: tx},
			PaymentRepo: &paymentRepository{db: tx},
		}
		return fn(factory)
	})
}
