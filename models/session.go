package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PlayerSession struct {
	bun.BaseModel `bun:"table:player_sessions,alias:ps"`

	ID        uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Player    *Player   `bun:"rel:belongs-to,join:player_id=id"`
	PlayerID  uint64    `json:"-"`
	ExpiresAt time.Time `bun:"expires_at"`
	IssuedAt  time.Time `bun:"created_at"`
}
