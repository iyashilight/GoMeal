package service

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/iyashilight/GoMeal/internal/cache"
	"github.com/iyashilight/GoMeal/internal/model"
	"github.com/iyashilight/GoMeal/internal/repository"
	"github.com/iyashilight/GoMeal/internal/utils"
)

// OrderService 订单业务，包含 Redis 预扣库存 + MySQL 事务两阶段提交
type OrderService struct {
	foodRepo    repository.FoodRepository
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	addressRepo repository.AddressRepository
	transactor  *repository.Transactor
	stockCache  *cache.StockCache
}

// NewOrderService 创建订单服务
func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	addressRepo repository.AddressRepository,
	foodRepo repository.FoodRepository,
	transactor *repository.Transactor,
	stockCache *cache.StockCache,
) *OrderService {
	return &OrderService{
		foodRepo:    foodRepo,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		addressRepo: addressRepo,
		transactor:  transactor,
		stockCache:  stockCache,
	}
}

// decrementItem 记录已扣减库存的商品，用于失败时回滚
type decrementItem struct {
	FoodID   uint
	Quantity int
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	AddressID   uint    `json:"address_id" binding:"required"`
	DeliveryFee float64 `json:"delivery_fee" binding:"min=0"`
	Remark      string  `json:"remark" binding:"max=200"`
}

// OrderItemResponse 订单商品响应
type OrderItemResponse struct {
	FoodID    uint    `json:"food_id"`
	FoodName  string  `json:"food_name"`
	FoodImage string  `json:"food_image"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID          uint                `json:"id"`
	OrderNo     string              `json:"order_no"`
	Status      int                 `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	DeliveryFee float64             `json:"delivery_fee"`
	Remark      string              `json:"remark"`
	Items       []OrderItemResponse `json:"items"`
	CreatedAt   time.Time           `json:"created_at"`
}

// OrderListResponse 订单列表分页响应
type OrderListResponse struct {
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
	Items []OrderResponse `json:"items"`
}

// generateOrderNo 生成订单号：yyyymmddhhmmss + 6位随机数，共20位
func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("%s%06d",
		now.Format("20060102150405"),
		rand.Intn(1000000),
	)
}

// CreateOrder 创建订单（核心流程）
//
// 两阶段库存扣减：
//   Phase 1 — Redis 预扣库存（高速缓存层）
//   Phase 2 — MySQL 事务写入订单 + 扣减实际库存 + 清空购物车
//
// 如果 Redis 操作失败，降级为只走 MySQL（redisFallback），保证可用性。
// 如果 MySQL 事务失败，回滚 Redis 预扣的库存。
func (s *OrderService) CreateOrder(userID uint, req *CreateOrderRequest) (*OrderResponse, error) {
	carts, err := s.cartRepo.GetCartByUser(userID)
	if err != nil || len(carts) == 0 {
		slog.Warn("create order failed: empty cart", "user_id", userID)
		return nil, ErrCartEmpty
	}

	addr, err := s.addressRepo.FindByID(req.AddressID)
	if err != nil || addr.UserID != userID {
		slog.Warn("create order failed: address not found", "user_id", userID, "address_id",
			req.AddressID)
		return nil, ErrAddressNotFound
	}

	merchantID := carts[0].MerchantID

	items := make([]OrderItemResponse, len(carts))
	totalAmount := 0.0
	for i, c := range carts {
		totalAmount += c.Food.Price * float64(c.Quantity)
		items[i] = OrderItemResponse{
			FoodID:    c.FoodID,
			FoodName:  c.Food.Name,
			FoodImage: c.Food.Image,
			Price:     c.Food.Price,
			Quantity:  c.Quantity,
		}
	}

	orderItems := make([]model.OrderItem, len(carts))
	for i, c := range carts {
		orderItems[i] = model.OrderItem{
			FoodID:    c.FoodID,
			FoodName:  c.Food.Name,
			FoodImage: c.Food.Image,
			Price:     c.Food.Price,
			Quantity:  c.Quantity,
		}
	}
	order := &model.Order{
		OrderNo:     generateOrderNo(),
		UserID:      userID,
		MerchantID:  merchantID,
		AddressID:   req.AddressID,
		TotalAmount: totalAmount + req.DeliveryFee,
		DeliveryFee: req.DeliveryFee,
		Remark:      req.Remark,
		Status:      model.OrderStatusPending,
		Items:       orderItems,
	}

	// ---- Phase 1: Redis 预扣库存 ----
	var decrement []decrementItem
	redisOK := true
	for _, c := range carts {
		res, err := s.stockCache.Decrease(c.FoodID, c.Quantity)
		if err != nil {
			goto redisFallback
		}
		switch res {
		case 1:
			decrement = append(decrement, decrementItem{c.FoodID, c.Quantity})
		case 0:
			for _, d := range decrement {
				_ = s.stockCache.Increase(d.FoodID, d.Quantity)
			}
			return nil, ErrInsufficientStock
		case -1:
			// key 不存在，从 MySQL 加载库存后重试
			food, err := s.foodRepo.FindFoodByID(c.FoodID)
			if err != nil {
				goto redisFallback
			}
			_ = s.stockCache.Init(c.FoodID, food.Stock)
			res, err = s.stockCache.Decrease(c.FoodID, c.Quantity)
			if err != nil || res == -1 {
				goto redisFallback
			}
			if res == 0 {
				for _, d := range decrement {
					_ = s.stockCache.Increase(d.FoodID, d.Quantity)
				}
				return nil, ErrInsufficientStock
			}
			decrement = append(decrement, decrementItem{c.FoodID, c.Quantity})
		}
	}
	goto execTx

redisFallback:
	// Redis 不可用时降级，仅依赖 MySQL
	for _, d := range decrement {
		_ = s.stockCache.Increase(d.FoodID, d.Quantity)
	}
	decrement = nil
	redisOK = false

execTx:
	// ---- Phase 2: MySQL 事务写入 ----
	err = s.transactor.ExecTx(func(tx *repository.TxFactory) error {
		if err := tx.OrderRepo.Create(order); err != nil {
			return err
		}
		for _, c := range carts {
			if err := tx.FoodRepo.DecreaseStock(c.FoodID, c.Quantity); err != nil {
				return err
			}
		}
		return tx.CartRepo.ClearByUser(userID)
	})
	if err != nil {
		if redisOK {
			for _, d := range decrement {
				_ = s.stockCache.Increase(d.FoodID, d.Quantity)
			}
		}
		slog.Error("create order transaction failed", "user_id", userID, "error", err)
		return nil, err
	}
	slog.Info("order created", "order_id", order.ID, "order_no", order.OrderNo,
		"user_id", userID, "amount", order.TotalAmount)
	return &OrderResponse{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		DeliveryFee: order.DeliveryFee,
		Remark:      order.Remark,
		Items:       items,
		CreatedAt:   order.CreatedAt,
	}, nil
}

