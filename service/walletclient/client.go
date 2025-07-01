package walletclient

import (
	"net/http"
	"time"
)

type WalletClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewWalletClient(baseURL, token string) *WalletClient {
	return &WalletClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type BalanceResponse struct {
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

type DepositRequest struct {
	UserID       int                         `json:"userId" binding:"required"`
	Currency     string                      `json:"currency" binding:"required"`
	Transactions []DepositRequestTransaction `json:"transactions" binding:"required"`
}

type DepositRequestTransaction struct {
	Amount    float64 `json:"amount" binding:"required"`
	BetID     int     `json:"betId" binding:"required"`
	Reference string  `json:"reference" binding:"required"`
}

type WithdrawRequest struct {
	UserID       int                          `json:"userId" binding:"required"`
	Currency     string                       `json:"currency" binding:"required"`
	Transactions []WithdrawRequestTransaction `json:"transactions" binding:"required"`
}

type WithdrawRequestTransaction struct {
	Amount    float64 `json:"amount" binding:"required"`
	BetID     int     `json:"betId" binding:"required"`
	Reference string  `json:"reference" binding:"required"`
}

type OperationResponse struct {
	Balance      string                         `json:"balance"`
	Transactions []OperationResponseTransaction `json:"transactions"`
}

type OperationResponseTransaction struct {
	ID        int    `json:"id"`
	Reference string `json:"reference"`
}

type ErrorResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
