package rest

import (
	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/labstack/echo/v4"
)

type errorCode string

const (
	ValidationError errorCode = "RequestValidationError"
)

func SetupRoutes(e *echo.Echo, srv service.Service) {
	api := e.Group("/api")

	v1 := api.Group("/v1")
	v1.POST("/auth", echo.MethodNotAllowedHandler)
	v1.GET("/user-info/:id", echo.MethodNotAllowedHandler)
	v1.POST("/withdraw", echo.MethodNotAllowedHandler)
	v1.POST("/deposit", echo.MethodNotAllowedHandler)
	v1.POST("/cancel", echo.MethodNotAllowedHandler)
}
