package service

import (
	"errors"
	"log/slog"

	"github.com/iyashilight/GoMeal/internal/config"
	"github.com/iyashilight/GoMeal/internal/model"
	"github.com/iyashilight/GoMeal/internal/repository"
	"github.com/iyashilight/GoMeal/internal/utils"
	"gorm.io/gorm"
)

// UserService 用户业务，处理注册、登录、信息查询
type UserService struct {
	userRepo     repository.UserRepository
	merchantRepo repository.MerchantRepository
	cfg          *config.Config
}

// NewUserService 创建用户服务
func NewUserService(
	userRepo repository.UserRepository,
	merchantRepo repository.MerchantRepository,
	cfg *config.Config,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		merchantRepo: merchantRepo,
		cfg:          cfg,
	}
}

// RegisterRequest 注册请求参数，binding 标签用于 Gin 参数校验
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,len=11,numeric"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID         uint   `json:"id"`
	Phone      string `json:"phone"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	UserType   int    `json:"user_type"`
	MerchantID uint   `json:"merchant_id"` // 0=非商家用户, >0=商家用户
}

// RegisterResponse 注册响应，包含用户信息和 JWT Token
type RegisterResponse struct {
	User  *UserInfoResponse `json:"user"`
	Token string            `json:"token"`
}

// LoginResponse 登录响应，包含用户信息和 JWT Token
type LoginResponse struct {
	User  *UserInfoResponse `json:"user"`
	Token string            `json:"token"`
}

// Register 用户注册
// 1. 检查手机号是否已注册
// 2. 密码 bcrypt 加密
// 3. 创建用户，生成 JWT Token
// 4. 默认昵称：用户 + 手机号后4位
func (s *UserService) Register(req RegisterRequest) (*RegisterResponse, error) {
	existing, err := s.userRepo.FindByPhone(req.Phone)
	if err == nil {
		_ = existing
		slog.Warn("register failed: phone already exists", "phone", req.Phone)
		return nil, ErrPhoneExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("register failed: db error", "phone", req.Phone, "error", err)
		return nil, err
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		slog.Error("register failed: hash password error", "error", err)
		return nil, err
	}
	user := &model.User{
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: "用户" + req.Phone[len(req.Phone)-4:],
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Phone, user.UserType,
		s.cfg.JWT.SecretKey, s.cfg.JWT.ExpireHours)
	if err != nil {
		slog.Error("register failed: generate token error", "user_id", user.ID, "error", err)
		return nil, err
	}
	slog.Info("user registered", "user_id", user.ID, "phone", user.Phone)
	return &RegisterResponse{
		User: &UserInfoResponse{
			ID:       user.ID,
			Phone:    user.Phone,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			UserType: user.UserType,
		},
		Token: token,
	}, nil
}

// Login 用户登录
// 1. 根据手机号查找用户
// 2. 验证密码
// 3. 生成 JWT Token
// 注意：手机号不存在和密码错误返回相同错误，防止枚举
func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		slog.Warn("login failed: user not found", "phone", req.Phone)
		return nil, ErrWrongPassword
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		slog.Warn("login failed: wrong password", "phone", req.Phone)
		return nil, ErrWrongPassword
	}
	token, err := utils.GenerateToken(user.ID, user.Phone, user.UserType,
		s.cfg.JWT.SecretKey, s.cfg.JWT.ExpireHours)
	if err != nil {
		slog.Error("login failed: generate token error", "user_id", user.ID, "error", err)
		return nil, err
	}
	slog.Info("user logged in", "user_id", user.ID, "phone", user.Phone)
	return &LoginResponse{
		User: &UserInfoResponse{
			ID:       user.ID,
			Phone:    user.Phone,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			UserType: user.UserType,
		},
		Token: token,
	}, nil
}

// GetUserInfo 获取当前登录用户的信息
func (s *UserService) GetUserInfo(userID uint) (*UserInfoResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	resp := &UserInfoResponse{
		ID:       user.ID,
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		UserType: user.UserType,
	}
	// 查询商家信息
	merchant, err := s.merchantRepo.FindByUserID(userID)
	if err == nil {
		resp.MerchantID = merchant.ID
	}
	return resp, nil
}
