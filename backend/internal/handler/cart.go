package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type CartHandler struct {
	cartSvc *service.CartService
}

func NewCartHandler(cartSvc *service.CartService) *CartHandler {
	return &CartHandler{cartSvc: cartSvc}
}

// AddItem 添加购物车
// @Summary      添加购物车
// @Description  将商品添加到购物车
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        req body service.AddCartRequest true "添加请求"
// @Success      200 {object} utils.Response
// @Security     ApiKeyAuth
// @Router       /cart [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	err := h.cartSvc.AddItem(userID, &req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}

// UpdateQuantity 更新购物车商品数量
// @Summary      更新购物车商品数量
// @Description  更新购物车中指定商品的数量
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        id  path int true "购物车 ID"
// @Param        req body service.UpdateCartQuantityRequest true "更新请求"
// @Success      200 {object} utils.Response
// @Security     ApiKeyAuth
// @Router       /cart/{id} [put]
func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	userID := middleware.GetUserID(c)

	cartID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid cart id")
		return
	}

	var req service.UpdateCartQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	err = h.cartSvc.UpdateQuantity(userID, uint(cartID), req.Quantity)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}

// RemoveItem 删除购物车商品
// @Summary      删除购物车商品
// @Description  从购物车中移除指定商品
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        id path int true "购物车 ID"
// @Success      200 {object} utils.Response
// @Security     ApiKeyAuth
// @Router       /cart/{id} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID := middleware.GetUserID(c)

	cartID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid cart id")
		return
	}

	err = h.cartSvc.RemoveItem(userID, uint(cartID))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}

// GetCart 获取购物车
// @Summary      获取购物车
// @Description  获取当前用户的购物车及商品列表
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response{data=service.CartResponse}
// @Security     ApiKeyAuth
// @Router       /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.cartSvc.GetCart(userID)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// ClearCart 清空购物车
// @Summary      清空购物车
// @Description  清空当前用户的购物车
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response
// @Security     ApiKeyAuth
// @Router       /cart/clear [post]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	err := h.cartSvc.ClearCart(userID)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}
