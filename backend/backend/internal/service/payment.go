package service

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/iyashilight/GoMeal/internal/cache"
	"github.com/iyashilight/GoMeal/internal/model"
	"github.com/iyashilight/GoMeal/internal/repository"
	"gorm.io/gorm"
)

// PaymentService 支付业务
type PaymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	transactor  *repository.Transactor
	cache       *cache.Cache
}

// NewPaymentService 创建支付服务
func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	transactor *repository.Transactor,
	cache *cache.Cache,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		transactor:  transactor,
		cache:       cache,
	}
}

// PaymentResponse 支付响应
type PaymentResponse struct {
	PaymentID uint    `json:"payment_id"`
	OrderID   uint    `json:"order_id"`
	Amount    float64 `json:"amount"`
	Method    string  `json:"method"`
	TradeNo   string  `json:"trade_no"`
	Status    int     `json:"status"`
	PaidAt    string  `json:"paid_at,omitempty"`
}

// generateTradeNo 生成支付流水号：yyyyMMddHHmmss + 8位随机数
func generateTradeNo() string {
	now := time.Now()
	return fmt.Sprintf("P%s%08d", now.Format("20060102150405"), rand.Intn(100000000))
}

// toPaymentResponse 将 Payment 模型转为响应
func toPaymentResponse(p *model.Payment) *PaymentResponse {
	resp := &PaymentResponse{
		PaymentID: p.ID,
		OrderID:   p.OrderID,
		Amount:    p.Amount,
		Method:    p.Method,
		TradeNo:   p.TradeNo,
		Status:    p.Status,
	}
	if p.PaidAt != nil {
		resp.PaidAt = p.PaidAt.Format("2006-01-02 15:04:05")
	}
	return resp
}

// PayOrder 发起支付（mock 模拟）
// 幂等：同一 order_id 已有支付记录则直接返回已有结果
func (s *PaymentService) PayOrder(userID uint, orderID uint) (*PaymentResponse, error) {
	// 1. 校验订单
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotFound
	}
	if order.Status != model.OrderStatusPending {
		return nil, ErrInvalidOrderStatus
	}

	// 2. 幂等检查（已存在支付记录则直接返回）
	existing, err := s.paymentRepo.FindByOrderID(orderID)
	if err == nil {
		slog.Info("payment already exists, idempotent return", "order_id", orderID)
		return toPaymentResponse(existing), nil
	}

	// 3. 事务：创建支付记录 + 更新订单状态
	now := time.Now()
	payment := &model.Payment{
		OrderID: orderID,
		UserID:  userID,
		Amount:  order.TotalAmount,
		Method:  "mock",
		TradeNo: generateTradeNo(),
		Status:  1,
		PaidAt:  &now,
	}

	err = s.transactor.ExecTx(func(tx *repository.TxFactory) error {
		if err := tx.PaymentRepo.Create(payment); err != nil {
			return err
		}
		if err := tx.OrderRepo.UpdatePaid(orderID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.Error("pay order transaction failed", "order_id", orderID, "error", err)
		return nil, err
	}

	// 4. 清理订单缓存
	_ = s.cache.Del(fmt.Sprintf("food:stock:%d", orderID))

	slog.Info("order paid", "order_id", orderID, "amount", order.TotalAmount, "trade_no", payment.TradeNo)
	return toPaymentResponse(payment), nil
}

// GetPayment 查询支付状态
func (s *PaymentService) GetPayment(userID uint, orderID uint) (*PaymentResponse, error) {
	// 校验订单归属
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil || order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	payment, err := s.paymentRepo.FindByOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 未支付，返回 unpaid 状态
			return &PaymentResponse{
				OrderID: orderID,
				Amount:  order.TotalAmount,
				Status:  0,
			}, nil
		}
		return nil, err
	}

	return toPaymentResponse(payment), nil
}
