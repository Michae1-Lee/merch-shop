package models

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coin_history"`
}
type CoinHistory struct {
	Received []TransactionHistory `json:"received"`
	Sent     []TransactionHistory `json:"sent"`
}
type TransactionHistory struct {
	FromUser string `json:"from_user,omitempty"`
	ToUser   string `json:"to_user,omitempty"`
	Amount   int    `json:"amount"`
}
type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
