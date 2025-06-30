package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Repository interface {
	PlayerRepository
	TransactionRepository
}

type PlayerRepository interface {
	GetPlayerByID(ctx context.Context, id int) (*models.Player, error)
	GetPlayerByUsername(ctx context.Context, username string) (*models.Player, error)
	GetPlayerBySession(ctx context.Context, session uuid.UUID) (*models.Player, error)

	CreatePlayerSession(ctx context.Context, playerID uint64) (*models.PlayerSession, error)
}

type TransactionRepository interface{}

type RepoPostgresSQLProvider struct {
	PlayerRepository
	TransactionRepository
}

func Connect(databaseUrl string) (*RepoPostgresSQLProvider, error) {
	slog.Debug(databaseUrl)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseUrl)))
	db := bun.NewDB(sqldb, pgdialect.New())

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	return &RepoPostgresSQLProvider{
		NewPlayerProvider(db),
		NewTransactionProvider(db),
	}, nil
}
