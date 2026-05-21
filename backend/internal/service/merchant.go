package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/iyashilight/GoMeal/internal/cache"
	"github.com/iyashilight/GoMeal/internal/repository"
)

// MerchantService 商家业务，包含 Redis 缓存逻辑
type MerchantService struct {
	merchantRepo repository.MerchantRepository
	categoryRepo repository.CategoryRepository
	foodRepo     repository.FoodRepository
	cache        *cache.Cache       // 通用 JSON 缓存
	stockCache   *cache.StockCache   // 库存缓存
}

// NewMerchantService 创建商家服务
func NewMerchantService(
	merchantRepo repository.MerchantRepository,
	categoryRepo repository.CategoryRepository,
	foodRepo repository.FoodRepository,
	cache *cache.Cache,
	stockCache *cache.StockCache,
) *MerchantService {
	return &MerchantService{
		merchantRepo: merchantRepo,
		categoryRepo: categoryRepo,
		foodRepo:     foodRepo,
		cache:        cache,
		stockCache:   stockCache,
	}
}

// MerchantResponse 商家列表响应
type MerchantResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Logo        string  `json:"logo"`
	Notice      string  `json:"notice"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	MinPrice    float64 `json:"min_price"`
	DeliveryFee float64 `json:"delivery_fee"`
	Status      int     `json:"status"`
	Rating      float64 `json:"rating"`
	Sales       int     `json:"sales"`
}

// CategoryWithFoods 分类及商品响应
type CategoryWithFoods struct {
	ID    uint           `json:"id"`
	Name  string         `json:"name"`
	Foods []FoodResponse `json:"foods"`
}

// FoodResponse 商品响应
type FoodResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	OldPrice    float64 `json:"old_price"`
	Stock       int     `json:"stock"`
	Sales       int     `json:"sales"`
}

// MerchantDetailResponse 商家详情响应（含分类和商品）
type MerchantDetailResponse struct {
	MerchantResponse
	Categories []CategoryWithFoods `json:"categories"`
}

// GetMerchantList 获取商家列表，优先从缓存读取，缓存 TTL 5 分钟
func (s *MerchantService) GetMerchantList() ([]MerchantResponse, error) {
	var result []MerchantResponse
	hit, err := s.cache.GetJSON("merchant:list", &result)
	if err != nil {
		slog.Warn("cache get failed", "key", "merchant:list", "error", err)
	}
	if hit {
		return result, nil
	}
	merchants, err := s.merchantRepo.FindAll()
	if err != nil {
		return nil, ErrMerchantNotFound
	}
	result = make([]MerchantResponse, len(merchants))
	for i, m := range merchants {
		result[i] = MerchantResponse{
			ID:          m.ID,
			Name:        m.Name,
			Logo:        m.Logo,
			Notice:      m.Notice,
			Phone:       m.Phone,
			Address:     m.Address,
			MinPrice:    m.MinPrice,
			DeliveryFee: m.DeliveryFee,
			Status:      m.Status,
			Rating:      m.Rating,
			Sales:       m.Sales,
		}
	}
	if err := s.cache.SetJSON("merchant:list", result, 5*time.Minute); err != nil {
		slog.Warn("cache set failed", "key", "merchant:list", "error", err)
	}
	return result, nil
}

// GetFoodDetail 获取商品详情，优先从缓存读取
// 同时将库存同步到 Redis，供下单时预扣使用
func (s *MerchantService) GetFoodDetail(id uint) (*FoodResponse, error) {
	var result FoodResponse
	hit, err := s.cache.GetJSON(fmt.Sprintf("food:%d", id), &result)
	if err != nil {
		slog.Warn("cache get failed", "key", fmt.Sprintf("food:%d", id), "error", err)
	}
	if hit {
		return &result, nil
	}
	food, err := s.foodRepo.FindFoodByID(id)
	if err != nil {
		return nil, ErrFoodNotFound
	}
	result = FoodResponse{
		ID:          food.ID,
		Name:        food.Name,
		Description: food.Description,
		Image:       food.Image,
		Price:       food.Price,
		OldPrice:    food.OldPrice,
		Stock:       food.Stock,
		Sales:       food.Sales,
	}
	if err := s.cache.SetJSON(fmt.Sprintf("food:%d", id), &result, 5*time.Minute); err != nil {
		slog.Warn("cache set failed", "key", fmt.Sprintf("food:%d", id), "error", err)
	}
	_ = s.stockCache.Init(id, food.Stock)
	return &result, nil
}

// GetMerchantDetail 获取商家详情（含分类和商品），优先从缓存读取
func (s *MerchantService) GetMerchantDetail(id uint) (*MerchantDetailResponse, error) {
	key := fmt.Sprintf("merchant:detail:%d", id)
	var result MerchantDetailResponse
	hit, err := s.cache.GetJSON(key, &result)
	if err != nil {
		slog.Warn("cache get failed", "key", key, "error", err)
	}
	if hit {
		return &result, nil
	}
	merchant, err := s.merchantRepo.FindByIDWithDetail(id)
	if err != nil {
		return nil, ErrMerchantNotFound
	}
	categoriesWithFoods := make([]CategoryWithFoods, len(merchant.Categories))
	for i, c := range merchant.Categories {
		foods := make([]FoodResponse, len(c.Foods))
		for j, f := range c.Foods {
			foods[j] = FoodResponse{
				ID:          f.ID,
				Name:        f.Name,
				Description: f.Description,
				Image:       f.Image,
				Price:       f.Price,
				OldPrice:    f.OldPrice,
				Stock:       f.Stock,
				Sales:       f.Sales,
			}
		}
		categoriesWithFoods[i] = CategoryWithFoods{
			ID:    c.ID,
			Name:  c.Name,
			Foods: foods,
		}
	}
	result = MerchantDetailResponse{
		MerchantResponse: MerchantResponse{
			ID:          merchant.ID,
			Name:        merchant.Name,
			Logo:        merchant.Logo,
			Notice:      merchant.Notice,
			Phone:       merchant.Phone,
			Address:     merchant.Address,
			MinPrice:    merchant.MinPrice,
			DeliveryFee: merchant.DeliveryFee,
			Status:      merchant.Status,
			Rating:      merchant.Rating,
			Sales:       merchant.Sales,
		},
		Categories: categoriesWithFoods,
	}
	if err := s.cache.SetJSON(key, &result, 5*time.Minute); err != nil {
		slog.Warn("cache set failed", "key", key, "error", err)
	}
	return &result, nil
}
