package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/uptrace/bun"
)

type TransactionProvider struct {
	*bun.DB
}

func NewTransactionProvider(db *bun.DB) TransactionProvider {
	return TransactionProvider{db}
}

func (t TransactionProvider) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	_, err := t.NewInsert().Model(transaction).Exec(ctx)
	return err
}

func (t TransactionProvider) GetTransactionByProviderID(ctx context.Context, providerID uint64) (*models.Transaction, error) {
	transaction := &models.Transaction{}
	err := t.NewSelect().Model(transaction).Where("provider_id = ?", providerID).Scan(ctx)
	return transaction, err
}

func (t TransactionProvider) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	_, err := t.NewUpdate().Model(transaction).WherePK().Exec(ctx)
	return err
}

func (t TransactionProvider) GetTransactionByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	transaction := &models.Transaction{}
	err := t.NewSelect().Model(transaction).Where("id = ?", id).Scan(ctx)
	return transaction, err
}

func (t TransactionProvider) GetFirstProcessingTransactionsByPlayerID(ctx context.Context, playerID uint64) (*models.Transaction, error) {
	var transactions *models.Transaction
	err := t.NewSelect().
		Model(&transactions).
		Where("player_id = ? AND status = ?", playerID, models.TransactionStatusProcessing).
		Order("created_at ASC").
		Limit(1).
		Scan(ctx)
	return transactions, err
}

func (t TransactionProvider) GetFirstPendingTransactionsByPlayerID(ctx context.Context, playerID uint64) (*models.Transaction, error) {
	var transactions *models.Transaction
	err := t.NewSelect().
		Model(&transactions).
		Where("player_id = ? AND status = ?", playerID, models.TransactionStatusPending).
		Order("created_at ASC").
		Limit(1).
		Scan(ctx)
	return transactions, err
}

// GetNextProcessableTransaction returns the first pending transaction for a user
// that doesn't have any other transaction currently being processed
func (t TransactionProvider) GetNextProcessableTransaction(ctx context.Context) (*models.Transaction, error) {
	// Find the oldest pending transaction where the player doesn't have any processing transactions
	var transaction models.Transaction
	err := t.NewSelect().
		Model(&transaction).
		Where("status = ?", models.TransactionStatusPending).
		Where("player_id NOT IN (SELECT DISTINCT player_id FROM transactions WHERE status = ?)", models.TransactionStatusProcessing).
		Order("created_at ASC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// StartProcessingTransaction atomically marks a transaction as processing
func (t TransactionProvider) StartProcessingTransaction(ctx context.Context, transactionID uuid.UUID) error {
	// Use a db transaction to atomically check and update
	return t.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// First, verify the transaction is still pending
		var transaction models.Transaction
		err := tx.NewSelect().
			Model(&transaction).
			Where("id = ? AND status = ?", transactionID, models.TransactionStatusPending).
			Scan(ctx)

		if err != nil {
			return err
		}

		// Check if the player has any processing transactions
		var count int
		count, err = tx.NewSelect().
			Model((*models.Transaction)(nil)).
			Where("player_id = ? AND status = ?", transaction.PlayerID, models.TransactionStatusProcessing).
			Count(ctx)
		if err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("player %d already has a transaction being processed", transaction.PlayerID)
		}

		// Update status to processing
		_, err = tx.NewUpdate().
			Model(&transaction).
			Set("status = ?", models.TransactionStatusProcessing).
			Set("updated_at = NOW()").
			Where("id = ?", transactionID).
			Exec(ctx)

		return err
	})
}
