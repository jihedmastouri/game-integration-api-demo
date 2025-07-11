package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Currency string

type TransactionType string
type TransactionStatus string

const (
	// Transaction Types
	TransactionTypeWithdraw TransactionType = "WITHDRAW"
	TransactionTypeDeposit  TransactionType = "DEPOSIT"
	TransactionTypeCancel   TransactionType = "CANCEL"

	// Transaction Status
	TransactionStatusPending    TransactionStatus = "PENDING"
	TransactionStatusConfirmed  TransactionStatus = "CONFIRMED"
	TransactionStatusFailed     TransactionStatus = "FAILED"
	TransactionStatusFinalized  TransactionStatus = "FINAL"
	TransactionStatusProcessing TransactionStatus = "PROCESSING"

	// Currencies
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyKES Currency = "KES"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID                 uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Player             *Player   `bun:"rel:belongs-to,join:player_id=id"`
	PlayerID           uint64    `json:"-"`
	ProviderID         uint64    `bun:"provider_id,nullzero"`
	WithdrawProviderID uint64    `bun:"withdraw_provider_id,nullzero"`
	Amount             string
	Currency           Currency
	Status             TransactionStatus
	Type               TransactionType
	Attempts           int
	CreatedAt          time.Time `bun:"created_at"`
	UpdatedAt          time.Time `bun:"updated_at"`
}