// GetOrders 分页获取用户订单，支持按状态筛选
func (s *OrderService) GetOrders(userID uint, status int, pageReq utils.PageReq) (*OrderListResponse, error) {
	pageReq.Normalize()
	total, err := s.orderRepo.CountByUser(userID, status)
	if err != nil {
		return nil, err
	}
	orders, err := s.orderRepo.FindByUser(userID, status, pageReq.Page, pageReq.Size)
	if err != nil {
		return nil, err
	}
	items := make([]OrderResponse, len(orders))
	for i, o := range orders {
		itemResponses := make([]OrderItemResponse, len(o.Items))
		for j, item := range o.Items {
			itemResponses[j] = OrderItemResponse{
				FoodID:    item.FoodID,
				FoodName:  item.FoodName,
				FoodImage: item.FoodImage,
				Price:     item.Price,
				Quantity:  item.Quantity,
			}
		}
		items[i] = OrderResponse{
			ID:          o.ID,
			OrderNo:     o.OrderNo,
			Status:      o.Status,
			TotalAmount: o.TotalAmount,
			DeliveryFee: o.DeliveryFee,
			Remark:      o.Remark,
			Items:       itemResponses,
			CreatedAt:   o.CreatedAt,
		}
	}

	return &OrderListResponse{
		Total: total,
		Page:  pageReq.Page,
		Size:  pageReq.Size,
		Items: items,
	}, nil
}

// GetOrderDetail 获取单个订单详情，校验用户归属
func (s *OrderService) GetOrderDetail(userID, orderID uint) (*OrderResponse, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrOrderNotFound
	}

	items := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemResponse{
			FoodID:    item.FoodID,
			FoodName:  item.FoodName,
			FoodImage: item.FoodImage,
			Price:     item.Price,
			Quantity:  item.Quantity,
		}
	}

	return &OrderResponse{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		DeliveryFee: order.DeliveryFee,
		Remark:      order.Remark,
		Items:       items,
		CreatedAt:   order.CreatedAt,
	}, nil
}

// CancelOrder 取消订单
// 仅待支付状态可取消，同时恢复 MySQL 和 Redis 库存
func (s *OrderService) CancelOrder(userID, orderID uint) error {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		slog.Warn("cancel order failed: not found", "order_id", orderID, "user_id", userID)
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrOrderNotFound
	}
	if order.Status != model.OrderStatusPending {
		slog.Warn("cancel order failed: invalid status", "order_id", orderID, "user_id",
			userID, "status", order.Status)
		return ErrInvalidOrderStatus
	}
	err = s.transactor.ExecTx(func(tx *repository.TxFactory) error {
		err := tx.OrderRepo.UpdateStatus(orderID, model.OrderStatusCancelled)
		if err != nil {
			return err
		}
		for _, item := range order.Items {
			if err := tx.FoodRepo.IncreaseStock(item.FoodID, item.Quantity); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		slog.Error("cancel order transaction failed", "order_id", orderID, "user_id",
			userID, "error", err)
		return err
	}
	// 恢复 Redis 库存
	for _, item := range order.Items {
		_ = s.stockCache.Increase(item.FoodID, item.Quantity)
	}
	slog.Info("order cancelled", "order_id", orderID, "user_id", userID)
	return nil
}
