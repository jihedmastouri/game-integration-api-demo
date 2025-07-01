package service

import (
	"github.com/jihedmastouri/game-integration-api-demo/internal"
	"github.com/jihedmastouri/game-integration-api-demo/repository"
	"github.com/jihedmastouri/game-integration-api-demo/service/walletclient"
)

type Service struct {
	repository.Repository
	WalletClient *walletclient.WalletClient
}

func NewService(repo repository.Repository) *Service {
	walletClient := walletclient.NewWalletClient(internal.Config.WALLET_API_URL, internal.Config.WALLET_API_KEY)
	return &Service{
		Repository:   repo,
		WalletClient: walletClient,
	}
}
