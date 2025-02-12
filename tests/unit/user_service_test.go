package unit

import (
	"avito-shop/services"
	"avito-shop/tests/unit/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryInterface)

	// Настройка ожидаемого поведения мока
	mockRepo.On("UserExists", "testuser").Return(false, nil)    // Пользователь не существует
	mockRepo.On("UserExists", "existinguser").Return(true, nil) // Пользователь существует
	mockRepo.On("CreateUser", "testuser", "password").Return(nil)

	// Инициализация UserService с моком
	userService := services.UserService{UserRepo: mockRepo}

	t.Run("Create new user", func(t *testing.T) {
		err := userService.CreateUser("testuser", "password")
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "UserExists", "testuser")             // Убедитесь, что UserExists был вызван
		mockRepo.AssertCalled(t, "CreateUser", "testuser", "password") // Убедитесь, что CreateUser был вызван
	})

	t.Run("Fail to create existing user", func(t *testing.T) {
		err := userService.CreateUser("existinguser", "password")
		assert.EqualError(t, err, "user already exists")                      // Ожидаем ошибку "user already exists"
		mockRepo.AssertCalled(t, "UserExists", "existinguser")                // Убедитесь, что UserExists был вызван
		mockRepo.AssertNotCalled(t, "CreateUser", "existinguser", "password") // Убедитесь, что CreateUser НЕ вызывался
	})
}
