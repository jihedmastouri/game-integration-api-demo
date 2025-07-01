package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Player struct {
	bun.BaseModel `bun:"table:players,alias:p"`

	ID        uint64    `bun:",pk,autoincrement"`
	Username  string    `bun:"username"`
	Password  string    `bun:"password"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`

	PlayerSessions []*PlayerSession `bun:"rel:has-many,join:id=player_id"`
	Transactions   []*Transaction   `bun:"rel:has-many,join:id=player_id"`
}
