package handler

import (
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/iyashilight/GoMeal/internal/middleware"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes 配置所有 HTTP 路由
// 公开路由：健康检查、Swagger 文档、注册登录、商家查询（限流 30次/分钟）
// 认证路由：用户信息、购物车、地址、订单操作（需 JWT + 限流 60次/分钟）
func SetupRoutes(
	r *gin.Engine,
	userH *UserHandler,
	merchantH *MerchantHandler,
	cartH *CartHandler,
	addrH *AddressHandler,
	orderH *OrderHandler,
	merchantManageH *MerchantManageHandler,
	paymentH *PaymentHandler,
	merchantRegisterH *MerchantRegisterHandler,
	jwtSecret string,
	rdb *redis.Client,
) {
	r.Use(middleware.AccessLog(), gin.Recovery())

	pprof.Register(r)

	r.GET("/health", health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 公开接口限流：30次/分钟
	publicLimit := middleware.RateLimit(rdb, 30, time.Minute)
	r.POST("/auth/register", publicLimit, userH.Register)
	r.POST("/auth/login", publicLimit, userH.Login)

	r.GET("/merchants", publicLimit, merchantH.GetMerchantList)
	r.GET("/merchants/:id", publicLimit, merchantH.GetMerchantDetail)
	r.GET("/foods/:id", publicLimit, merchantH.GetFoodDetail)
	r.POST("/payment/notify", publicLimit, paymentH.Notify)

	// 认证接口限流：60次/分钟
	authLimit := middleware.RateLimit(rdb, 60, time.Minute)
	auth := r.Group("")
	auth.Use(middleware.JWTAuth(jwtSecret), authLimit)
	{
		auth.GET("/user/info", userH.GetUserInfo)

		auth.GET("/cart", cartH.GetCart)
		auth.POST("/cart", cartH.AddItem)
		auth.PUT("/cart/:id", cartH.UpdateQuantity)
		auth.DELETE("/cart/:id", cartH.RemoveItem)
		auth.POST("/cart/clear", cartH.ClearCart)

		auth.GET("/addresses", addrH.GetAddresses)
		auth.POST("/addresses", addrH.Create)
		auth.PUT("/addresses/:id", addrH.Update)
		auth.DELETE("/addresses/:id", addrH.Delete)
		auth.PUT("/addresses/:id/default", addrH.SetDefault)

		auth.POST("/orders", orderH.CreateOrder)
		auth.GET("/orders", orderH.GetOrders)
		auth.GET("/orders/:id", orderH.GetOrderDetail)
		auth.PUT("/orders/:id/cancel", orderH.CancelOrder)
		auth.POST("/merchant/register", merchantRegisterH.Register)
		auth.POST("/orders/:id/pay", paymentH.PayOrder)
		auth.GET("/orders/:id/payment", paymentH.GetPayment)
		auth.GET("/merchant/foods", merchantManageH.GetMyFoods)
		auth.POST("/merchant/foods", merchantManageH.CreateFood)
		auth.PUT("/merchant/foods/:id", merchantManageH.UpdateFood)
		auth.DELETE("/merchant/foods/:id", merchantManageH.DeleteFood)
		auth.PUT("/merchant/foods/:id/status", merchantManageH.SetFoodStatus)
	}
}
