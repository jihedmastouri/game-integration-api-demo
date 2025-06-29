package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TransactionType string
type TransactionStatus string

const (
	// Transaction Types
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeCancel   TransactionType = "CANCEL"

	// Transaction Status
	TransactionStatusPending     TransactionStatus = "PENDING"
	TransactionStatusConfirmed   TransactionStatus = "CONFIRMED"
	TransactionStatusFailed      TransactionStatus = "FAILED"
	TransactionStatusCompensated TransactionStatus = "COMPENSATED"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID         uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Currency   Currency
	ProviderID int
	PlayerID   *Player `bun:"rel:belongs-to,join:player_id=id"`
	Status     TransactionStatus
	Type       TransactionType
	Amount     uint
}
