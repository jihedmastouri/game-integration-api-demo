package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/uptrace/bun"
)

type PlayerProvider struct {
	*bun.DB
}

func NewPlayerProvider(db *bun.DB) PlayerProvider {
	return PlayerProvider{db}
}

func (p PlayerProvider) GetPlayerBySession(ctx context.Context, session uuid.UUID) (*models.Player, error) {
	player := &models.Player{}
	err := p.NewSelect().
		Model(player).
		Relation("PlayerSessions", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("ps.id = ?", session)
		}).
		Scan(ctx)
	return player, err
}

func (p PlayerProvider) GetPlayerByID(ctx context.Context, id uint64) (*models.Player, error) {
	player := &models.Player{}
	err := p.NewSelect().Model(player).Where("id = ?", id).Scan(ctx)
	return player, err
}

func (p PlayerProvider) GetPlayerByUsername(ctx context.Context, username string) (*models.Player, error) {
	player := &models.Player{}
	err := p.NewSelect().Model(player).Where("username = ?", username).Scan(ctx)
	return player, err
}

func (p PlayerProvider) CreatePlayer(ctx context.Context, player *models.Player) error {
	_, err := p.NewInsert().Model(player).Exec(ctx)
	return err
}

func (p PlayerProvider) CreatePlayerSession(ctx context.Context, playerID uint64) (*models.PlayerSession, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	PlayerSession := &models.PlayerSession{
		ExpiresAt: expirationTime,
		IssuedAt:  time.Now(),
		PlayerID:  playerID,
	}
	_, err := p.NewInsert().Model(PlayerSession).Exec(ctx)
	return PlayerSession, err
}
