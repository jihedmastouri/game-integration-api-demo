package shared

import (
	"github.com/google/uuid"
	"github.com/jihedmastouri/game-integration-api-demo/models"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type PlayerInfoResponse struct {
	PlayerID uint64 `json:"user_id" example:"1"`
	Balance  string `json:"balance" example:"1000.50"`
	Currency string `json:"currency" example:"USD"`
}

type DepositRequest struct {
	Currency                       models.Currency `json:"currency" validate:"required" example:"USD"`
	Amount                         float64         `json:"amount" validate:"min=0" example:"1000.00"`
	ProviderTransactionID          uint64          `json:"provider_transaction_id" validate:"required" example:"12345"`
	ProviderWithdrawnTransactionID uint64          `json:"provider_withdrawn_transaction_id" validate:"required" example:"12344"`
}

type WithdrawRequest struct {
	Currency              models.Currency `json:"currency" validate:"required" example:"USD"`
	Amount                float64         `json:"amount" validate:"required,gt=0" example:"100"`
	ProviderTransactionID uint64          `json:"provider_transaction_id" validate:"required" example:"12345"`
}

type CancelRequest struct {
	ProviderTransactionID uint64 `json:"provider_transaction_id" validate:"required" example:"12345"`
}

type BetOperationResponse struct {
	TransactionID         uuid.UUID                `json:"transaction_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProviderTransactionID uint64                   `json:"provider_transaction_id" example:"12345"`
	OldBalance            string                   `json:"old_balance" example:"1900.50"`
	NewBalance            string                   `json:"new_balance" example:"1000.50"`
	Status                models.TransactionStatus `json:"status" example:"CONFIRMED"`
}

type ErrorResponse struct {
	Code errorCode `json:"code" example:"Invalid request"`
	Msg  string    `json:"msg,omitempty" example:"Validation failed"`
}
