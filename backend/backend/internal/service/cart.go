package service

import (
	"log/slog"

	"github.com/iyashilight/GoMeal/internal/model"
	"github.com/iyashilight/GoMeal/internal/repository"
)

// CartService 购物车业务
type CartService struct {
	cartRepo repository.CartRepository
	foodRepo repository.FoodRepository
}

// NewCartService 创建购物车服务
func NewCartService(cartRepo repository.CartRepository,
	foodRepo repository.FoodRepository,
) *CartService {
	return &CartService{
		cartRepo: cartRepo,
		foodRepo: foodRepo,
	}
}

// AddCartRequest 添加购物车请求
type AddCartRequest struct {
	FoodID     uint `json:"food_id" binding:"required"`
	MerchantID uint `json:"merchant_id" binding:"required"`
	Quantity   int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartQuantityRequest 更新数量请求
type UpdateCartQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// CartItemResponse 购物车商品响应
type CartItemResponse struct {
	ID         uint              `json:"id"`
	FoodID     uint              `json:"food_id"`
	FoodName   string            `json:"food_name"`
	FoodImage  string            `json:"food_image"`
	Price      float64           `json:"price"`
	Quantity   int               `json:"quantity"`
	MerchantID uint              `json:"merchant_id"`
	Merchant   *MerchantResponse `json:"merchant"`
}

// CartResponse 购物车响应，按商家分组展示
type CartResponse struct {
	MerchantID    uint               `json:"merchant_id"`
	Merchant      *MerchantResponse  `json:"merchant"`
	Items         []CartItemResponse `json:"items"`
	TotalPrice    float64            `json:"total_price"`
	TotalQuantity int                `json:"total_quantity"`
}

// AddItem 添加商品到购物车
// 1. 检查商品是否存在且属于指定商家
// 2. 如果购物车已有同款商品，增加数量
// 3. 如果购物车已有其他商家商品，阻止跨商家添加
// 4. 否则新建购物车条目
func (s *CartService) AddItem(userID uint, req *AddCartRequest) error {
	food, err := s.foodRepo.FindFoodByID(req.FoodID)
	if err != nil {
		return ErrFoodNotFound
	}
	if food.MerchantID != req.MerchantID {
		return ErrFoodNotFound
	}

	existing, err := s.cartRepo.GetCartItem(userID, req.FoodID)
	if err == nil {
		existing.Quantity += req.Quantity
		return s.cartRepo.Update(existing)
	}

	carts, _ := s.cartRepo.GetCartByUser(userID)
	if len(carts) > 0 && carts[0].MerchantID != req.MerchantID {
		slog.Warn("add item failed: merchant conflict", "user_id", userID, "food_id", req.FoodID)
		return ErrMerchantConflict
	}
	cart := &model.Cart{
		UserID:     userID,
		MerchantID: req.MerchantID,
		FoodID:     req.FoodID,
		Quantity:   req.Quantity,
	}
	err = s.cartRepo.Create(cart)
	if err != nil {
		return err
	}
	slog.Info("cart item added", "user_id", userID, "food_id", req.FoodID, "quantity", req.Quantity)
	return nil
}

// UpdateQuantity 更新购物车商品数量
// quantity = 0 时删除该条目（方便减到零时自动清除）
func (s *CartService) UpdateQuantity(userID, cartID uint, quantity int) error {
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	if quantity == 0 {
		return s.cartRepo.Delete(userID, cartID)
	}

	carts, err := s.cartRepo.GetCartByUser(userID)
	if err != nil {
		return ErrCartNotFound
	}

	var target *model.Cart
	for i := range carts {
		if carts[i].ID == cartID {
			target = &carts[i]
			break
		}
	}
	if target == nil {
		return ErrCartNotFound
	}

	target.Quantity = quantity
	return s.cartRepo.Update(target)
}

// RemoveItem 从购物车删除指定商品
func (s *CartService) RemoveItem(userID, cartID uint) error {
	return s.cartRepo.Delete(userID, cartID)
}

// GetCart 获取用户购物车，计算总价和总数量
// 返回按商家聚合的购物车视图
func (s *CartService) GetCart(userID uint) (*CartResponse, error) {
	carts, err := s.cartRepo.GetCartByUser(userID)
	if err != nil {
		return nil, err
	}
	if len(carts) == 0 {
		return nil, ErrCartEmpty
	}

	items := make([]CartItemResponse, len(carts))
	for i, c := range carts {
		items[i] = CartItemResponse{
			ID:         c.ID,
			FoodID:     c.FoodID,
			FoodName:   c.Food.Name,
			FoodImage:  c.Food.Image,
			Price:      c.Food.Price,
			Quantity:   c.Quantity,
			MerchantID: c.MerchantID,
			Merchant: &MerchantResponse{
				ID:          c.Merchant.ID,
				Name:        c.Merchant.Name,
				Logo:        c.Merchant.Logo,
				Notice:      c.Merchant.Notice,
				Phone:       c.Merchant.Phone,
				Address:     c.Merchant.Address,
				MinPrice:    c.Merchant.MinPrice,
				DeliveryFee: c.Merchant.DeliveryFee,
				Status:      c.Merchant.Status,
				Rating:      c.Merchant.Rating,
				Sales:       c.Merchant.Sales,
			},
		}
	}

	totalPrice := 0.0
	totalQuantity := 0
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
		totalQuantity += item.Quantity
	}

	return &CartResponse{
		MerchantID:    carts[0].MerchantID,
		Merchant:      items[0].Merchant,
		Items:         items,
		TotalPrice:    totalPrice,
		TotalQuantity: totalQuantity,
	}, nil
}

// ClearCart 清空用户购物车
func (s *CartService) ClearCart(userID uint) error {
	return s.cartRepo.ClearByUser(userID)
}
