// Package config 负责配置文件的加载和解析，以及数据库、Redis 的初始化
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 顶层配置结构，对应 config.yaml
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig HTTP 服务配置
type ServerConfig struct {
	Port int `yaml:"port"` // 监听端口，默认 8080
}

// DatabaseConfig MySQL 数据库连接配置
type DatabaseConfig struct {
	User     string `yaml:"user"`     // 数据库用户名
	Password string `yaml:"password"` // 数据库密码
	Host     string `yaml:"host"`     // 主机地址，默认 127.0.0.1
	Port     int    `yaml:"port"`     // 端口，默认 3306
	DBName   string `yaml:"dbname"`   // 数据库名
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`     // 地址，默认 127.0.0.1:6379
	Password string `yaml:"password"` // 密码
	DB       int    `yaml:"db"`       // 数据库编号
}

// JWTConfig JWT 令牌配置
type JWTConfig struct {
	SecretKey   string `yaml:"secret_key"`   // 签名密钥
	ExpireHours int    `yaml:"expire_hours"` // 过期时间（小时），默认 72
}

// setDefaults 为零值字段设置默认值
func (c *Config) setDefaults() {
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Database.Host == "" {
		c.Database.Host = "127.0.0.1"
	}
	if c.Database.Port == 0 {
		c.Database.Port = 3306
	}
	if c.Redis.Addr == "" {
		c.Redis.Addr = "127.0.0.1:6379"
	}
	if c.JWT.ExpireHours == 0 {
		c.JWT.ExpireHours = 72
	}
}

// DSN 生成 MySQL 连接字符串
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}

// Load 从 YAML 文件加载配置，解析后设置默认值
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file:%w", err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config:%w", err)
	}
	cfg.setDefaults()
	cfg.overrideFromEnv()
	return &cfg, nil
}

// overrideFromEnv 用环境变量覆盖配置中的敏感字段
func (c *Config) overrideFromEnv() {
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.Database.Password = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		c.JWT.SecretKey = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		c.Redis.Password = v
	}
}
