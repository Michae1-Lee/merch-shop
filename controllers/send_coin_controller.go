package controllers

import (
	"avito-shop/models"
	"avito-shop/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type SendCoinController struct {
	UserService        *services.UserService
	TransactionService *services.TransactionService
}

func NewSendCoinController(
	userService *services.UserService,
	transactionService *services.TransactionService,
) *SendCoinController {
	return &SendCoinController{
		UserService:        userService,
		TransactionService: transactionService,
	}
}

func (c *SendCoinController) SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println(req.ToUser)

	if err := c.UserService.UpdateUserCoins(username, -req.Amount); err != nil {
		fmt.Println("dsa")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.UserService.UpdateUserCoins(req.ToUser, req.Amount); err != nil {
		err := c.UserService.UpdateUserCoins(username, req.Amount)
		if err != nil {
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.TransactionService.CreateTransaction(username, req.ToUser, req.Amount); err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
