package controllers

import (
	"avito-shop/models"
	"avito-shop/services"
	"encoding/json"
	"net/http"
)

type InfoController struct {
	UserService        *services.UserService
	TransactionService *services.TransactionService
	InventoryService   *services.InventoryService
}

func NewInfoController(
	userService *services.UserService,
	transactionService *services.TransactionService,
	inventoryService *services.InventoryService,
) *InfoController {
	return &InfoController{
		UserService:        userService,
		TransactionService: transactionService,
		InventoryService:   inventoryService,
	}
}

func (c *InfoController) InfoHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	user, err := c.UserService.GetUser(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	coinHistory, err := c.TransactionService.GetTransactionsForUser(username)
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	inventory, err := c.InventoryService.GetInventory(username)
	if err != nil {
		http.Error(w, "Failed to fetch inventory", http.StatusInternalServerError)
		return
	}

	response := models.InfoResponse{
		Coins:       user.Coins,
		Inventory:   inventory,
		CoinHistory: coinHistory,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
