package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/service/walletclient"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
)

// Helper method to check if player has pending or processing transactions
func (s *Service) hasPendingTransactions(ctx context.Context, playerID uint64) (bool, error) {
	pendingTxs, err := s.GetFirstPendingTransactionsByPlayerID(ctx, playerID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	// Also check for processing transactions
	processingTxs, err := s.GetFirstProcessingTransactionsByPlayerID(ctx, playerID)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return processingTxs != nil || pendingTxs != nil, nil
}

func (s *Service) ProcessBet(ctx context.Context, player *models.Player, req shared.WithdrawRequest) (*shared.BetOperationResponse, error) {
	prevTx, err := s.GetTransactionByProviderID(ctx, req.ProviderTransactionID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check for any previous transactions: %w", err)
	}

	// reject duplicate transaction
	if prevTx != nil {
		return nil, fmt.Errorf("Duplicate transaction")
	}

	// Check if player has any pending transactions
	hasPending, err := s.hasPendingTransactions(ctx, player.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check pending transactions: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		PlayerID:   player.ID,
		ProviderID: req.ProviderTransactionID,
		Amount:     strconv.FormatFloat(req.Amount, 'f', -1, 64),
		Currency:   req.Currency,
		Status:     models.TransactionStatusPending,
		Type:       models.TransactionTypeWithdraw,
		Attempts:   0,
	}

	err = s.Repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// If player has pending transactions, keep this one pending too
	if hasPending {
		slog.Info("Player has pending transactions, keeping bet transaction pending", "player_id", player.ID, "transaction_id", transaction.ID)
		return &shared.BetOperationResponse{
			TransactionID:         transaction.ID,
			ProviderTransactionID: req.ProviderTransactionID,
			Status:                transaction.Status, // PENDING
		}, nil
	}

	// Try to process with wallet service
	balanceResp, err := s.WalletClient.GetBalance(player.ID)
	if err != nil {
		slog.Error("Failed to get balance, keeping transaction pending", "error", err, "player_id", player.ID, "transaction_id", transaction.ID)
		return &shared.BetOperationResponse{
			TransactionID:         transaction.ID,
			ProviderTransactionID: req.ProviderTransactionID,
			Status:                transaction.Status, // PENDING
		}, nil
	}
	oldBalance := balanceResp.Balance

	// Process withdrawal through wallet client
	withdrawReq := walletclient.WithdrawRequest{
		UserID:   int(player.ID),
		Currency: string(req.Currency),
		Transactions: []walletclient.WithdrawRequestTransaction{
			{
				Amount:    float64(req.Amount),
				BetID:     req.ProviderTransactionID,
				Reference: transaction.ID.String(),
			},
		},
	}

	withdrawResp, err := s.WalletClient.Withdraw(withdrawReq)
	if err != nil {
		slog.Error("Failed to process withdrawal, keeping transaction pending", "error", err, "player_id", player.ID, "transaction_id", transaction.ID)
		return &shared.BetOperationResponse{
			TransactionID:         transaction.ID,
			ProviderTransactionID: req.ProviderTransactionID,
			OldBalance:            oldBalance,
			NewBalance:            oldBalance,         // No change since withdrawal failed
			Status:                transaction.Status, // PENDING
		}, nil
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
	prevTx, err := s.GetTransactionByProviderID(ctx, req.ProviderTransactionID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check for any previous transactions: %w", err)
	}

	// reject duplicate transaction
	if prevTx != nil {
		return nil, fmt.Errorf("Duplicate transaction")
	}

	oldTx, err := s.GetTransactionByProviderID(ctx, req.ProviderWithdrawnTransactionID)
	if err != nil || oldTx == nil {
		return nil, fmt.Errorf("Failed to get previous bet. make sure to include a valide previous bet_id to settle it.")
	}

	if oldTx.Status == models.TransactionStatusFinalized {
		return nil, fmt.Errorf("This bet is already settled.")
	}

	if oldTx.Status == models.TransactionStatusFailed {
		return nil, fmt.Errorf("The bet failed. you can not settle failed bets")
	}

	// Check if player has any pending transactions
	hasPending, err := s.hasPendingTransactions(ctx, player.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check pending transactions: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		PlayerID:   player.ID,
		ProviderID: req.ProviderTransactionID,
		Amount:     strconv.FormatFloat(req.Amount, 'f', -1, 64),
		Currency:   req.Currency,
		Status:     models.TransactionStatusPending,
		Type:       models.TransactionTypeDeposit,
		Attempts:   0,
	}

	err = s.Repository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// If player has pending transactions, keep this one pending too
	if hasPending {
		slog.Info("Player has pending transactions, keeping settle transaction pending", "player_id", player.ID, "transaction_id", transaction.ID)
		return &shared.BetOperationResponse{
			TransactionID:         transaction.ID,
			ProviderTransactionID: req.ProviderTransactionID,
			Status:                transaction.Status, // PENDING
		}, nil
	}

	// Try to process with wallet service
	balanceResp, err := s.WalletClient.GetBalance(player.ID)
	if err != nil {
		slog.Error("Failed to get balance, keeping transaction pending", "error", err, "player_id", player.ID, "transaction_id", transaction.ID)
		return &shared.BetOperationResponse{
			TransactionID:         transaction.ID,
			ProviderTransactionID: req.ProviderTransactionID,
			Status:                transaction.Status, // PENDING
		}, nil
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
					Amount:    float64(req.Amount),
					BetID:     req.ProviderWithdrawnTransactionID,
					Reference: transaction.ID.String(),
				},
			},
		}

		depositResp, err := s.WalletClient.Deposit(depositReq)
		if err != nil {
			slog.Error("Failed to process deposit, keeping transaction pending", "error", err, "player_id", player.ID, "transaction_id", transaction.ID)
			return &shared.BetOperationResponse{
				TransactionID:         transaction.ID,
				ProviderTransactionID: req.ProviderTransactionID,
				OldBalance:            oldBalance,
				NewBalance:            oldBalance,         // No change since deposit failed
				Status:                transaction.Status, // PENDING
			}, nil
		}

		oldTx.Status = models.TransactionStatusFinalized
		err = s.UpdateTransaction(ctx, oldTx)
		if err != nil {
			slog.Error("failed to update withdraw transaction", "error", err)
		}

		newBalance = depositResp.Balance
	} else {
		// If amount is 0, bet is lost - no deposit needed
		newBalance = oldBalance

		oldTx.Status = models.TransactionStatusFinalized
		err = s.UpdateTransaction(ctx, oldTx)
		if err != nil {
			slog.Error("failed to update withdraw transaction", "error", err)
		}
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
	originalTx, err := s.GetTransactionByProviderID(ctx, req.ProviderTransactionID)
	if err != nil || originalTx == nil {
		slog.Error("original transaction not found", "error", err)
		return nil, fmt.Errorf("original transaction not found")
	}

	// Validate the transaction belongs to this player
	if originalTx.PlayerID != player.ID {
		return nil, errors.New("transaction does not belong to this player")
	}

	// Validate the transaction belongs to this player
	if originalTx.Status == models.TransactionStatusFinalized {
		return nil, errors.New("transaction already finalized")
	}

	// Create cancel transaction record
	cancelTx := &models.Transaction{
		PlayerID:           player.ID,
		WithdrawProviderID: req.ProviderTransactionID,
		Amount:             originalTx.Amount,
		Currency:           originalTx.Currency,
		Status:             models.TransactionStatusPending,
		Type:               models.TransactionTypeCancel,
		Attempts:           0,
	}

	err = s.Repository.CreateTransaction(ctx, cancelTx)
	if err != nil {
		return nil, fmt.Errorf("failed to create cancel transaction: %w", err)
	}

	// Check if player has any pending transactions
	hasPending, err := s.hasPendingTransactions(ctx, player.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check pending transactions: %w", err)
	}

	// If player has pending transactions, keep this one pending too
	if hasPending {
		slog.Info("Player has pending transactions, keeping cancel transaction pending", "player_id", player.ID, "transaction_id", cancelTx.ID)
		return &shared.BetOperationResponse{
			TransactionID:         cancelTx.ID,
			ProviderTransactionID: req.ProviderTransactionID,
			Status:                cancelTx.Status, // PENDING
		}, nil
	}

	// Try to process with wallet service
	balanceResp, err := s.WalletClient.GetBalance(player.ID)
	if err != nil {
		slog.Error("Failed to get balance, keeping transaction pending", "error", err, "player_id", player.ID, "transaction_id", cancelTx.ID)
		return &shared.BetOperationResponse{
			TransactionID:         cancelTx.ID,
			ProviderTransactionID: req.ProviderTransactionID,
			Status:                cancelTx.Status, // PENDING
		}, nil
	}
	oldBalance := balanceResp.Balance

	ogAmount, err := strconv.ParseFloat(originalTx.Amount, 64)
	if err != nil {
		return nil, err
	}

	var newBalance string

	originalTx.Status = models.TransactionStatusFinalized
	err = s.Repository.UpdateTransaction(ctx, originalTx)
	if err != nil {
		return nil, fmt.Errorf("failed to update original transaction: %w", err)
	}

	// Reverse the original transaction
	if originalTx.Type == models.TransactionTypeWithdraw {
		// Original was a withdrawal (bet), so we need to deposit back
		depositReq := walletclient.DepositRequest{
			UserID:   int(player.ID),
			Currency: string(originalTx.Currency),
			Transactions: []walletclient.DepositRequestTransaction{
				{
					Amount:    ogAmount,
					BetID:     originalTx.ProviderID,
					Reference: fmt.Sprintf("cancel-%d", req.ProviderTransactionID),
				},
			},
		}

		depositResp, err := s.WalletClient.Deposit(depositReq)
		if err != nil {
			slog.Error("Failed to process cancel deposit, keeping transaction pending", "error", err, "player_id", player.ID, "transaction_id", cancelTx.ID)
			return &shared.BetOperationResponse{
				TransactionID:         cancelTx.ID,
				ProviderTransactionID: req.ProviderTransactionID,
				OldBalance:            oldBalance,
				NewBalance:            oldBalance,      // No change since deposit failed
				Status:                cancelTx.Status, // PENDING
			}, nil
		}
		newBalance = depositResp.Balance

	} else if originalTx.Type == models.TransactionTypeDeposit {
		// Original was a deposit (settle), so we need to withdraw back
		withdrawReq := walletclient.WithdrawRequest{
			UserID:   int(player.ID),
			Currency: string(originalTx.Currency),
			Transactions: []walletclient.WithdrawRequestTransaction{
				{
					Amount:    ogAmount,
					BetID:     originalTx.ProviderID,
					Reference: fmt.Sprintf("cancel-%d", req.ProviderTransactionID),
				},
			},
		}

		withdrawResp, err := s.WalletClient.Withdraw(withdrawReq)
		if err != nil {
			slog.Error("Failed to process cancel withdrawal, keeping transaction pending", "error", err, "player_id", player.ID, "transaction_id", cancelTx.ID)
			return &shared.BetOperationResponse{
				TransactionID:         cancelTx.ID,
				ProviderTransactionID: req.ProviderTransactionID,
				OldBalance:            oldBalance,
				NewBalance:            oldBalance,      // No change since withdrawal failed
				Status:                cancelTx.Status, // PENDING
			}, nil
		}
		newBalance = withdrawResp.Balance
	} else {
		return nil, errors.New("cannot cancel a cancel transaction")
	}

	// Update transaction statuses
	cancelTx.Status = models.TransactionStatusConfirmed
	err = s.Repository.UpdateTransaction(ctx, cancelTx)
	if err != nil {
		return nil, fmt.Errorf("failed to update cancel transaction: %w", err)
	}

	return &shared.BetOperationResponse{
		TransactionID:         cancelTx.ID,
		ProviderTransactionID: req.ProviderTransactionID,
		OldBalance:            oldBalance,
		NewBalance:            newBalance,
		Status:                cancelTx.Status,
	}, nil
}
