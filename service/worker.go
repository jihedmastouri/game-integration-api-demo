package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/service/walletclient"
)

const (
	MaxRetryAttempts = 3
	RetryInterval    = 30 * time.Second
)

// StartPendingTransactionWorker starts a background worker to process pending transactions
func (s *Service) StartPendingTransactionWorker(ctx context.Context) {
	ticker := time.NewTicker(RetryInterval)
	defer ticker.Stop()

	slog.Info("Starting pending transaction worker")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping pending transaction worker")
			return
		case <-ticker.C:
			s.processPendingTransactions(ctx)
		}
	}
}

// processPendingTransactions processes pending transactions one at a time per user
func (s *Service) processPendingTransactions(ctx context.Context) {
	// Process transactions one at a time until no more processable transactions
	for {
		// Get the next processable transaction
		tx, err := s.Repository.GetNextProcessableTransaction(ctx)
		if err != nil {
			// If no rows found, it means no processable transactions available
			if err == sql.ErrNoRows {
				return // No more transactions to process
			}
			slog.Error("Failed to get next processable transaction", "error", err)
			return
		}

		if tx == nil {
			return // No more transactions to process
		}

		// Check if transaction exceeded max retry attempts
		if tx.Attempts >= MaxRetryAttempts {
			slog.Warn("Transaction exceeded max retry attempts, marking as failed", "transaction_id", tx.ID, "attempts", tx.Attempts)
			tx.Status = models.TransactionStatusFailed
			if err := s.Repository.UpdateTransaction(ctx, tx); err != nil {
				slog.Error("Failed to update failed transaction", "error", err, "transaction_id", tx.ID)
			}
			continue // Try next transaction
		}

		// Atomically start processing this transaction
		if err := s.Repository.StartProcessingTransaction(ctx, tx.ID); err != nil {
			slog.Warn("Failed to start processing transaction (may already be processing)", "error", err, "transaction_id", tx.ID)
			continue // Try next transaction
		}

		slog.Info("Processing transaction", "transaction_id", tx.ID, "type", tx.Type, "player_id", tx.PlayerID)

		// Increment attempt counter
		tx.Attempts++

		// Process the transaction
		retrySTatus := models.TransactionStatusPending
		switch tx.Type {
		case models.TransactionTypeWithdraw:
			retrySTatus = s.retryWithdraw(ctx, tx)
		case models.TransactionTypeDeposit:
			retrySTatus = s.retryDeposit(ctx, tx)
		case models.TransactionTypeCancel:
			retrySTatus = s.retryCancel(ctx, tx)
		default:
			slog.Error("Unknown transaction type", "transaction_id", tx.ID, "type", tx.Type)
		}

		tx.Status = retrySTatus
		if err := s.Repository.UpdateTransaction(ctx, tx); err != nil {
			slog.Error("Failed to update transaction after retry", "error", err, "transaction_id", tx.ID)
		}

		// Small delay between transactions to prevent overwhelming the wallet service
		time.Sleep(1 * time.Second)
	}
}

// retryWithdraw retries a pending withdrawal transaction
func (s *Service) retryWithdraw(ctx context.Context, tx *models.Transaction) models.TransactionStatus {
	ogAmount, err := strconv.ParseFloat(tx.Amount, 64)
	if err != nil {
		slog.Error(err.Error())
		return models.TransactionStatusFailed
	}

	withdrawReq := walletclient.WithdrawRequest{
		UserID:   int(tx.PlayerID),
		Currency: string(tx.Currency),
		Transactions: []walletclient.WithdrawRequestTransaction{
			{
				Amount:    ogAmount,
				BetID:     tx.ProviderID,
				Reference: tx.ID.String(),
			},
		},
	}

	_, err = s.WalletClient.Withdraw(withdrawReq)
	if err != nil {
		slog.Error("Failed to retry withdrawal", "error", err, "transaction_id", tx.ID)
		return models.TransactionStatusPending
	}

	return models.TransactionStatusConfirmed
}

