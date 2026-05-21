package config

import (
	"fmt"
	"log/slog"

	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// autoMigrate 自动迁移所有模型到数据库表中
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Merchant{},
		&model.Category{},
		&model.Food{},
		&model.Address{},
		&model.Cart{},
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
	)
}

// InitDB 初始化数据库连接，创建连接池并执行自动迁移
func InitDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database:%w", err)
	}
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("auto migrate:%w", err)
	}
	slog.Info("database connected", "host", cfg.Database.Host, "port", cfg.Database.Port,
		"db", cfg.Database.DBName)
	return db, nil
}
