package transactions

import "time"

// Transaction is the model for the transactions table
type Transaction struct {
	ID        uint64    `json:"-" db:"id"`
	AccountID uint64    `json:"-" db:"account_id"`
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

type ListPayload struct {
	Page          int        `query:"page" json:"page" validate:"required"`
	Size          int        `query:"size" json:"size" validate:"required"`
	AccountNumber string     `json:"accountNumber" validate:"required,numeric,len=11"`
	StartDate     *time.Time `query:"startDate" json:"startDate" validate:"omitempty"`
	EndDate       *time.Time `query:"endDate" json:"endDate" validate:"omitempty"`
	AccountID     uint64     `json:"-"`
}

type PaginationMeta struct {
	Page            int   `json:"page"`
	TotalPage       int   `json:"totalPage"`
	TotalData       int64 `json:"totalData"`
	TotalDataOnPage int   `json:"totalDataOnPage"`
}

type ListRepo struct {
	Data         []Transaction `json:"data"`
	TotalRecords int64         `json:"totalRecords"`
}

type Pagination struct {
	Data []Transaction  `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
