package service

import (
	"log/slog"

	"github.com/iyashilight/GoMeal/internal/model"
	"github.com/iyashilight/GoMeal/internal/repository"
)

// AddressService 收货地址业务
type AddressService struct {
	addressRepo repository.AddressRepository
}

// NewAddressService 创建地址服务
func NewAddressService(
	addressRepo repository.AddressRepository,
) *AddressService {
	return &AddressService{
		addressRepo: addressRepo,
	}
}

// CreateAddressRequest 创建地址请求
type CreateAddressRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Phone     string `json:"phone" binding:"required,len=11,numeric"`
	Address   string `json:"address" binding:"required"`
	Tag       string `json:"tag"`
	IsDefault bool   `json:"is_default"`
}

// UpdateAddressRequest 更新地址请求（指针类型字段用于区分"不传"和"传false"）
type UpdateAddressRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Tag       string `json:"tag"`
	IsDefault *bool  `json:"is_default"`
}

// AddressResponse 地址响应
type AddressResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Tag       string `json:"tag"`
	IsDefault bool   `json:"is_default"`
}

// Create 创建收货地址
// 如果新地址设为默认，先清除用户其他地址的默认标记
// 如果用户还没有任何地址，自动设为默认
func (s *AddressService) Create(userID uint, req *CreateAddressRequest) (*AddressResponse, error) {
	if req.IsDefault {
		_ = s.addressRepo.ClearDefault(userID)
	}

	addr := &model.Address{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Address:   req.Address,
		Tag:       req.Tag,
		IsDefault: req.IsDefault,
	}

	if !req.IsDefault {
		existing, _ := s.addressRepo.FindByUser(userID)
		if len(existing) == 0 {
			addr.IsDefault = true
		}
	}

	if err := s.addressRepo.Create(addr); err != nil {
		return nil, err
	}
	slog.Info("address created", "address_id", addr.ID, "user_id", userID)
	return &AddressResponse{
		ID:        addr.ID,
		Name:      addr.Name,
		Phone:     addr.Phone,
		Address:   addr.Address,
		Tag:       addr.Tag,
		IsDefault: addr.IsDefault,
	}, nil
}

// Update 更新收货地址
// 校验地址归属（只能更新自己的地址）
// 只更新非空字段
func (s *AddressService) Update(userID, addressID uint, req *UpdateAddressRequest) (*AddressResponse, error) {
	addr, err := s.addressRepo.FindByID(addressID)
	if err != nil {
		return nil, ErrAddressNotFound
	}
	if addr.UserID != userID {
		return nil, ErrAddressNotFound
	}

	if req.Name != "" {
		addr.Name = req.Name
	}
	if req.Phone != "" {
		addr.Phone = req.Phone
	}
	if req.Address != "" {
		addr.Address = req.Address
	}
	if req.Tag != "" {
		addr.Tag = req.Tag
	}
	if req.IsDefault != nil && *req.IsDefault {
		_ = s.addressRepo.ClearDefault(userID)
		addr.IsDefault = true
	} else if req.IsDefault != nil && !*req.IsDefault {
		addr.IsDefault = false
	}

	if err := s.addressRepo.Update(addr); err != nil {
		return nil, err
	}
	slog.Info("address updated", "address_id", addr.ID, "user_id", userID)
	return &AddressResponse{
		ID:        addr.ID,
		Name:      addr.Name,
		Phone:     addr.Phone,
		Address:   addr.Address,
		Tag:       addr.Tag,
		IsDefault: addr.IsDefault,
	}, nil
}

// Delete 删除收货地址
func (s *AddressService) Delete(userID, addressID uint) error {
	err := s.addressRepo.Delete(userID, addressID)
	if err != nil {
		return err
	}
	slog.Info("address deleted", "address_id", addressID, "user_id", userID)
	return nil
}

// GetAddresses 获取用户所有收货地址
func (s *AddressService) GetAddresses(userID uint) ([]AddressResponse, error) {
	addresses, err := s.addressRepo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]AddressResponse, len(addresses))
	for i, a := range addresses {
		result[i] = AddressResponse{
			ID:        a.ID,
			Name:      a.Name,
			Phone:     a.Phone,
			Address:   a.Address,
			Tag:       a.Tag,
			IsDefault: a.IsDefault,
		}
	}
	return result, nil
}

// SetDefault 设置默认地址
// 清除用户其他地址的默认标记，再将指定地址设为默认
func (s *AddressService) SetDefault(userID, addressID uint) error {
	addr, err := s.addressRepo.FindByID(addressID)
	if err != nil {
		return ErrAddressNotFound
	}
	if addr.UserID != userID {
		return ErrAddressNotFound
	}

	if err := s.addressRepo.ClearDefault(userID); err != nil {
		return err
	}

	addr.IsDefault = true
	err = s.addressRepo.Update(addr)
	if err != nil {
		return err
	}
	slog.Info("address set default", "address_id", addressID, "user_id", userID)
	return nil
}
