package balances

import "time"

// Balance is the model for the balances table
type Balance struct {
	ID        uint64    `json:"id" db:"id"`
	AccountID uint64    `json:"accountID" db:"account_id"`
	Balance   float64   `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// TableName function returns the table name for the Balance model
func (b *Balance) TableName() string {
	return "balances"
}

type DetailBalance struct {
	AccountNumber string  `json:"accountNumber" db:"account_number"`
	Name          string  `json:"name" db:"name"`
	Balance       float64 `json:"balance" db:"balance"`
}
