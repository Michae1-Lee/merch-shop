package integration

import (
	db2 "avito-shop/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"avito-shop/controllers"
	"avito-shop/repositories"
	"avito-shop/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationAndAuthorization(t *testing.T) {
	db2.InitDB()
	db := db2.GetDB()
	defer db2.CloseDB()

	r := mux.NewRouter()

	// Инициализация репозитория, сервиса и контроллера
	userRepo := repositories.UserRepository{DB: db}
	userService := services.UserService{UserRepo: &userRepo}
	userController := controllers.NewAuthController(&userService)

	// Настройка маршрутов
	r.HandleFunc("/api/auth", userController.AuthHandler).Methods("POST")
	r.HandleFunc("/api/auth", userController.AuthHandler).Methods("POST")

	t.Run("Successful registration and login", func(t *testing.T) {
		// Регистрация пользователя
		registerPayload := `{"username": "testuser", "password": "password"}`
		registerReq := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(registerPayload))
		registerReq.Header.Set("Content-Type", "application/json")
		registerRecorder := httptest.NewRecorder()

		r.ServeHTTP(registerRecorder, registerReq)

		// Проверка статуса ответа при регистрации
		assert.Equal(t, http.StatusCreated, registerRecorder.Code)

		// Авторизация пользователя
		loginPayload := `{"username": "testuser", "password": "password"}`
		loginReq := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(loginPayload))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRecorder := httptest.NewRecorder()

		r.ServeHTTP(loginRecorder, loginReq)

		// Проверка статуса ответа при авторизации
		assert.Equal(t, http.StatusOK, loginRecorder.Code)

		// Проверка, что токен возвращается
		var loginResponse map[string]string
		err := json.Unmarshal(loginRecorder.Body.Bytes(), &loginResponse)
		assert.NoError(t, err)
		assert.Contains(t, loginResponse, "token")
	})

	t.Run("Login with invalid credentials", func(t *testing.T) {
		// Авторизация с неверными данными
		loginPayload := `{"username": "testuser", "password": "wrongpassword"}`
		loginReq := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewBufferString(loginPayload))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRecorder := httptest.NewRecorder()

		r.ServeHTTP(loginRecorder, loginReq)

		// Проверка статуса ответа
		assert.Equal(t, http.StatusUnauthorized, loginRecorder.Code)
	})
	_, err := db.Exec("DELETE FROM users WHERE username='testuser'")
	if err != nil {
		return
	}
}
