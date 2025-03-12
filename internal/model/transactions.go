package models

import "time"

// Transaction is the model for the transactions table
type Transaction struct {
	ID        uint64    `json:"id"`
	AccountID uint64    `json:"account_id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName function returns the table name for the Transaction model
func (t *Transaction) TableName() string {
	return "transactions"
}
