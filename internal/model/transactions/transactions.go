package transactions

import "time"

// Transaction is the model for the transactions table
type Transaction struct {
	ID        uint64    `json:"id" db:"id"`
	AccountID uint64    `json:"accountID" db:"account_id"`
	Type      string    `json:"type" db:"type"`
	Amount    float64   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// TableName function returns the table name for the Transaction model
func (t *Transaction) TableName() string {
	return "transactions"
}

type Payload struct {
	AccountNumber string  `json:"accountNumber" validate:"required,numeric,len=11"`
	Type          string  `json:"type" validate:"oneof=deposit withdraw"`
	Amount        float64 `json:"amount" validate:"required,number,min=1"`
}
