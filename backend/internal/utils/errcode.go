package utils

// 统一业务错误码定义
// 规范：1xxx 通用错误，2xxx 用户模块，3xxx 资源未找到，4xxx 业务逻辑错误
const (
	CodeSuccess = 0

	// 通用错误 1xxx
	ErrInvalidParams = 1001 // 请求参数校验失败
	ErrUnauthorized  = 1002 // 未认证（未登录或 token 无效）
	ErrAccessDenied  = 1003 // 无权限（角色不足）

	// 用户模块错误 2xxx
	ErrPhoneExists   = 2001 // 手机号已注册
	ErrWrongPassword = 2002 // 密码错误
	ErrUserNotFound  = 2003 // 用户不存在

	// 资源未找到错误 3xxx
	ErrMerchantNotFound = 3001 // 商家不存在
	ErrFoodNotFound     = 3002 // 商品不存在
	ErrAddressNotFound  = 3003 // 地址不存在
	ErrOrderNotFound    = 3004 // 订单不存在
	ErrCartNotFound     = 3005 // 购物车条目不存在
	ErrCartEmpty        = 3006 // 购物车为空

	// 业务逻辑错误 4xxx
	ErrInvalidQuantity    = 4001 // 无效数量
	ErrInsufficientStock  = 4002 // 库存不足
	ErrInvalidOrderStatus = 4003 // 非法订单状态变更
	ErrMerchantConflict   = 4004 // 购物车跨商家冲突

	// 商家管理错误 5xxx
	ErrNotMerchant         = 5001 // 用户不是商家
	ErrFoodNotBelong       = 5002 // 商品不属于该商家
	ErrCategoryNotFound    = 5003 // 分类不存在
	ErrMerchantNameExists  = 5004 // 商家名称已存在

	// 支付错误 6xxx
	ErrPaymentNotFound = 6001 // 支付记录不存在
)
