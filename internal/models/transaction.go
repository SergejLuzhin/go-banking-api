package models

import "time"

type Transaction struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferRequest struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferByUsernamesRequest struct {
	FromUsername string  `json:"from_username"`
	ToUsername   string  `json:"to_username"`
	Amount       float64 `json:"amount"`
}
