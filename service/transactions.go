package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/service/walletclient"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
)

func (s *Service) ProcessBet(ctx context.Context, player *models.Player, req shared.WithdrawRequest) (*shared.BetOperationResponse, error) {

	// Create transaction record
	transaction := &models.Transaction{
		PlayerID:   player.ID,
		ProviderID: req.ProviderTransactionID,
		Amount:     uint64(req.Amount),
		Currency:   req.Currency,
		Status:     models.TransactionStatusPending,
		Type:       models.TransactionTypeWithdraw,
		Attempts:   0,
	}

	err := s.Repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Get current balance before withdrawal
	balanceResp, err := s.WalletClient.GetBalance(player.ID)
	if err != nil {
		transaction.Status = models.TransactionStatusFailed
		s.Repository.UpdateTransaction(ctx, transaction)
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	oldBalance := balanceResp.Balance

	// Process withdrawal through wallet client
	withdrawReq := walletclient.WithdrawRequest{
		UserID:   int(player.ID),
		Currency: string(req.Currency),
		Transactions: []walletclient.WithdrawRequestTransaction{
			{
				Amount:    req.Amount,
				BetID:     req.ProviderTransactionID,
				Reference: strconv.Itoa(req.ProviderTransactionID),
			},
		},
	}

	withdrawResp, err := s.WalletClient.Withdraw(withdrawReq)
	if err != nil {
		transaction.Status = models.TransactionStatusFailed
		s.Repository.UpdateTransaction(ctx, transaction)
		return nil, fmt.Errorf("failed to process withdrawal: %w", err)
	}

	// Update transaction status
	transaction.Status = models.TransactionStatusConfirmed
	err = s.Repository.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &shared.BetOperationResponse{
		TransactionID:         transaction.ID,
		ProviderTransactionID: req.ProviderTransactionID,
		OldBalance:            oldBalance,
		NewBalance:            withdrawResp.Balance,
		Status:                transaction.Status,
	}, nil
}

func (s *Service) ProcessSettle(ctx context.Context, player *models.Player, req shared.DepositRequest) (*shared.BetOperationResponse, error) {
	// Create transaction record
	transaction := &models.Transaction{
		PlayerID:   player.ID,
		ProviderID: req.ProviderTransactionID,
		Amount:     uint64(req.Amount),
		Currency:   req.Currency,
		Status:     models.TransactionStatusPending,
		Type:       models.TransactionTypeDeposit,
		Attempts:   0,
	}

	err := s.Repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Get current balance before deposit
	balanceResp, err := s.WalletClient.GetBalance(player.ID)
	if err != nil {
		transaction.Status = models.TransactionStatusFailed
		s.Repository.UpdateTransaction(ctx, transaction)
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	oldBalance := balanceResp.Balance

	// Process deposit through wallet client (only if amount > 0)
	var newBalance string
	if req.Amount > 0 {
		depositReq := walletclient.DepositRequest{
			UserID:   int(player.ID),
			Currency: string(req.Currency),
			Transactions: []walletclient.DepositRequestTransaction{
				{
					Amount:    req.Amount,
					BetID:     req.ProviderWithdrawnTransactionID,
					Reference: strconv.Itoa(req.ProviderTransactionID),
				},
			},
		}

		depositResp, err := s.WalletClient.Deposit(depositReq)
		if err != nil {
			transaction.Status = models.TransactionStatusFailed
			s.Repository.UpdateTransaction(ctx, transaction)
			return nil, fmt.Errorf("failed to process deposit: %w", err)
		}
		newBalance = depositResp.Balance
	} else {
		// If amount is 0, bet is lost - no deposit needed
		newBalance = oldBalance
	}

	// Update transaction status
	transaction.Status = models.TransactionStatusConfirmed
	err = s.Repository.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &shared.BetOperationResponse{
		TransactionID:         transaction.ID,
		ProviderTransactionID: req.ProviderTransactionID,
		OldBalance:            oldBalance,
		NewBalance:            newBalance,
		Status:                transaction.Status,
	}, nil
}

func (s *Service) ProcessCancel(ctx context.Context, player *models.Player, req shared.CancelRequest) (*shared.BetOperationResponse, error) {
	// Find the original transaction to cancel
	originalTx, err := s.Repository.GetTransactionByProviderID(ctx, req.ProviderTransactionID)
	if err != nil {
		return nil, fmt.Errorf("original transaction not found: %w", err)
	}

	// Validate the transaction belongs to this player
	if originalTx.PlayerID != player.ID {
		return nil, errors.New("transaction does not belong to this player")
	}

	// Create cancel transaction record
	cancelTx := &models.Transaction{
		PlayerID:   player.ID,
		ProviderID: req.ProviderTransactionID,
		Amount:     originalTx.Amount,
		Currency:   originalTx.Currency,
		Status:     models.TransactionStatusPending,
		Type:       models.TransactionTypeCancel,
		Attempts:   0,
	}

	err = s.Repository.CreateTransaction(ctx, cancelTx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cancel transaction: %w", err)
	}

	// Get current balance
	balanceResp, err := s.WalletClient.GetBalance(player.ID)
	if err != nil {
		cancelTx.Status = models.TransactionStatusFailed
		s.Repository.UpdateTransaction(ctx, cancelTx)
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	oldBalance := balanceResp.Balance

	var newBalance string

	// Reverse the original transaction
	if originalTx.Type == models.TransactionTypeWithdraw {
		// Original was a withdrawal (bet), so we need to deposit back
		depositReq := walletclient.DepositRequest{
			UserID:   int(player.ID),
			Currency: string(originalTx.Currency),
			Transactions: []walletclient.DepositRequestTransaction{
				{
					Amount:    float64(originalTx.Amount),
					BetID:     originalTx.ProviderID,
					Reference: fmt.Sprintf("cancel-%d", req.ProviderTransactionID),
				},
			},
		}

		depositResp, err := s.WalletClient.Deposit(depositReq)
		if err != nil {
			cancelTx.Status = models.TransactionStatusFailed
			s.Repository.UpdateTransaction(ctx, cancelTx)
			return nil, fmt.Errorf("failed to process cancel deposit: %w", err)
		}
		newBalance = depositResp.Balance
	} else if originalTx.Type == models.TransactionTypeDeposit {
		// Original was a deposit (settle), so we need to withdraw back
		withdrawReq := walletclient.WithdrawRequest{
			UserID:   int(player.ID),
			Currency: string(originalTx.Currency),
			Transactions: []walletclient.WithdrawRequestTransaction{
				{
					Amount:    float64(originalTx.Amount),
					BetID:     originalTx.ProviderID,
					Reference: fmt.Sprintf("cancel-%d", req.ProviderTransactionID),
				},
			},
		}

		withdrawResp, err := s.WalletClient.Withdraw(withdrawReq)
		if err != nil {
			cancelTx.Status = models.TransactionStatusFailed
			s.Repository.UpdateTransaction(ctx, cancelTx)
			return nil, fmt.Errorf("failed to process cancel withdrawal: %w", err)
		}
		newBalance = withdrawResp.Balance
	} else {
		return nil, errors.New("cannot cancel a cancel transaction")
	}

	// Update transaction statuses
	cancelTx.Status = models.TransactionStatusConfirmed
	originalTx.Status = models.TransactionStatusCompensated

	err = s.Repository.UpdateTransaction(ctx, cancelTx)
	if err != nil {
		return nil, fmt.Errorf("failed to update cancel transaction: %w", err)
	}

	err = s.Repository.UpdateTransaction(ctx, originalTx)
	if err != nil {
		return nil, fmt.Errorf("failed to update original transaction: %w", err)
	}

	return &shared.BetOperationResponse{
		TransactionID:         cancelTx.ID,
		ProviderTransactionID: req.ProviderTransactionID,
		OldBalance:            oldBalance,
		NewBalance:            newBalance,
		Status:                cancelTx.Status,
	}, nil
}
