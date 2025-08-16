package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/flowerize/wb-l0/cache"
	"github.com/flowerize/wb-l0/config"
	"github.com/flowerize/wb-l0/handlers"
	"github.com/flowerize/wb-l0/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func waitForPostgres(dsn string, maxRetries int, delay time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		log.Printf("Попытка подключения к PostgreSQL: %d из %d", i+1, maxRetries)

		conn, err := pgx.Connect(context.Background(), dsn)
		if err == nil {
			conn.Close(context.Background())
			log.Println("Подключение к PostgreSQL успешно!")
			return nil
		}

		log.Printf("Не удалось подключиться к PostgreSQL: %v", err)
		time.Sleep(delay)
	}

	return fmt.Errorf("не удалось подключиться к PostgreSQL после %d попыток", maxRetries)
}

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	if err := waitForPostgres(dsn, 30, 5*time.Second); err != nil {
		log.Fatalf("Ошибка ожидания PostgreSQL: %v", err)
	}

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
	r.GET("/orders/:id", orderHandler.GetOrder)

	// Подключение статики
	r.StaticFile("/", "./static/index.html")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
