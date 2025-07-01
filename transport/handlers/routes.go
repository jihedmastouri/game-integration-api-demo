package handlers

import (
	"github.com/jihedmastouri/game-integration-api-demo/service"
	v1 "github.com/jihedmastouri/game-integration-api-demo/transport/handlers/rest_v1"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, srv *service.Service) {
	v1Handlers := v1.NewHandlers(srv)

	api := e.Group("/api")
	v1Group := api.Group("/v1")
	{
		v1Group.POST("/auth", v1Handlers.Authenticate)

		authv1 := v1Group.Group("", AuthMiddlewareFactory(srv))
		{
			authv1.GET("/player-info", v1Handlers.PlayerInfo)
			authv1.POST("/withdraw", echo.MethodNotAllowedHandler)
			authv1.POST("/deposit", echo.MethodNotAllowedHandler)
			authv1.POST("/cancel", echo.MethodNotAllowedHandler)
		}
	}
}
