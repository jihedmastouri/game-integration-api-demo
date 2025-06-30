package v1

import "github.com/jihedmastouri/game-integration-api-demo/service"

type Handlers struct {
	srv *service.Service
}

func NewHandlers(srv *service.Service) *Handlers {
	return &Handlers{
		srv,
	}
}
