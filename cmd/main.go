package main

import (
	"log"

	"github.com/flowerize/wb-l0/cache"
	"github.com/flowerize/wb-l0/internal/config"
	"github.com/flowerize/wb-l0/internal/handlers"
	"github.com/flowerize/wb-l0/internal/pkg/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := storage.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	orderCache := cache.NewInMemoryCache(cfg.CacheSize)

	orders, err := db.GetAllOrders()
	if err != nil {
		log.Printf("Предупреждение: не удалось загрузить кеш из БД: %v", err)
	} else {
		orderCache.LoadFromDB(orders)
	}

	kafkaConsumer := storage.NewKafkaConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		db,
		orderCache,
	)
	go func() {
		if err := kafkaConsumer.Start(); err != nil {
			log.Fatalf("Ошибка запуска Kafka-потребителя: %v", err)
		}
	}()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	orderHandler := handlers.NewOrderHandler(orderCache, db)
	r.GET("/orders/:order_uid", orderHandler.GetOrder)

	// Подключение статики
	r.StaticFile("/", "./static/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
