// Package repository 实现数据访问层，提供对数据库的 CRUD 操作
package repository

import (
	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓储接口，定义用户相关的数据库操作契约
type UserRepository interface {
	Create(user *model.User) error             // 创建用户
	FindByID(id uint) (*model.User, error)     // 按 ID 查询
	FindByPhone(phone string) (*model.User, error) // 按手机号查询
}

// userRepository 用户仓储的 GORM 实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 插入新用户记录
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID 按主键 ID 查询用户
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone 按手机号查询用户（手机号有唯一索引，只返回一条）
func (r *userRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone=?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
