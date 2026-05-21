package service

import (
	"fmt"
	"log/slog"

	"github.com/iyashilight/GoMeal/internal/cache"
	"github.com/iyashilight/GoMeal/internal/model"
	"github.com/iyashilight/GoMeal/internal/repository"
)

// MerchantManageService 商家管理业务
type MerchantManageService struct {
	merchantRepo repository.MerchantRepository
	categoryRepo repository.CategoryRepository
	foodRepo     repository.MerchantFoodRepository
	cache        *cache.Cache
}

// NewMerchantManageService 创建商家管理服务
func NewMerchantManageService(
	merchantRepo repository.MerchantRepository,
	categoryRepo repository.CategoryRepository,
	foodRepo repository.MerchantFoodRepository,
	cache *cache.Cache,
) *MerchantManageService {
	return &MerchantManageService{
		merchantRepo: merchantRepo,
		categoryRepo: categoryRepo,
		foodRepo:     foodRepo,
		cache:        cache,
	}
}

// --- Request / Response ---

type CreateFoodRequest struct {
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name" binding:"required,max=100"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	OldPrice    float64 `json:"old_price"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Description string  `json:"description" binding:"max=500"`
	Image       string  `json:"image" binding:"max=255"`
}

type UpdateFoodRequest struct {
	CategoryID  uint     `json:"category_id"`
	Name        string   `json:"name" binding:"max=100"`
	Price       float64  `json:"price" binding:"gt=0"`
	OldPrice    *float64 `json:"old_price"`
	Stock       *int     `json:"stock" binding:"omitempty,gte=0"`
	Description string   `json:"description" binding:"max=500"`
	Image       string   `json:"image" binding:"max=255"`
}

