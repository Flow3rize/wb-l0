package handlers

import (
	"log"
	"net/http"

	"github.com/flowerize/wb-l0/internal/models"
	"github.com/gin-gonic/gin"
)

type Cache interface {
	Get(orderUID string) (models.Order, bool)
	Set(orderUID string, order models.Order)
	LoadFromDB(orders []models.Order)
}

type Storage interface {
	GetOrder(orderUID string) (*models.Order, error)
	GetAllOrders() ([]models.Order, error)
	SaveOrder(order *models.Order) error
}

// Валидация входящих данных
type OrderInput struct {
	OrderUID    string `json:"order_uid" binding:"required"`
	TrackNumber string `json:"track_number" binding:"required"`
	Entry       string `json:"entry" binding:"required"`

	Delivery *models.Delivery `json:"delivery" binding:"required"`
	Payment  *models.Payment  `json:"payment" binding:"required"`
	Items    []models.Item    `json:"items" binding:"required,dive,required"`
}

type OrderHandler struct {
	cache Cache
	db    Storage
}

func NewOrderHandler(cache Cache, db Storage) *OrderHandler {
	return &OrderHandler{
		cache: cache,
		db:    db,
	}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderUID := c.Param("order_uid")

	cachedOrder, exists := h.cache.Get(orderUID)
	if exists {
		log.Printf("Данные взяты из кэша: %s", orderUID)
		c.JSON(http.StatusOK, cachedOrder)
		return
	}

	log.Printf("Данные взяты из БД: %s", orderUID)

	order, err := h.db.GetOrder(orderUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Заказ не найден"})
		return
	}

	h.cache.Set(orderUID, *order)

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input OrderInput
	if err := c.ShouldBindJSON(&input); err != nil { //использую ShouldBindJSON из пакета validator
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := &models.Order{
		OrderUID:    input.OrderUID,
		TrackNumber: input.TrackNumber,
		Entry:       input.Entry,
	}

	if err := h.db.SaveOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении заказа"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func StartServer(addr string, cache Cache, db Storage) error {
	r := gin.Default()

	orderHandler := NewOrderHandler(cache, db)
	r.GET("/orders/:order_uid", orderHandler.GetOrder)

	return r.Run(addr)
}
