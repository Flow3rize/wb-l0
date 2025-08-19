package storage

import (
	"fmt"

	"github.com/flowerize/wb-l0/internal/config"
	"github.com/flowerize/wb-l0/internal/models"
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
	models := []interface{}{
		(*models.Order)(nil),
		(*models.Delivery)(nil),
		(*models.Payment)(nil),
		(*models.Item)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}

	return &PostgresStorage{DB: db}, nil
}

func (s *PostgresStorage) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	// Получаем все заказы
	err := s.DB.Model(&orders).Select()
	if err != nil {
		return nil, err
	}

	for i := range orders {
		var delivery models.Delivery
		err = s.DB.Model(&delivery).
			Where("order_uid = ?", orders[i].OrderUID).
			Select()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения доставки: %w", err)
		}
		orders[i].Delivery = &delivery

		var payment models.Payment
		err = s.DB.Model(&payment).
			Where("order_uid = ?", orders[i].OrderUID).
			Select()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения оплаты: %w", err)
		}
		orders[i].Payment = &payment

		var items []models.Item
		err = s.DB.Model(&items).
			Where("order_uid = ?", orders[i].OrderUID).
			Select()
		if err != nil {
			return nil, fmt.Errorf("ошибка получения товаров: %w", err)
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (s *PostgresStorage) SaveOrder(order *models.Order) error {

	_, err := s.DB.Model(order).Insert(order)
	if err != nil {
		return err
	}

	delivery := order.Delivery
	delivery.OrderUID = order.OrderUID
	_, err = s.DB.Model(delivery).Insert(delivery)
	if err != nil {
		return err
	}

	payment := order.Payment
	payment.OrderUID = order.OrderUID
	_, err = s.DB.Model(payment).Insert(payment)
	if err != nil {
		return err
	}

	for i := range order.Items {
		item := &order.Items[i]
		item.OrderUID = order.OrderUID
		_, err = s.DB.Model(item).Insert(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStorage) GetOrder(orderUID string) (*models.Order, error) {
	var order models.Order

	err := s.DB.Model(&order).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		return nil, err
	}

	var delivery models.Delivery
	err = s.DB.Model(&delivery).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения доставки: %w", err)
	}
	order.Delivery = &delivery

	var payment models.Payment
	err = s.DB.Model(&payment).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения оплаты: %w", err)
	}
	order.Payment = &payment

	var items []models.Item
	err = s.DB.Model(&items).
		Where("order_uid = ?", orderUID).
		Select()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товаров: %w", err)
	}
	order.Items = items

	return &order, nil
}
