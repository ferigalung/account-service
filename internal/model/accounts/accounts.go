package accounts

import "time"

// Account is the model for the accounts table
type Account struct {
	ID            uint64    `json:"id" db:"id"`
	AccountNumber string    `json:"accountNumber" db:"account_number"`
	Name          string    `json:"name" db:"name"`
	NIK           string    `json:"nik" db:"nik"`
	Phone         string    `json:"phone" db:"phone"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

// TableName function returns the table name for the Account model
func (a *Account) TableName() string {
	return "accounts"
}

// Register struct is used to register a new account
type Register struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	NIK   string `json:"nik" validate:"required,numeric,len=16"`
	Phone string `json:"phone" validate:"required,numeric,min=10,max=13"`
}

// UniqueAccount params to check is data unique
type UniqueAccount struct {
	NIK   string `json:"nik" db:"name"`
	Phone string `json:"phone" db:"phone"`
}
