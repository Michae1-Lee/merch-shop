package integration

import (
	db2 "avito-shop/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"avito-shop/controllers"
	"avito-shop/middleware"
	"avito-shop/repositories"
	"avito-shop/services"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		return "", err
	}
	return "Bearer " + tokenString, nil
}
func TestBuyItemIntegration(t *testing.T) {
	db2.InitDB()
	db := db2.GetDB()
	defer db2.CloseDB()

	// Создание репозиториев
	merchRepo := repositories.MerchRepository{DB: db}
	userRepo := repositories.UserRepository{DB: db}
	inventoryRepo := repositories.InventoryRepository{DB: db}

	// Создание сервисов
	userService := services.UserService{UserRepo: &userRepo}
	inventoryService := services.InventoryService{
		InventoryRepo: &inventoryRepo,
		MerchRepo:     &merchRepo,
		UserRepo:      &userService,
	}

	// Создание контроллера
	buyItemController := controllers.NewBuyItemController(&inventoryService)

	// Создание роутера
	r := mux.NewRouter()
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/buy/{item}", buyItemController.BuyItemHandler).Methods("GET")
	t.Run("Successful purchase with Bearer token", func(t *testing.T) {
		// Создание пользователя
		err := userService.CreateUser("testuser", "password")
		if err != nil {
			fmt.Println("dsaw")
			return
		}
		// Генерация JWT-токена
		token, err := GenerateJWT("testuser")
		if err != nil {
			t.Fatalf("Failed to generate JWT token: %v", err)
		}

		// Создание HTTP-запроса
		req := httptest.NewRequest(http.MethodGet, "/api/buy/t-shirt", nil)
		req.Header.Set("Authorization", ""+token)
		recorder := httptest.NewRecorder()
		// Вызов роутера
		r.ServeHTTP(recorder, req)
		// Проверка статуса ответа
		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200, got %d", recorder.Code)

		// Проверка баланса пользователя
		user, _ := userService.GetUser("testuser")
		assert.Equal(t, 920, user.Coins, "Expected coins to be 920, got %d", user.Coins)

		// Проверка наличия товара в инвентаре
		inventory, _ := inventoryService.GetInventory("testuser")
		assert.Equal(t, 1, len(inventory), "Expected inventory length to be 1, got %d", len(inventory))
		assert.Equal(t, "t-shirt", inventory[0].Type, "Expected item type 'cup', got '%s'", inventory[0].Type)
		assert.Equal(t, 1, inventory[0].Quantity, "Expected quantity to be 1, got %d", inventory[0].Quantity)
		_, err = db.Exec("DELETE FROM users WHERE username='testuser'")
		if err != nil {
			return
		}
		_, err = db.Exec("DELETE FROM inventory WHERE username='testuser'")
		if err != nil {
			return
		}
	})
}
