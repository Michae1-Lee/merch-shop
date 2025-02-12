package repositories

import (
	"database/sql"
	"errors"
)

type MerchRepository struct {
	DB *sql.DB
}

func (r *MerchRepository) GetMerchPrice(itemName string) (int, error) {
	var price int

	err := r.DB.QueryRow("SELECT price FROM merch WHERE name = $1", itemName).Scan(&price)

	if errors.Is(err, sql.ErrNoRows) {

		return 0, nil
	}
	return price, err
}
