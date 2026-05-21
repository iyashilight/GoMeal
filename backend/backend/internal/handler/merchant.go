package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type MerchantHandler struct {
	merchantSvc *service.MerchantService
}

func NewMerchantHandler(merchantSvc *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{merchantSvc: merchantSvc}
}

// GetMerchantList 获取商家列表
// @Summary 获取商家列表
// @Description 获取所有营业中的商家
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response{data=[]service.MerchantResponse}
// @Router       /merchants [get]
func (h *MerchantHandler) GetMerchantList(c *gin.Context) {
	resp, err := h.merchantSvc.GetMerchantList()
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// GetMerchantDetail 获取商家详情
// @Summary      获取商家详情
// @Description  获取商家信息及分类、食品列表
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        id path int true "商家 ID"
// @Success      200 {object} utils.Response{data=service.MerchantDetailResponse}
// @Router       /merchants/{id} [get]
func (h *MerchantHandler) GetMerchantDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid merchant id")
		return
	}

	resp, err := h.merchantSvc.GetMerchantDetail(uint(id))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// GetFoodDetail 获取食品详情
// @Summary      获取食品详情
// @Description  获取单个食品的详细信息
// @Tags         merchant
// @Accept       json
// @Produce      json
// @Param        id path int true "食品 ID"
// @Success      200 {object} utils.Response{data=service.FoodResponse}
// @Router       /foods/{id} [get]
func (h *MerchantHandler) GetFoodDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid food id")
		return
	}

	resp, err := h.merchantSvc.GetFoodDetail(uint(id))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}
