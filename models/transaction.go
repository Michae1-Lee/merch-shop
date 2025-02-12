package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	FromUser  string    `json:"from_user"`
	ToUser    string    `json:"to_user"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
