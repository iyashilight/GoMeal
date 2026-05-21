package service

import (
	"errors"
	"testing"

	"github.com/iyashilight/GoMeal/internal/model"
)

var errNotFound = errors.New("not found")

type mockCartRepo struct {
	carts []model.Cart
	seq   uint
}

func (m *mockCartRepo) GetCartByUser(userID uint) ([]model.Cart, error) {
	var result []model.Cart
	for _, c := range m.carts {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockCartRepo) GetCartItem(userID, foodID uint) (*model.Cart, error) {
	for _, c := range m.carts {
		if c.UserID == userID && c.FoodID == foodID {
			return &c, nil
		}
	}
	return nil, errNotFound
}

func (m *mockCartRepo) Create(cart *model.Cart) error {
	m.seq++
	cart.ID = m.seq
	m.carts = append(m.carts, *cart)
	return nil
}

func (m *mockCartRepo) Update(cart *model.Cart) error {
	for i, c := range m.carts {
		if c.ID == cart.ID {
			m.carts[i] = *cart
			return nil
		}
	}
	return nil
}

func (m *mockCartRepo) Delete(userID, cartID uint) error {
	for i, c := range m.carts {
		if c.ID == cartID && c.UserID == userID {
			m.carts = append(m.carts[:i], m.carts[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockCartRepo) ClearByUser(userID uint) error {
	var kept []model.Cart
	for _, c := range m.carts {
		if c.UserID != userID {
			kept = append(kept, c)
		}
	}
	m.carts = kept
	return nil
}

func newMockCartRepo() *mockCartRepo {
	return &mockCartRepo{}
}

// cart 测试用的 mock foodRepo
type mockFoodRepoCart struct {
	foods map[uint]*model.Food
}

func (m *mockFoodRepoCart) FindFoodsByCategory(categoryID uint) ([]model.Food, error) {
	return nil, nil
}

func (m *mockFoodRepoCart) FindFoodByID(id uint) (*model.Food, error) {
	if f, ok := m.foods[id]; ok {
		return f, nil
	}
	return nil, errNotFound
}

func (m *mockFoodRepoCart) DecreaseStock(foodID uint, quantity int) error {
	return nil
}

func (m *mockFoodRepoCart) IncreaseStock(foodID uint, quantity int) error {
	return nil
}

func TestAddCartItem_Success(t *testing.T) {
	cartRepo := newMockCartRepo()
	foodRepo := &mockFoodRepoCart{
		foods: map[uint]*model.Food{
			1: {ID: 1, Name: "测试食品", Price: 10.0, MerchantID: 1},
		},
	}
	svc := NewCartService(cartRepo, foodRepo)

	err := svc.AddItem(1, &AddCartRequest{
		FoodID:     1,
		MerchantID: 1,
		Quantity:   2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	carts, _ := cartRepo.GetCartByUser(1)
	if len(carts) != 1 {
		t.Errorf("expected 1 cart item, got %d", len(carts))
	}
	if carts[0].Quantity != 2 {
		t.Errorf("expected quantity 2, got %d", carts[0].Quantity)
	}
}

func TestClearCart(t *testing.T) {
	cartRepo := newMockCartRepo()
	foodRepo := &mockFoodRepoCart{
		foods: map[uint]*model.Food{
			1: {ID: 1, Name: "食品1", Price: 10.0, MerchantID: 1},
		},
	}
	svc := NewCartService(cartRepo, foodRepo)

	svc.AddItem(1, &AddCartRequest{FoodID: 1, MerchantID: 1, Quantity: 1})

	err := svc.ClearCart(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	carts, _ := cartRepo.GetCartByUser(1)
	if len(carts) != 0 {
		t.Error("expected empty cart after clear")
	}
}

// 用于 mock 的错误标记
