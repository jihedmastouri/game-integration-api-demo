package repository

import "github.com/uptrace/bun"

type PlayerProvider struct {
	*bun.DB
}

func NewPlayerProvider(db *bun.DB) PlayerProvider {
	return PlayerProvider{db}
}
