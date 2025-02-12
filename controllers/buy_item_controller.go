package controllers

import (
	"avito-shop/services"
	"github.com/gorilla/mux"
	"net/http"
)

type BuyItemController struct {
	InventoryService *services.InventoryService
}

func NewBuyItemController(inventoryService *services.InventoryService) *BuyItemController {
	return &BuyItemController{
		InventoryService: inventoryService,
	}
}

/*
/
*/
func (c *BuyItemController) BuyItemHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	vars := mux.Vars(r)
	itemType := vars["item"]

	err := c.InventoryService.BuyItem(username, itemType)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(itemType))
	if err != nil {
		return
	}

}
