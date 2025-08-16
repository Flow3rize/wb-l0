package handlers

import (
	"net/http"

	"github.com/flow3rize/wb-l0/cache"
	"github.com/flow3rize/wb-l0/storage"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	cache *cache.InMemoryCache
	db    *storage.PostgresStorage
}

func NewOrderHandler(cache *cache.InMemoryCache, db *storage.PostgresStorage) *OrderHandler {
	return &OrderHandler{
		cache: cache,
		db:    db,
	}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	if order, ok := h.cache.Get(orderID); ok {
		c.JSON(http.StatusOK, order)
		return
	}

	order, err := h.db.GetOrder(orderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	h.cache.Set(orderID, *order)

	c.JSON(http.StatusOK, order)
}

func StartServer(addr string, cache *cache.InMemoryCache, db *storage.PostgresStorage) error {
	r := gin.Default()

	orderHandler := NewOrderHandler(cache, db)
	r.GET("/orders/:id", orderHandler.GetOrder)

	return r.Run(addr)
}
