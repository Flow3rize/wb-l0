package storage

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/flowerize/wb-l0/cache"
	"github.com/flowerize/wb-l0/internal/models"
)

type KafkaConsumer struct {
	brokers string
	topic   string
	db      *PostgresStorage
	cache   *cache.InMemoryCache
}

func NewKafkaConsumer(brokers, topic string, db *PostgresStorage, cache *cache.InMemoryCache) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topic:   topic,
		db:      db,
		cache:   cache,
	}
}

func (c *KafkaConsumer) Start() error {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": c.brokers,
		"group.id":          "order-service-group",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(conf)
	if err != nil {
		return fmt.Errorf("не удалось создать консьюмера: %w", err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		return fmt.Errorf("не удалось подписаться на топик: %w", err)
	}

	for {
		msg, err := consumer.ReadMessage(-1) // -1 = блокирующее ожидание
		if err != nil {
			fmt.Printf("Ошибка чтения сообщения: %v\n", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			fmt.Printf("Ошибка парсинга сообщения: %v\n", err)
			continue
		}

		if err := c.db.SaveOrder(&order); err != nil {
			fmt.Printf("Ошибка сохранения в БД: %v\n", err)
			continue
		}

		c.cache.Set(order.OrderUID, order)
	}
}
