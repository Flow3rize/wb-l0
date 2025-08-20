package handlers

import (
	"net/http"

	"github.com/flowerize/wb-l0/cache"
	"github.com/flowerize/wb-l0/internal/models"
	"github.com/flowerize/wb-l0/internal/pkg/storage"
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
	orderUID := c.Param("order_uid")

	var order models.Order
	var delivery models.Delivery
	var payment models.Payment
	var items []models.Item

	err := h.db.DB.Model(&order).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}

	err = h.db.DB.Model(&delivery).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения доставки"})
		return
	}

	err = h.db.DB.Model(&payment).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения оплаты"})
		return
	}

	err = h.db.DB.Model(&items).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения товаров"})
		return
	}

	order.Delivery = &delivery
	order.Payment = &payment
	order.Items = items

	c.JSON(http.StatusOK, order)
}

func StartServer(addr string, cache *cache.InMemoryCache, db *storage.PostgresStorage) error {
	r := gin.Default()

	orderHandler := NewOrderHandler(cache, db)
	r.GET("/orders/:id", orderHandler.GetOrder)

	return r.Run(addr)
}
