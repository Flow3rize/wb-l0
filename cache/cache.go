package cache

import (
	"fmt"
	"sync"

	"github.com/flowerize/wb-l0/models"
)

type InMemoryCache struct {
	data        map[string]models.Order
	lru         []string // Очередь LRU (от старого к новому)
	maxSize     int
	currentSize int
	mu          sync.RWMutex
}

func NewInMemoryCache(maxSize int) *InMemoryCache {
	return &InMemoryCache{
		data:        make(map[string]models.Order),
		lru:         make([]string, 0, maxSize),
		maxSize:     maxSize,
		currentSize: 0,
	}
}

func (c *InMemoryCache) LoadFromDB(orders []models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		if c.currentSize == c.maxSize {
			break
		}
		c.data[order.OrderUID] = order
		c.lru = append(c.lru, order.OrderUID)
		c.currentSize++
	}
}

func (c *InMemoryCache) Get(orderUID string) (models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	order, ok := c.data[orderUID]
	if !ok {
		return models.Order{}, false
	}

	c.moveToEnd(orderUID)
	return order, true
}

func (c *InMemoryCache) Set(orderUID string, order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.data[orderUID]; exists {
		// Обновляем данные, перемещаем ключ в конец очереди LRU
		c.data[orderUID] = order
		c.moveToEnd(orderUID)
		return
	}

	if c.currentSize == c.maxSize {
		oldest := c.lru[0]
		c.lru = c.lru[1:]
		delete(c.data, oldest)
		c.currentSize--
	}

	c.data[orderUID] = order
	c.lru = append(c.lru, orderUID)
	c.currentSize++
}

// moveToEnd перемещает ключ в конец LRU очереди (самый свежий)
func (c *InMemoryCache) moveToEnd(key string) {
	idx := -1
	for i, v := range c.lru {
		if v == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		c.lru = append(c.lru, key)
		return
	}
	c.lru = append(c.lru[:idx], c.lru[idx+1:]...)
	c.lru = append(c.lru, key)
}

func (c *InMemoryCache) Delete(orderUID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.data[orderUID]; !exists {
		return
	}

	delete(c.data, orderUID)
	c.currentSize--

	idx := -1
	for i, v := range c.lru {
		if v == orderUID {
			idx = i
			break
		}
	}
	if idx != -1 {
		c.lru = append(c.lru[:idx], c.lru[idx+1:]...)
	}
}

func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]models.Order)
	c.lru = c.lru[:0]
	c.currentSize = 0
}

func (c *InMemoryCache) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return fmt.Sprintf("InMemoryCache(size=%d, maxSize=%d, keys=%v)", c.currentSize, c.maxSize, c.lru)
}
