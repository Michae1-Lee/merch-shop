package services

import (
	"avito-shop/models"
	"avito-shop/repositories"
	"errors"
)

type UserService struct {
	UserRepo repositories.UserRepositoryInterface
}

func (s *UserService) CreateUser(username, password string) error {
	exists, err := s.UserRepo.UserExists(username)
	if err != nil {
		return errors.New("failed to check user existence")
	}
	if exists {
		return errors.New("user already exists")
	}

	return s.UserRepo.CreateUser(username, password)
}

func (s *UserService) GetUser(username string) (models.User, error) {
	return s.UserRepo.GetUser(username)
}

func (s *UserService) UpdateUserCoins(username string, coins int) error {
	user, err := s.UserRepo.GetUser(username)
	if err != nil {
		return err
	}
	if user.Coins+coins < 0 {
		return errors.New("not enough coins")
	}
	return s.UserRepo.UpdateUserCoins(username, coins)
}

func (s *UserService) AuthenticateUser(username, password string) (bool, error) {
	return s.UserRepo.AuthenticateUser(username, password)
}

func (s *UserService) UserExists(username string) (bool, error) {
	return s.UserRepo.UserExists(username)
}
