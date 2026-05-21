package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type OrderHandler struct {
	orderSvc *service.OrderService
}

func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

// CreateOrder 创建订单
// @Summary      创建订单
// @Description  提交当前购物车中的商品创建订单
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        req body service.CreateOrderRequest true "下单请求"
// @Success      200 {object} utils.Response{data=service.OrderResponse}
// @Security     ApiKeyAuth
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.orderSvc.CreateOrder(userID, &req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// GetOrders 获取订单列表
// @Summary      获取订单列表
// @Description  获取当前用户的订单列表，支持按状态筛选和分页
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        status query int false "订单状态" Enums(0,1,2,3,4,5,6)
// @Param        page   query int false "页码" default(1)
// @Param        size   query int false "每页数量" default(10)
// @Success      200 {object} utils.Response{data=service.OrderListResponse}
// @Security     ApiKeyAuth
// @Router       /orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	var pageReq utils.PageReq
	userID := middleware.GetUserID(c)
	status := -1
	if s := c.Query("status"); s != "" {
		status, _ = strconv.Atoi(s)
	}
	err := c.ShouldBindQuery(&pageReq)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}
	resp, err := h.orderSvc.GetOrders(userID, status, pageReq)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// GetOrderDetail 获取订单详情
// @Summary      获取订单详情
// @Description  获取指定订单的详细信息
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id path int true "订单 ID"
// @Success      200 {object} utils.Response{data=service.OrderResponse}
// @Security     ApiKeyAuth
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}
	resp, err := h.orderSvc.GetOrderDetail(userID, uint(id))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)

}

// CancelOrder 取消订单
// @Summary      取消订单
// @Description  取消指定订单（仅待支付状态可取消）
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id path int true "订单 ID"
// @Success      200 {object} utils.Response
// @Security     ApiKeyAuth
// @Router       /orders/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}
	err = h.orderSvc.CancelOrder(userID, uint(id))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}
