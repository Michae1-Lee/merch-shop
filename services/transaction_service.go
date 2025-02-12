package services

import (
	"avito-shop/models"
	"avito-shop/repositories"
)

//go:generate mockery --name=UserRepositoryInterface

type TransactionService struct {
	TransactionRepo repositories.TransactionRepositoryInterface
}

func (s *TransactionService) GetTransactionsForUser(username string) (models.CoinHistory, error) {
	received, sent, err := s.TransactionRepo.GetTransactionsForUser(username)
	if err != nil {
		return models.CoinHistory{}, err
	}
	return models.CoinHistory{
		Received: received,
		Sent:     sent,
	}, nil
}

func (s *TransactionService) CreateTransaction(fromUser, toUser string, amount int) error {
	return s.TransactionRepo.CreateTransaction(fromUser, toUser, amount)
}
