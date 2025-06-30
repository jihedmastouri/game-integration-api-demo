package rest

import (
	"github.com/jihedmastouri/game-integration-api-demo/service"
	v1 "github.com/jihedmastouri/game-integration-api-demo/transport/rest/v1"
	"github.com/labstack/echo/v4"
)

type errorCode string

const (
	ValidationError errorCode = "RequestValidationError"
)

func SetupRoutes(e *echo.Echo, srv *service.Service) {
	v1Handlers := v1.NewHandlers(srv)

	api := e.Group("/api")
	v1Group := api.Group("/v1")

	v1Group.POST("/auth", v1Handlers.Authenticate)

	authv1 := v1Group.Group("", AuthMiddlewareFactory(srv))
	authv1.GET("/user-info/:id", echo.MethodNotAllowedHandler)
	authv1.POST("/withdraw", echo.MethodNotAllowedHandler)
	authv1.POST("/deposit", echo.MethodNotAllowedHandler)
	authv1.POST("/cancel", echo.MethodNotAllowedHandler)
}
