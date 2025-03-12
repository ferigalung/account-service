package models

import "time"

// Balance is the model for the balances table
type Balance struct {
	ID        uint64    `json:"id"`
	AccountID uint64    `json:"account_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName function returns the table name for the Balance model
func (b *Balance) TableName() string {
	return "balances"
}
