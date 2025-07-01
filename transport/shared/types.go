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
	Amount                         float64         `json:"amount" validate:"required,min=0" example:"1000"`
	ProviderTransactionID          int             `json:"provider_transaction_id" validate:"required" example:"12345"`
	ProviderWithdrawnTransactionID int             `json:"provider_withdrawn_transaction_id" validate:"required" example:"12344"`
}

type DepositResponse struct {
	TransactionID         uuid.UUID                `json:"transaction_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProviderTransactionID int                      `json:"provider_transaction_id" example:"12345"`
	OldBalance            float64                  `json:"old_balance" example:"5000"`
	NewBalance            float64                  `json:"new_balance" example:"6000"`
	Status                models.TransactionStatus `json:"status" example:"CONFIRMED"`
}

type ErrorResponse struct {
	Code errorCode `json:"code" example:"Invalid request"`
	Msg  string    `json:"msg,omitempty" example:"Validation failed"`
}
