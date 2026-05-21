package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type PaymentHandler struct {
	paymentSvc *service.PaymentService
}

func NewPaymentHandler(paymentSvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

// PayOrder 发起支付
func (h *PaymentHandler) PayOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid order id")
		return
	}

	resp, err := h.paymentSvc.PayOrder(userID, uint(orderID))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// GetPayment 查询支付状态
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	userID := middleware.GetUserID(c)

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid order id")
		return
	}

	resp, err := h.paymentSvc.GetPayment(userID, uint(orderID))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// Notify 支付回调（预留，供支付网关调用）
func (h *PaymentHandler) Notify(c *gin.Context) {
	utils.Success(c, gin.H{"message": "notify received"})
}