// retryDeposit retries a pending deposit transaction
func (s *Service) retryDeposit(ctx context.Context, tx *models.Transaction) models.TransactionStatus {
	ogAmount, err := strconv.ParseFloat(tx.Amount, 64)
	if err != nil {
		slog.Error(err.Error())
		return models.TransactionStatusFailed
	}

	oldTx, err := s.GetTransactionByProviderID(ctx, tx.WithdrawProviderID)
	if err != nil || oldTx == nil {
		slog.Error("failed to get withdraw transaction", "error", err)
		return models.TransactionStatusFailed
	}

	if oldTx.Status == models.TransactionStatusFailed {
		slog.Error("Withdraw transaction is failed. you cannot deposite")
		return models.TransactionStatusFailed
	}

	if ogAmount < 0 {
		slog.Error("og amount should not be negative", "transaction_id", tx.ID.String())
		return models.TransactionStatusFailed
	}

	if ogAmount > 0 {
		depositReq := walletclient.DepositRequest{
			UserID:   int(tx.PlayerID),
			Currency: string(tx.Currency),
			Transactions: []walletclient.DepositRequestTransaction{
				{
					Amount:    ogAmount,
					BetID:     tx.ProviderID,
					Reference: tx.ID.String(),
				},
			},
		}
		_, err := s.WalletClient.Deposit(depositReq)
		if err != nil {
			slog.Error("Failed to retry deposit", "error", err, "transaction_id", tx.ID)
			return models.TransactionStatusPending
		}
	}

	oldTx.Status = models.TransactionStatusFinalized
	err = s.UpdateTransaction(ctx, oldTx)
	if err != nil {
		slog.Error("failed to update withdraw transaction", "error", err)
	}

	return models.TransactionStatusConfirmed
}

// retryCancel retries a pending cancel transaction
func (s *Service) retryCancel(ctx context.Context, tx *models.Transaction) models.TransactionStatus {
	ogAmount, err := strconv.ParseFloat(tx.Amount, 64)
	if err != nil {
		slog.Error(err.Error())
		return models.TransactionStatusFailed
	}

	// Find the original transaction to understand what to reverse
	originalTx, err := s.Repository.GetTransactionByProviderID(ctx, tx.ProviderID)
	if err != nil {
		slog.Error("Failed to find original transaction for cancel", "error", err, "transaction_id", tx.ID)
		return models.TransactionStatusFailed
	}

	// Reverse the original transaction
	if originalTx.Type == models.TransactionTypeWithdraw {
		// Original was a withdrawal (bet), so we need to deposit back
		depositReq := walletclient.DepositRequest{
			UserID:   int(tx.PlayerID),
			Currency: string(tx.Currency),
			Transactions: []walletclient.DepositRequestTransaction{
				{
					Amount:    ogAmount,
					BetID:     originalTx.ProviderID,
					Reference: fmt.Sprintf("cancel-%d", tx.ProviderID),
				},
			},
		}

		_, err := s.WalletClient.Deposit(depositReq)
		if err != nil {
			slog.Error("Failed to retry cancel deposit", "error", err, "transaction_id", tx.ID)
			return models.TransactionStatusPending
		}

	} else if originalTx.Type == models.TransactionTypeDeposit {
		// Original was a deposit (settle), so we need to withdraw back
		withdrawReq := walletclient.WithdrawRequest{
			UserID:   int(tx.PlayerID),
			Currency: string(tx.Currency),
			Transactions: []walletclient.WithdrawRequestTransaction{
				{
					Amount:    ogAmount,
					BetID:     originalTx.ProviderID,
					Reference: fmt.Sprintf("cancel-%d", tx.ProviderID),
				},
			},
		}

		_, err := s.WalletClient.Withdraw(withdrawReq)
		if err != nil {
			slog.Error("Failed to retry cancel withdrawal", "error", err, "transaction_id", tx.ID)
			return models.TransactionStatusPending
		}
	} else {
		slog.Error("Cannot cancel a cancel transaction", "transaction_id", tx.ID, "original_type", originalTx.Type)
		return models.TransactionStatusFailed
	}

	originalTx.Status = models.TransactionStatusFinalized
	if err := s.Repository.UpdateTransaction(ctx, originalTx); err != nil {
		slog.Error("Failed to update original transaction status", "error", err, "transaction_id", originalTx.ID)
	}

	return models.TransactionStatusConfirmed
}
