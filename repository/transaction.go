package repository

import (
	"context"

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

func (t TransactionProvider) GetTransactionByProviderID(ctx context.Context, providerID int) (*models.Transaction, error) {
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
