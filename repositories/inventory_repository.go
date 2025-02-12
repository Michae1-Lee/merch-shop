package repositories

import (
	"avito-shop/models"
	"database/sql"
)

//go:generate mockery --name=InventoryRepositoryInterface
type InventoryRepositoryInterface interface {
	GetInventory(username string) ([]models.InventoryItem, error)
	AddItem(username, itemType string, quantity int) error
}

type InventoryRepository struct {
	DB *sql.DB
}

func (r *InventoryRepository) GetInventory(username string) ([]models.InventoryItem, error) {
	rows, err := r.DB.Query("SELECT type, quantity FROM inventory WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, err
		}
		inventory = append(inventory, item)
	}

	return inventory, nil
}

func (r *InventoryRepository) AddItem(username, itemType string, quantity int) error {
	_, err := r.DB.Exec(`
        INSERT INTO inventory (username, type, quantity)
        VALUES ($1, $2, $3)
        ON CONFLICT (username, type) DO UPDATE SET quantity = inventory.quantity + $3
    `, username, itemType, quantity)
	return err
}