type FoodManageResponse struct {
	ID          uint    `json:"id"`
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	OldPrice    float64 `json:"old_price"`
	Stock       int     `json:"stock"`
	Sales       int     `json:"sales"`
	Status      int     `json:"status"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// --- 工具方法 ---

// getMerchantByUserID 根据用户 ID 获取商家，若不是商家则返回 ErrNotMerchant
func (s *MerchantManageService) getMerchantByUserID(userID uint) (*model.Merchant, error) {
	merchant, err := s.merchantRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrNotMerchant
	}
	return merchant, nil
}

// toFoodManageResponse 将 Food 模型转为响应结构
func toFoodManageResponse(food *model.Food) FoodManageResponse {
	return FoodManageResponse{
		ID:          food.ID,
		CategoryID:  food.CategoryID,
		Name:        food.Name,
		Price:       food.Price,
		OldPrice:    food.OldPrice,
		Stock:       food.Stock,
		Sales:       food.Sales,
		Status:      food.Status,
		Description: food.Description,
		Image:       food.Image,
		CreatedAt:   food.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   food.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// invalidateFoodCache 清除与商品相关的缓存
func (s *MerchantManageService) invalidateFoodCache(foodID uint, merchantID uint) {
	_ = s.cache.Del(fmt.Sprintf("food:%d", foodID))
	_ = s.cache.Del("merchant:list")
	_ = s.cache.Del(fmt.Sprintf("merchant:detail:%d", merchantID))
}

// --- 核心业务方法 ---

// GetMyFoods 获取当前商家所有商品
func (s *MerchantManageService) GetMyFoods(userID uint) ([]FoodManageResponse, error) {
	merchant, err := s.getMerchantByUserID(userID)
	if err != nil {
		return nil, err
	}

	foods, err := s.foodRepo.FindByMerchantID(merchant.ID)
	if err != nil {
		slog.Error("failed to get merchant foods", "merchant_id", merchant.ID, "error", err)
		return nil, err
	}

	result := make([]FoodManageResponse, len(foods))
	for i, f := range foods {
		result[i] = toFoodManageResponse(&f)
	}
	return result, nil
}

// CreateFood 创建商品
func (s *MerchantManageService) CreateFood(userID uint, req CreateFoodRequest) (*FoodManageResponse, error) {
	merchant, err := s.getMerchantByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 若未指定分类，自动使用商家的第一个分类
	categoryID := req.CategoryID
	if categoryID == 0 {
		categories, err := s.categoryRepo.FindByMerchantID(merchant.ID)
		if err == nil && len(categories) > 0 {
			categoryID = categories[0].ID
		}
	}

	food := &model.Food{
		MerchantID:  merchant.ID,
		CategoryID:  categoryID,
		Name:        req.Name,
		Price:       req.Price,
		OldPrice:    req.OldPrice,
		Stock:       req.Stock,
		Description: req.Description,
		Image:       req.Image,
		Status:      1,
	}

	if err := s.foodRepo.Create(food); err != nil {
		slog.Error("failed to create food", "merchant_id", merchant.ID, "error", err)
		return nil, err
	}

	s.invalidateFoodCache(food.ID, merchant.ID)

	slog.Info("food created", "food_id", food.ID, "merchant_id", merchant.ID)
	result := toFoodManageResponse(food)
	return &result, nil
}

// UpdateFood 修改商品
func (s *MerchantManageService) UpdateFood(userID uint, foodID uint, req UpdateFoodRequest) (*FoodManageResponse, error) {
	merchant, err := s.getMerchantByUserID(userID)
	if err != nil {
		return nil, err
	}

	food, err := s.foodRepo.FindByID(foodID)
	if err != nil {
		return nil, ErrFoodNotFound
	}

	if food.MerchantID != merchant.ID {
		return nil, ErrFoodNotBelong
	}

	if req.Name != "" {
		food.Name = req.Name
	}
	if req.Price > 0 {
		food.Price = req.Price
	}
	if req.OldPrice != nil {
		food.OldPrice = *req.OldPrice
	}
	if req.Stock != nil {
		food.Stock = *req.Stock
	}
	if req.CategoryID > 0 {
		food.CategoryID = req.CategoryID
	}
	if req.Description != "" {
		food.Description = req.Description
	}
	if req.Image != "" {
		food.Image = req.Image
	}

	if err := s.foodRepo.Update(food); err != nil {
		slog.Error("failed to update food", "food_id", foodID, "error", err)
		return nil, err
	}

	s.invalidateFoodCache(foodID, merchant.ID)

	slog.Info("food updated", "food_id", foodID, "merchant_id", merchant.ID)
	result := toFoodManageResponse(food)
	return &result, nil
}

// DeleteFood 删除商品
func (s *MerchantManageService) DeleteFood(userID uint, foodID uint) error {
	merchant, err := s.getMerchantByUserID(userID)
	if err != nil {
		return err
	}

	food, err := s.foodRepo.FindByID(foodID)
	if err != nil {
		return ErrFoodNotFound
	}

	if food.MerchantID != merchant.ID {
		return ErrFoodNotBelong
	}

	if err := s.foodRepo.Delete(foodID, merchant.ID); err != nil {
		slog.Error("failed to delete food", "food_id", foodID, "error", err)
		return err
	}

	s.invalidateFoodCache(foodID, merchant.ID)

	slog.Info("food deleted", "food_id", foodID, "merchant_id", merchant.ID)
	return nil
}

// SetFoodStatus 设置商品上下架状态（0=下架, 1=上架）
func (s *MerchantManageService) SetFoodStatus(userID uint, foodID uint, status int) (*FoodManageResponse, error) {
	merchant, err := s.getMerchantByUserID(userID)
	if err != nil {
		return nil, err
	}

	food, err := s.foodRepo.FindByID(foodID)
	if err != nil {
		return nil, ErrFoodNotFound
	}

	if food.MerchantID != merchant.ID {
		return nil, ErrFoodNotBelong
	}

	if err := s.foodRepo.UpdateStatus(foodID, merchant.ID, status); err != nil {
		slog.Error("failed to update food status", "food_id", foodID, "error", err)
		return nil, err
	}

	food.Status = status
	s.invalidateFoodCache(foodID, merchant.ID)

	slog.Info("food status updated", "food_id", foodID, "status", status)
	result := toFoodManageResponse(food)
	return &result, nil
}
