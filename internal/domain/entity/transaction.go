package entity

import "time"

const (
	TransactionTypeDebit = "DEBIT"
)

type Transaction struct {
	ID        int64
	UserID    int64
	Type      string
	Amount    int64
	CreatedAt time.Time
}
