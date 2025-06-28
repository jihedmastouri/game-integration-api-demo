package models

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

