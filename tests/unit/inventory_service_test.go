package unit

import (
	"avito-shop/models"
	"avito-shop/services"
	"avito-shop/tests/unit/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInventoryService_AddItem(t *testing.T) {
	// Создание мока репозитория
	mockRepo := mocks.NewInventoryRepositoryInterface(t)

	// Создание сервиса с моком репозитория
	service := services.InventoryService{
		InventoryRepo: mockRepo,
	}

	t.Run("AddItem success", func(t *testing.T) {
		// Настройка ожиданий мока
		mockRepo.On("AddItem", "testuser", "cup", 1).Return(nil)

		// Вызов метода сервиса
		err := service.AddItem("testuser", "cup", 1)

		// Проверка результатов
		assert.NoError(t, err)

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "AddItem", "testuser", "cup", 1)
	})

	t.Run("AddItem error", func(t *testing.T) {
		// Настройка ожиданий мока
		mockRepo.On("AddItem", "testuser", "invalid-item", 1).Return(errors.New("item not found"))

		// Вызов метода сервиса
		err := service.AddItem("testuser", "invalid-item", 1)

		// Проверка результатов
		assert.Error(t, err)
		assert.Equal(t, "item not found", err.Error())

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "AddItem", "testuser", "invalid-item", 1)
	})
}
func TestInventoryService_GetInventory(t *testing.T) {
	// Создание мока репозитория
	mockRepo := mocks.NewInventoryRepositoryInterface(t)

	// Создание сервиса с моком репозитория
	service := services.InventoryService{
		InventoryRepo: mockRepo,
	}

	t.Run("GetInventory success", func(t *testing.T) {
		// Настройка ожиданий мока
		inventory := []models.InventoryItem{
			{Type: "cup", Quantity: 1},
			{Type: "t-shirt", Quantity: 2},
		}
		mockRepo.On("GetInventory", "testuser").Return(inventory, nil)

		// Вызов метода сервиса
		result, err := service.GetInventory("testuser")

		// Проверка результатов
		assert.NoError(t, err)
		assert.Equal(t, inventory, result)

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "GetInventory", "testuser")
	})

	t.Run("GetInventory error", func(t *testing.T) {
		// Настройка ожиданий мока
		mockRepo.On("GetInventory", "invaliduser").Return(nil, errors.New("failed to fetch inventory"))

		// Вызов метода сервиса
		result, err := service.GetInventory("invaliduser")

		// Проверка результатов
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "failed to fetch inventory", err.Error())

		// Проверка, что мок был вызван
		mockRepo.AssertCalled(t, "GetInventory", "invaliduser")
	})
}
