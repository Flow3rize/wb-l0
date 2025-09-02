package mocks

import (
	"sync"

	"github.com/flowerize/wb-l0/internal/models"
	"github.com/flowerize/wb-l0/internal/pkg/storage"
)

type MockStorage struct {
	data map[string]*models.Order
	mu   sync.RWMutex
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string]*models.Order),
	}
}

func (m *MockStorage) GetOrder(orderUID string) (*models.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if order, exists := m.data[orderUID]; exists {
		return order, nil
	}
	return nil, storage.ErrOrderNotFound
}

func (m *MockStorage) GetAllOrders() ([]models.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var orders []models.Order
	for _, order := range m.data {
		orders = append(orders, *order)
	}
	return orders, nil
}
func (m *MockStorage) SaveOrder(order *models.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[order.OrderUID] = order
	return nil
}
