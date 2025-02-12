package main

import (
	"avito-shop/db"
	"avito-shop/middleware"
	"log"
	"net/http"

	"avito-shop/controllers"
	"avito-shop/repositories"
	"avito-shop/services"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	r := mux.NewRouter()

	merchRepo := repositories.MerchRepository{DB: db.GetDB()}
	userRepo := repositories.UserRepository{DB: db.GetDB()}
	transactionRepo := repositories.TransactionRepository{DB: db.GetDB()}
	inventoryRepo := repositories.InventoryRepository{DB: db.GetDB()}

	userService := services.UserService{UserRepo: &userRepo}
	transactionService := services.TransactionService{TransactionRepo: &transactionRepo}
	inventoryService := services.InventoryService{
		InventoryRepo: &inventoryRepo,
		MerchRepo:     &merchRepo,
		UserRepo:      &userService,
	}

	authController := controllers.NewAuthController(&userService)
	infoController := controllers.NewInfoController(&userService, &transactionService, &inventoryService)
	sendCoinController := controllers.NewSendCoinController(&userService, &transactionService)
	buyItemController := controllers.NewBuyItemController(&inventoryService)

	r.HandleFunc("/api/auth", authController.AuthHandler).Methods("POST")

	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/info", infoController.InfoHandler).Methods("GET")
	protected.HandleFunc("/sendCoin", sendCoinController.SendCoinHandler).Methods("POST")
	protected.HandleFunc("/buy/{item}", buyItemController.BuyItemHandler).Methods("GET")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
