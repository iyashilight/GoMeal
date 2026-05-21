package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type AddressHandler struct {
	addressSvc *service.AddressService
}

func NewAddressHandler(addressSvc *service.AddressService) *AddressHandler {
	return &AddressHandler{addressSvc: addressSvc}
}

// Create 创建收货地址
// @Summary      创建收货地址
// @Description  添加一个新的收货地址
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        req body service.CreateAddressRequest true "地址信息"
// @Success      200 {object} utils.Response{data=service.AddressResponse}
// @Security     ApiKeyAuth
// @Router       /addresses [post]
func (h *AddressHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.addressSvc.Create(userID, &req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// Update 更新收货地址
// @Summary      更新收货地址
// @Description  修改指定的收货地址
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        id  path int true "地址 ID"
// @Param        req body service.UpdateAddressRequest true "更新信息"
// @Success      200 {object} utils.Response{data=service.AddressResponse}
// @Security     ApiKeyAuth
// @Router       /addresses/{id} [put]
func (h *AddressHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)

	addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid address id")
		return
	}

	var req service.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.addressSvc.Update(userID, uint(addressID), &req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// Delete 删除收货地址
// @Summary      删除收货地址
// @Description  删除指定的收货地址
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        id path int true "地址 ID"
// @Success      200 {object} utils.Response
// @Security     ApiKeyAuth
// @Router       /addresses/{id} [delete]
func (h *AddressHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)

	addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid address id")
		return
	}

	err = h.addressSvc.Delete(userID, uint(addressID))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}

// GetAddresses 获取收货地址列表
// @Summary      获取收货地址列表
// @Description  获取当前用户的所有收货地址
// @Tags         address
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response{data=[]service.AddressResponse}
// @Security     ApiKeyAuth
// @Router       /addresses [get]
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.addressSvc.GetAddresses(userID)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// SetDefault 设置默认地址
// @Summary      设置默认地址
// @Description  将指定地址设为默认收货地址
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        id path int true "地址 ID"
// @Success      200 {object} utils.Response
// @Security     ApiKeyAuth
// @Router       /addresses/{id}/default [put]
func (h *AddressHandler) SetDefault(c *gin.Context) {
	userID := middleware.GetUserID(c)

	addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, utils.ErrInvalidParams, "invalid address id")
		return
	}

	err = h.addressSvc.SetDefault(userID, uint(addressID))
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, nil)
}
