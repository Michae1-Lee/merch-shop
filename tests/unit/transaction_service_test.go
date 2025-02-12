package unit

import (
	"errors"
	"testing"

	"avito-shop/models"
	"avito-shop/services"
	"avito-shop/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService(t *testing.T) {
	// Создание мока репозитория
	mockRepo := mocks.NewTransactionRepositoryInterface(t)

	// Создание сервиса с моком репозитория
	service := services.TransactionService{
		TransactionRepo: mockRepo,
	}

	t.Run("GetTransactionsForUser success", func(t *testing.T) {
		// Настройка ожиданий мока
		received := []models.TransactionHistory{
			{FromUser: "sender1", ToUser: "testuser", Amount: 100},
			{FromUser: "sender2", ToUser: "testuser", Amount: 200},
		}
		sent := []models.TransactionHistory{
			{FromUser: "testuser", ToUser: "receiver1", Amount: 50},
			{FromUser: "testuser", ToUser: "receiver2", Amount: 75},
		}
		mockRepo.On("GetTransactionsForUser", "testuser").Return(received, sent, nil)

		// Вызов метода сервиса
		history, err := service.GetTransactionsForUser("testuser")

		// Проверка результатов
		assert.NoError(t, err)
		assert.Equal(t, received, history.Received)
		assert.Equal(t, sent, history.Sent)

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "GetTransactionsForUser", "testuser")
	})

	t.Run("GetTransactionsForUser error", func(t *testing.T) {
		// Настройка ожиданий мока
		mockRepo.On("GetTransactionsForUser", "invaliduser").Return(nil, nil, errors.New("user not found"))

		// Вызов метода сервиса
		history, err := service.GetTransactionsForUser("invaliduser")

		// Проверка результатов
		assert.Error(t, err)
		assert.Equal(t, models.CoinHistory{}, history)
		assert.Equal(t, "user not found", err.Error())

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "GetTransactionsForUser", "invaliduser")
	})

	t.Run("CreateTransaction success", func(t *testing.T) {
		// Настройка ожиданий мока
		mockRepo.On("CreateTransaction", "sender", "receiver", 100).Return(nil)

		// Вызов метода сервиса
		err := service.CreateTransaction("sender", "receiver", 100)

		// Проверка результатов
		assert.NoError(t, err)

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "CreateTransaction", "sender", "receiver", 100)
	})

	t.Run("CreateTransaction error", func(t *testing.T) {
		// Настройка ожиданий мока
		mockRepo.On("CreateTransaction", "sender", "receiver", -100).Return(errors.New("invalid amount"))

		// Вызов метода сервиса
		err := service.CreateTransaction("sender", "receiver", -100)

		// Проверка результатов
		assert.Error(t, err)
		assert.Equal(t, "invalid amount", err.Error())

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "CreateTransaction", "sender", "receiver", -100)
	})
}
