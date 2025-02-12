package integration

import (
	"avito-shop/controllers"
	db2 "avito-shop/db"
	"avito-shop/middleware"
	"avito-shop/repositories"
	"avito-shop/services"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSendCoinIntegration(t *testing.T) {
	db2.InitDB()
	db := db2.GetDB()
	defer db2.CloseDB()

	// Создание репозиториев
	userRepo := repositories.UserRepository{DB: db}
	transactionRepo := repositories.TransactionRepository{DB: db}
	// Создание сервисов
	userService := services.UserService{UserRepo: &userRepo}
	transactionService := services.TransactionService{TransactionRepo: &transactionRepo}

	// Создание контроллера
	sendCoinController := controllers.NewSendCoinController(&userService, &transactionService)

	// Создание роутера
	r := mux.NewRouter()
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/sendCoin", sendCoinController.SendCoinHandler).Methods("POST")

	// Создание пользователей
	err := userService.CreateUser("sender", "password")
	if err != nil {
		return
	}
	err = userService.CreateUser("receiver", "password")
	if err != nil {
		return
	}

	t.Run("Successful coin transfer", func(t *testing.T) {
		// Генерация JWT-токена для отправителя
		token, err := GenerateJWT("sender")
		if err != nil {
			t.Fatalf("Failed to generate JWT token: %v", err)
		}

		// Создание HTTP-запроса
		payload := `{"toUser": "receiver", "amount": 100}`
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBufferString(payload))
		req.Header.Set("Authorization", ""+token)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		// Вызов роутера
		r.ServeHTTP(recorder, req)
		// Проверка статуса ответа
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200, got %d", recorder.Code)

		// Проверка баланса отправителя
		sender, _ := userService.GetUser("sender")
		assert.Equal(t, 900, sender.Coins, "Expected sender's coins to be 900, got %d", sender.Coins)

		// Проверка баланса получателя
		receiver, _ := userService.GetUser("receiver")
		assert.Equal(t, 1100, receiver.Coins, "Expected receiver's coins to be 1100, got %d", receiver.Coins)
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		// Генерация JWT-токена для отправителя
		token, err := GenerateJWT("sender")
		if err != nil {
			t.Fatalf("Failed to generate JWT token: %v", err)
		}

		// Создание HTTP-запроса
		payload := `{"toUser": "receiver", "amount": 2000}`
		req := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewBufferString(payload))
		req.Header.Set("Authorization", ""+token)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		// Вызов роутера
		r.ServeHTTP(recorder, req)

		// Проверка статуса ответа
		assert.Equal(t, http.StatusBadRequest, recorder.Code, "Expected status code 400, got %d", recorder.Code)
	})
	_, err = db.Exec("DELETE FROM users WHERE username='receiver';")
	if err != nil {
		return
	}
	_, err = db.Exec("DELETE FROM users WHERE username='sender';")
	if err != nil {
		return
	}
	_, err = db.Exec("DELETE FROM transactions WHERE from_user='sender'")
	if err != nil {
		return
	}
}
