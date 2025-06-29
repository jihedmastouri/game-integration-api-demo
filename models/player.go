package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Player struct {
	bun.BaseModel `bun:"table:players,alias:p"`

	ID       uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Balance  int
	Currency Currency
}
