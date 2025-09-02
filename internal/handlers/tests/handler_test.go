package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cachemocks "github.com/flowerize/wb-l0/internal/cache/mocks"
	"github.com/flowerize/wb-l0/internal/handlers"
	"github.com/flowerize/wb-l0/internal/pkg/storage/mocks"

	"github.com/flowerize/wb-l0/internal/models"
	"github.com/gin-gonic/gin"
)

func TestGetOrder_CacheHit(t *testing.T) {
	cacheMock := cachemocks.NewMockCache()
	dbMock := mocks.NewMockStorage()

	orderUID := "test_order_uid"
	order := models.Order{OrderUID: orderUID}

	cacheMock.Set(orderUID, order)

	handler := handlers.NewOrderHandler(cacheMock, dbMock)

	r := gin.Default()
	r.GET("/orders/:order_uid", handler.GetOrder)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/orders/"+orderUID, nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	var response models.Order
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.OrderUID != orderUID {
		t.Errorf("Ожидался UID %s, получен %s", orderUID, response.OrderUID)
	}
}

func TestStartServer(t *testing.T) {
	cacheMock := cachemocks.NewMockCache()
	dbMock := mocks.NewMockStorage()
	orderUID := "order123"
	order := models.Order{OrderUID: orderUID}
	cacheMock.Set(orderUID, order)
	dbMock.SaveOrder(&order)

	errChan := make(chan error, 1)

	go func() {
		err := handlers.StartServer(":8081", cacheMock, dbMock)
		if err != nil {
			errChan <- err
		}
	}()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/orders/" + orderUID)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var gotOrder models.Order
	if err := json.NewDecoder(resp.Body).Decode(&gotOrder); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if gotOrder.OrderUID != orderUID {
		t.Errorf("Expected OrderUID %s, got %s", orderUID, gotOrder.OrderUID)
	}

	select {
	case err := <-errChan:
		t.Fatalf("Server error: %v", err)
	default:
	}
}
