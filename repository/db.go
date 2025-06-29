package repository

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Repository interface {
	PlayerRepository
	TransactionRepository
}

type PlayerRepository interface{}

type TransactionRepository interface{}

type RepoPostgresSQLProvider struct {
	PlayerRepository
	TransactionRepository
}

func Connect(databaseUrl string) RepoPostgresSQLProvider {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseUrl)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return RepoPostgresSQLProvider{
		NewPlayerProvider(db),
		NewTransactionProvider(db),
	}
}
