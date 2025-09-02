package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	KafkaBrokers string
	KafkaTopic   string

	CacheSize int
}

func LoadConfig() *Config {
	return &Config{
		KafkaBrokers: getEnv("KAFKA_BROKERS", "kafka:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "orders"),
		DBHost:       getEnv("DB_HOST", "postgres"),
		DBPort:       getEnvInt("DB_PORT", 5432),
		DBUser:       getEnv("DB_USER", "flowerize"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "order_db"),
		CacheSize:    getEnvInt("CACHE_SIZE", 1000),
	}
}

// Вспомогательная функция для получения строки
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Вспомогательная функция для получения целого числа
func getEnvInt(key string, defaultValue int) int {
	if valueStr := os.Getenv(key); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
