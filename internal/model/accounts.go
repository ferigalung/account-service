package models

import "time"

// Account is the model for the accounts table
type Account struct {
	ID            uint64    `json:"id"`
	AccountNumber string    `json:"account_number" validate:"required,numeric"`
	Name          string    `json:"name" validate:"required,min=3,max=100"`
	NIK           string    `json:"nik" validate:"required,numeric,len=16"`
	Phone         string    `json:"phone" validate:"required,numeric,min=10,max=13"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
