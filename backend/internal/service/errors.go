// Package service 实现业务逻辑层，调用 Repository 接口完成核心业务操作
package service

import "github.com/iyashilight/GoMeal/internal/utils"

// 业务错误定义，统一使用 AppError 返回给 Handler 层处理
var (
	ErrPhoneExists        = utils.NewAppError(utils.ErrPhoneExists, "phone already registered")
	ErrWrongPassword      = utils.NewAppError(utils.ErrWrongPassword, "invalid phone or password")
	ErrUserNotFound       = utils.NewAppError(utils.ErrUserNotFound, "user not found")
	ErrMerchantNotFound   = utils.NewAppError(utils.ErrMerchantNotFound, "merchant not found")
	ErrFoodNotFound       = utils.NewAppError(utils.ErrFoodNotFound, "food not found")
	ErrAddressNotFound    = utils.NewAppError(utils.ErrAddressNotFound, "address not found")
	ErrCartEmpty          = utils.NewAppError(utils.ErrCartEmpty, "cart is empty")
	ErrOrderNotFound      = utils.NewAppError(utils.ErrOrderNotFound, "order not found")
	ErrInvalidQuantity    = utils.NewAppError(utils.ErrInvalidQuantity, "quantity must be at least 1")
	ErrInsufficientStock  = utils.NewAppError(utils.ErrInsufficientStock, "insufficient stock")
	ErrInvalidOrderStatus = utils.NewAppError(utils.ErrInvalidOrderStatus, "invalid order status")
	ErrCartNotFound       = utils.NewAppError(utils.ErrCartNotFound, "cart not found")
	ErrMerchantConflict   = utils.NewAppError(utils.ErrMerchantConflict, "Merchant Conflict")
	ErrNotMerchant        = utils.NewAppError(utils.ErrNotMerchant, "not a merchant")
	ErrFoodNotBelong      = utils.NewAppError(utils.ErrFoodNotBelong, "food does not belong to the merchant")
	ErrCategoryNotFound   = utils.NewAppError(utils.ErrCategoryNotFound, "category not found")
	ErrPaymentNotFound    = utils.NewAppError(utils.ErrPaymentNotFound, "payment not found")
)
