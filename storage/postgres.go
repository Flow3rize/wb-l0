package storage

import (
	"encoding/json"
	"fmt"

	"github.com/flowerize/wl-l0/config"
	"github.com/flowerize/wl-l0/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type PostgresStorage struct {
	DB *pg.DB
}

func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	opt, _ := pg.ParseURL(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	))

	db := pg.Connect(opt)

	// Автомиграция
	model := &models.Order{}
	err := db.Model(model).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{DB: db}, nil
}

func (s *PostgresStorage) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	err := s.DB.Model(&orders).Select()
	if err != nil {
		return nil, err
	}

	for i := range orders {
		delivery := models.Delivery{}
		if err := json.Unmarshal(orders[i].DeliveryJSON, &delivery); err != nil {
			return nil, fmt.Errorf("ошибка десериализации Delivery: %w", err)
		}
		orders[i].Delivery = delivery

		payment := models.Payment{}
		if err := json.Unmarshal(orders[i].PaymentJSON, &payment); err != nil {
			return nil, fmt.Errorf("ошибка десериализации Payment: %w", err)
		}
		orders[i].Payment = payment

		items := []models.Item{}
		if err := json.Unmarshal(orders[i].ItemsJSON, &items); err != nil {
			return nil, fmt.Errorf("ошибка десериализации Items: %w", err)
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (s *PostgresStorage) SaveOrder(order *models.Order) error {
	_, err := s.DB.Model(order).Insert(order)
	return err
}

func (s *PostgresStorage) GetOrder(orderUID string) (*models.Order, error) {
	var order models.Order

	err := s.DB.Model(&order).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		return nil, err
	}

	delivery := models.Delivery{}
	if err := json.Unmarshal(order.DeliveryJSON, &delivery); err != nil {
		return nil, fmt.Errorf("ошибка десериализации Delivery: %w", err)
	}
	order.Delivery = delivery

	payment := models.Payment{}
	if err := json.Unmarshal(order.PaymentJSON, &payment); err != nil {
		return nil, fmt.Errorf("ошибка десериализации Payment: %w", err)
	}
	order.Payment = payment

	items := []models.Item{}
	if err := json.Unmarshal(order.ItemsJSON, &items); err != nil {
		return nil, fmt.Errorf("ошибка десериализации Items: %w", err)
	}
	order.Items = items

	return &order, nil
}
