package service

import (
	"log/slog"

	"github.com/iyashilight/GoMeal/internal/cache"
	"github.com/iyashilight/GoMeal/internal/model"
	"github.com/iyashilight/GoMeal/internal/repository"
)

// MerchantRegisterService 商家注册业务
type MerchantRegisterService struct {
	merchantRepo repository.MerchantRepository
	categoryRepo repository.CategoryRepository
	cache        *cache.Cache
}

// NewMerchantRegisterService 创建商家注册服务
func NewMerchantRegisterService(merchantRepo repository.MerchantRepository, categoryRepo repository.CategoryRepository, cache *cache.Cache) *MerchantRegisterService {
	return &MerchantRegisterService{merchantRepo: merchantRepo, categoryRepo: categoryRepo, cache: cache}
}

// RegisterMerchantRequest 商家注册请求
type RegisterMerchantRequest struct {
	Name        string  `json:"name" binding:"required,max=100"`
	Phone       string  `json:"phone" binding:"required"`
	Address     string  `json:"address" binding:"required,max=255"`
	Notice      string  `json:"notice" binding:"max=500"`
	MinPrice    float64 `json:"min_price"`
	DeliveryFee float64 `json:"delivery_fee"`
}

// RegisterMerchantResponse 商家注册响应
type RegisterMerchantResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	Notice      string  `json:"notice"`
	MinPrice    float64 `json:"min_price"`
	DeliveryFee float64 `json:"delivery_fee"`
	Status      int     `json:"status"`
}

// Register 商家注册
// 一个用户只能注册一个商家，重复注册返回 ErrMerchantConflict
func (s *MerchantRegisterService) Register(userID uint, req RegisterMerchantRequest) (*RegisterMerchantResponse, error) {
	// 1. 查重
	existing, err := s.merchantRepo.FindByUserID(userID)
	if err == nil {
		_ = existing
		slog.Warn("merchant register failed: user already has merchant", "user_id", userID)
		return nil, ErrMerchantConflict
	}

	// 2. 创建商家
	merchant := &model.Merchant{
		UserID:      userID,
		Name:        req.Name,
		Phone:       req.Phone,
		Address:     req.Address,
		Notice:      req.Notice,
		MinPrice:    req.MinPrice,
		DeliveryFee: req.DeliveryFee,
	}

	if err := s.merchantRepo.Create(merchant); err != nil {
		slog.Error("merchant register failed: create error", "user_id", userID, "error", err)
		return nil, err
	}

	// 创建默认分类
	defaultCat := &model.Category{
		MerchantID: merchant.ID,
		Name:       "全部商品",
		SortOrder:  1,
	}
	if err := s.categoryRepo.Create(defaultCat); err != nil {
		slog.Warn("failed to create default category", "merchant_id", merchant.ID, "error", err)
	}

	// 清除商家列表缓存
	_ = s.cache.Del("merchant:list")

	slog.Info("merchant registered", "merchant_id", merchant.ID, "user_id", userID, "name", req.Name)
	return &RegisterMerchantResponse{
		ID:          merchant.ID,
		Name:        merchant.Name,
		Phone:       merchant.Phone,
		Address:     merchant.Address,
		Notice:      merchant.Notice,
		MinPrice:    merchant.MinPrice,
		DeliveryFee: merchant.DeliveryFee,
		Status:      merchant.Status,
	}, nil
}
