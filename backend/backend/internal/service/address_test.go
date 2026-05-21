package service

import (
	"testing"

	"github.com/iyashilight/GoMeal/internal/model"
	"gorm.io/gorm"
)

type mockAddressRepo struct {
	addresses map[uint]*model.Address
	seq       uint
}

func (m *mockAddressRepo) Create(addr *model.Address) error {
	m.seq++
	addr.ID = m.seq
	m.addresses[addr.ID] = addr
	return nil
}

func (m *mockAddressRepo) FindByUser(userID uint) ([]model.Address, error) {
	var result []model.Address
	for _, a := range m.addresses {
		if a.UserID == userID {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (m *mockAddressRepo) FindByID(id uint) (*model.Address, error) {
	if a, ok := m.addresses[id]; ok {
		return a, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockAddressRepo) Update(addr *model.Address) error {
	m.addresses[addr.ID] = addr
	return nil
}

func (m *mockAddressRepo) Delete(userID, addressID uint) error {
	if a, ok := m.addresses[addressID]; ok && a.UserID == userID {
		delete(m.addresses, addressID)
	}
	return nil
}

func (m *mockAddressRepo) ClearDefault(userID uint) error {
	for _, a := range m.addresses {
		if a.UserID == userID {
			a.IsDefault = false
		}
	}
	return nil
}

func newMockAddressRepo() *mockAddressRepo {
	return &mockAddressRepo{addresses: make(map[uint]*model.Address)}
}

func TestCreateAddress_Success(t *testing.T) {
	repo := newMockAddressRepo()
	svc := NewAddressService(repo)

	resp, err := svc.Create(1, &CreateAddressRequest{
		Name:    "张三",
		Phone:   "13800138000",
		Address: "北京市朝阳区xxx路",
		Tag:     "家",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != "张三" {
		t.Errorf("expected name 张三, got %s", resp.Name)
	}
	if resp.Address != "北京市朝阳区xxx路" {
		t.Errorf("expected address 北京市朝阳区xxx路, got %s", resp.Address)
	}
}

func TestCreateAddress_Default(t *testing.T) {
	repo := newMockAddressRepo()
	svc := NewAddressService(repo)

	resp, err := svc.Create(1, &CreateAddressRequest{
		Name:      "家",
		Phone:     "13800138000",
		Address:   "北京市",
		IsDefault: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.IsDefault {
		t.Error("expected address to be default")
	}
}

func TestGetAddresses(t *testing.T) {
	repo := newMockAddressRepo()
	svc := NewAddressService(repo)

	svc.Create(1, &CreateAddressRequest{Name: "家", Phone: "13800138000", Address: "北京"})
	svc.Create(1, &CreateAddressRequest{Name: "公司", Phone: "13800138000", Address: "上海"})

	addresses, err := svc.GetAddresses(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addresses) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(addresses))
	}
}

func TestDeleteAddress(t *testing.T) {
	repo := newMockAddressRepo()
	svc := NewAddressService(repo)

	resp, _ := svc.Create(1, &CreateAddressRequest{
		Name: "家", Phone: "13800138000", Address: "北京",
	})

	err := svc.Delete(1, resp.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	addresses, _ := svc.GetAddresses(1)
	if len(addresses) != 0 {
		t.Error("expected empty addresses after delete")
	}
}
