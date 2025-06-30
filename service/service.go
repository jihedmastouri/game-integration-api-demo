package service

import (
	"github.com/jihedmastouri/game-integration-api-demo/repository"
)

type Service struct {
	repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo,
	}
}
