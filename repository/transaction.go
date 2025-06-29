package repository

import "github.com/uptrace/bun"

type TransactionProvider interface {
	*bun.DB
}

func NewTransactionProvider(db *bun.DB) PlayerProvider {
	return PlayerProvider{db}
}
