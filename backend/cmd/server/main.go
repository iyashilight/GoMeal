// Package main 应用入口，负责初始化所有依赖并启动 HTTP 服务
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/iyashilight/GoMeal/docs"
	"github.com/iyashilight/GoMeal/internal/cache"
	"github.com/iyashilight/GoMeal/internal/config"
	"github.com/iyashilight/GoMeal/internal/handler"
	"github.com/iyashilight/GoMeal/internal/repository"
	"github.com/iyashilight/GoMeal/internal/service"
)

func main() {
	// 1. 加载配置（支持 CONFIG_PATH 环境变量覆盖）
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config:%v", err)
	}

	// 2. 初始化数据库（含自动迁移）
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("init database:%v", err)
	}

	// 3. 设置结构化日志输出
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// 4. 初始化 Redis
	rdb, err := config.InitRedis(cfg)
	if err != nil {
		log.Fatalf("init redis:%v", err)
	}

	// 5. 初始化路由
	r := gin.New()

	// 7. 创建 Repository 层实例
	userRepo := repository.NewUserRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	cartRepo := repository.NewCartRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	foodRepo := repository.NewFoodRepository(db)
	merchantFoodRepo := repository.NewMerchantFoodRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// 8. 创建缓存实例
	cacheClient := cache.NewCache(rdb)
	stockCache := cache.NewStockCache(rdb)

	// 9. 创建 Service 层实例
	addressSvc := service.NewAddressService(addressRepo)
	cartSvc := service.NewCartService(cartRepo, foodRepo)
	merchantSvc := service.NewMerchantService(merchantRepo, categoryRepo, foodRepo, cacheClient, stockCache)
	transactor := repository.NewTransactor(db)
	orderSvc := service.NewOrderService(orderRepo, cartRepo, addressRepo, foodRepo, transactor, stockCache)
	userSvc := service.NewUserService(userRepo, merchantRepo, cfg)
	merchantManageSvc := service.NewMerchantManageService(merchantRepo, categoryRepo, merchantFoodRepo, cacheClient)
	paymentSvc := service.NewPaymentService(paymentRepo, orderRepo, transactor, cacheClient)
	merchantRegisterSvc := service.NewMerchantRegisterService(merchantRepo, categoryRepo, cacheClient)

	// 10. 创建 Handler 层实例
	addressH := handler.NewAddressHandler(addressSvc)
	cartH := handler.NewCartHandler(cartSvc)
	merchantH := handler.NewMerchantHandler(merchantSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	userH := handler.NewUserHandler(userSvc)
	merchantManageH := handler.NewMerchantManageHandler(merchantManageSvc)
	paymentH := handler.NewPaymentHandler(paymentSvc)
	merchantRegisterH := handler.NewMerchantRegisterHandler(merchantRegisterSvc)

	// 11. 注册路由并启动 HTTP 服务
	handler.SetupRoutes(r, userH, merchantH, cartH, addressH, orderH, merchantManageH, paymentH, merchantRegisterH, cfg.JWT.SecretKey, rdb)
	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
