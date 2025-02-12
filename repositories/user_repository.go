package repositories

import (
	"avito-shop/models"
	"database/sql"
	"errors"
	"fmt"
)

//go:generate mockery --name=UserRepositoryInterface
type UserRepositoryInterface interface {
	CreateUser(username, password string) error
	GetUser(username string) (models.User, error)
	UpdateUserCoins(username string, coins int) error
	UserExists(username string) (bool, error)
	AuthenticateUser(username, password string) (bool, error)
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(username, password string) error {
	_, err := r.DB.Exec("INSERT INTO users (username, password, coins) VALUES ($1, $2, 1000)", username, password)
	return err
}

func (r *UserRepository) GetUser(username string) (models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, username, password, coins FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		return models.User{}, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) UpdateUserCoins(username string, coins int) error {
	_, err := r.DB.Exec("UPDATE users SET coins = coins + $1 WHERE username = $2", coins, username)
	return err
}
func (r *UserRepository) UserExists(username string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) AuthenticateUser(username, password string) (bool, error) {
	var storedPassword string
	err := r.DB.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return storedPassword == password, nil
}
