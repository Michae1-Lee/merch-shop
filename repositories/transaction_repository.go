package repositories

import (
	"avito-shop/models"
	"database/sql"
)

//go:generate mockery --name=TransactionRepositoryInterface
type TransactionRepositoryInterface interface {
	GetTransactionsForUser(username string) ([]models.TransactionHistory, []models.TransactionHistory, error)
	CreateTransaction(fromUser, toUser string, amount int) error
}

type TransactionRepository struct {
	DB *sql.DB
}

func (r *TransactionRepository) GetTransactionsForUser(username string) ([]models.TransactionHistory, []models.TransactionHistory, error) {
	rows, err := r.DB.Query(`
        SELECT from_user, to_user, amount FROM transactions WHERE to_user = $1
    `, username)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var received []models.TransactionHistory
	for rows.Next() {
		var t models.TransactionHistory
		if err := rows.Scan(&t.FromUser, &t.ToUser, &t.Amount); err != nil {
			return nil, nil, err
		}
		received = append(received, t)
	}

	rows, err = r.DB.Query(`
        SELECT from_user, to_user, amount FROM transactions WHERE from_user = $1
    `, username)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var sent []models.TransactionHistory
	for rows.Next() {
		var t models.TransactionHistory
		if err := rows.Scan(&t.FromUser, &t.ToUser, &t.Amount); err != nil {
			return nil, nil, err
		}
		sent = append(sent, t)
	}

	return received, sent, nil
}

func (r *TransactionRepository) CreateTransaction(fromUser, toUser string, amount int) error {
	_, err := r.DB.Exec(`
        INSERT INTO transactions (from_user, to_user, amount, created_at)
        VALUES ($1, $2, $3, NOW())
    `, fromUser, toUser, amount)
	return err
}
