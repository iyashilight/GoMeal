package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/iyashilight/GoMeal/internal/service"
	"github.com/iyashilight/GoMeal/internal/utils"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// Register 用户注册
// @Summary      用户注册
// @Description  使用手机号和密码注册
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req body service.RegisterRequest true "注册请求"
// @Success      200 {object} utils.Response{data=service.RegisterResponse}
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.userSvc.Register(req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// Login 用户登录
// @Summary      用户登录
// @Description  使用手机号和密码登录
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req body service.LoginRequest true "登录请求"
// @Success      200 {object} utils.Response{data=service.LoginResponse}
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.ErrInvalidParams, utils.FormatValidationError(err))
		return
	}

	resp, err := h.userSvc.Login(req)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}

// GetUserInfo 获取用户信息
// @Summary      获取用户信息
// @Description  获取当前登录用户的个人信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response{data=service.UserInfoResponse}
// @Security     ApiKeyAuth
// @Router       /user/info [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.userSvc.GetUserInfo(userID)
	if err != nil {
		utils.ErrorFromErr(c, err)
		return
	}
	utils.Success(c, resp)
}
