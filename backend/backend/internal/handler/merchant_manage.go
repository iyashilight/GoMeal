package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type MerchantManageHandler struct {
	manageSvc *service.MerchantManageService
}

func NewMerchantManageHandler(manageSvc *service.MerchantManageService) *MerchantManageHandler {
	return &MerchantManageHandler{manageSvc: manageSvc}
}

// GetMyFoods 商家查看自己的商品列表
func (h *MerchantManageHandler) GetMyFoods(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.manageSvc.GetMyFoods(userID)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// CreateFood 商家创建商品
func (h *MerchantManageHandler) CreateFood(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.manageSvc.CreateFood(userID, req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// UpdateFood 商家修改商品
func (h *MerchantManageHandler) UpdateFood(c *gin.Context) {
	userID := middleware.GetUserID(c)

	foodID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid food id")
		return
	}

	var req service.UpdateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.manageSvc.UpdateFood(userID, uint(foodID), req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// DeleteFood 商家删除商品
func (h *MerchantManageHandler) DeleteFood(c *gin.Context) {
	userID := middleware.GetUserID(c)

	foodID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid food id")
		return
	}

	if err := h.manageSvc.DeleteFood(userID, uint(foodID)); err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}

// SetFoodStatus 设置商品上下架状态
func (h *MerchantManageHandler) SetFoodStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)

	foodID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid food id")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.manageSvc.SetFoodStatus(userID, uint(foodID), req.Status)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}
