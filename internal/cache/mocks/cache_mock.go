package cache

import (
	"sync"

	"github.com/flowerize/wb-l0/internal/models"
)

type CacheMock struct {
	data map[string]models.Order
	mu   sync.RWMutex
}

func NewMockCache() *CacheMock {
	return &CacheMock{
		data: make(map[string]models.Order),
	}
}

func (m *CacheMock) Get(orderUID string) (models.Order, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	order, ok := m.data[orderUID]
	return order, ok
}

func (m *CacheMock) Set(orderUID string, order models.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[orderUID] = order
}

func (m *CacheMock) LoadFromDB(orders []models.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, order := range orders {
		m.data[order.OrderUID] = order
	}
}
