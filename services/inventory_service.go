package services

import (
	"avito-shop/models"
	"avito-shop/repositories"
	"errors"
)

type InventoryService struct {
	InventoryRepo repositories.InventoryRepositoryInterface
	MerchRepo     *repositories.MerchRepository
	UserRepo      *UserService
}

func (s *InventoryService) AddItem(username, itemType string, quantity int) error {
	if quantity <= 0 {
		return errors.New("invalid quantity")
	}
	return s.InventoryRepo.AddItem(username, itemType, quantity)
}

func (s *InventoryService) GetInventory(username string) ([]models.InventoryItem, error) {
	inventory, err := s.InventoryRepo.GetInventory(username)
	if err != nil {
		return nil, errors.New("failed to fetch inventory")
	}
	return inventory, nil
}
func (s *InventoryService) BuyItem(username, itemType string) error {
	price, err := s.MerchRepo.GetMerchPrice(itemType)

	if err != nil {
		return errors.New("failed to fetch item price")
	}
	if price == 0 {
		return errors.New("invalid item")
	}

	if err := s.UserRepo.UpdateUserCoins(username, -price); err != nil {
		return err
	}

	if err := s.InventoryRepo.AddItem(username, itemType, 1); err != nil {
		err := s.UserRepo.UpdateUserCoins(username, price)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}
