package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/repository/migrations"
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

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err = migrator.Init(ctx)
	if err != nil {
		slog.Error("Migration Init", "error", err)
		return nil, err
	}

	if err := migrator.Lock(ctx); err != nil {
		slog.Error("Migration Lock", "error", err)
		return nil, err
	}
	defer migrator.Unlock(ctx) //nolint:errcheck

	group, err := migrator.Migrate(ctx)
	if err != nil {
		slog.Error("Migration", "error", err)
		return nil, err
	}
	if group.IsZero() {
		slog.Debug("there are no new migrations to run (database is up to date)\n")
	}

	slog.Debug(fmt.Sprintf("migrated to %d\n", group.ID))

	return &RepoPostgresSQLProvider{
		NewPlayerProvider(db),
		NewTransactionProvider(db),
	}, nil
}
