package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type MerchantRegisterHandler struct {
	registerSvc *service.MerchantRegisterService
}

func NewMerchantRegisterHandler(registerSvc *service.MerchantRegisterService) *MerchantRegisterHandler {
	return &MerchantRegisterHandler{registerSvc: registerSvc}
}

// Register 商家注册
func (h *MerchantRegisterHandler) Register(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.RegisterMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.registerSvc.Register(userID, req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}
